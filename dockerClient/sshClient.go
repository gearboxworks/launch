package dockerClient

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"launch/ux"
	"net"
	"os"
	"strings"
	"time"
)

type Ssh struct {
	// SSH client
	ClientInstance *ssh.Client
	ClientSession  *ssh.Session
	ClientAuth     *SshAuth

	// SSH server for SSHFS
	ServerConfig      *ssh.ServerConfig
	ServerListener    net.Listener
	ServerConnection  net.Conn
	ServerAuth        *SshAuth
	FsReadOnly        bool
	FsMount           string

	// Status line related.
	StatusLine  StatusLine
	GearName    string
	GearVersion string

	// Shell related.
	Shell      bool
	Env        Environment
	CmdArgs    []string

	State      *ux.State
	Debug      bool
}
type SshClientArgs Ssh

type Environment map[string]string

const DefaultUsername = "gearbox"
const DefaultPassword = "box"
const DefaultKeyFile = "./keyfile.pub"
const DefaultSshHost = "localhost"
const DefaultSshPort = "22"
const DefaultStatusLineUpdateDelay = time.Second * 2


func NewSshClient(args ...SshClientArgs) *Ssh {

	var _args SshClientArgs
	if len(args) > 0 {
		_args = args[0]
	}

	_args.ClientAuth = NewSshAuth(*_args.ClientAuth)

	if _args.StatusLine.UpdateDelay == 0 {
		_args.StatusLine.UpdateDelay = DefaultStatusLineUpdateDelay
	}

	sshClient := &Ssh{}
	*sshClient = Ssh(_args)

	return sshClient
}


func (s *Ssh) IsNil() *ux.State {
	if state := ux.IfNilReturnError(s); state.IsError() {
		return state
	}
	s.State = s.State.EnsureNotNil()
	return s.State
}

func (s *Ssh) IsValid() *ux.State {
	if state := ux.IfNilReturnError(s); state.IsError() {
		return state
	}

	for range OnlyOnce {
		s.State = s.State.EnsureNotNil()

		if s.GearName == "" {
			s.State.SetError("name is nil")
			break
		}

		if s.GearVersion == "" {
			s.State.SetError("version is nil")
			break
		}
	}

	return s.State
}


func (s *Ssh) Connect() error {
	var err error
	if state := s.IsNil(); state.IsError() {
		return state.GetError()
	}

	for range OnlyOnce {
		sshConfig := &ssh.ClientConfig{}

		var auth []ssh.AuthMethod

		// Try SSH key file first.
		var keyfile ssh.AuthMethod
		keyfile, err = readPublicKeyFile(s.ClientAuth.PublicKey)

		if err == nil && keyfile != nil {
			// Authenticate using SSH key.
			auth = []ssh.AuthMethod{keyfile}
		} else {
			// Authenticate using password
			auth = []ssh.AuthMethod{ssh.Password(s.ClientAuth.Password)}
		}

		sshConfig = &ssh.ClientConfig{
			User: s.ClientAuth.Username,
			Auth: auth,
			// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 10,
		}

		s.ClientInstance, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.ClientAuth.Host, s.ClientAuth.Port), sshConfig)
		if err != nil {
			break
		}

		s.ClientSession, err = s.ClientInstance.NewSession()
		//noinspection GoDeferInLoop,GoUnhandledErrorResult
		defer s.ClientSession.Close()
		//noinspection GoDeferInLoop,GoUnhandledErrorResult
		defer s.ClientInstance.Close()
		if err != nil {
			break
		}

		// Set IO
		s.ClientSession.Stdout = os.Stdout
		s.ClientSession.Stderr = os.Stderr
		s.ClientSession.Stdin = os.Stdin

		for k, v := range s.Env {
			err = s.ClientSession.Setenv(k, v)
			if err != nil {
				break
			}
		}

		if len(s.CmdArgs) == 0 {
			// Set up terminal modes
			modes := ssh.TerminalModes{
				ssh.ECHO:          1,
				ssh.TTY_OP_ISPEED: 19200,
				ssh.TTY_OP_OSPEED: 19200,
			}

			// Request pseudo terminal
			fileDescriptor := int(os.Stdin.Fd())
			if terminal.IsTerminal(fileDescriptor) {
				originalState, err := terminal.MakeRaw(fileDescriptor)
				if err != nil {
					break
				}
				//noinspection GoDeferInLoop,GoUnhandledErrorResult
				defer terminal.Restore(fileDescriptor, originalState)

				s.StatusLine.TermWidth, s.StatusLine.TermHeight, err = terminal.GetSize(fileDescriptor)
				if err != nil {
					break
				}

				// xterm-256color
				err = s.ClientSession.RequestPty("xterm-256color", s.StatusLine.TermHeight, s.StatusLine.TermWidth, modes)
				if err != nil {
					break
				}
			}

			go s.StatusLineUpdate()
			go s.statusLineWorker()

			// Start remote shell
			err = s.ClientSession.Shell()
			if err != nil {
				break
			}

			err = s.ClientSession.Wait()
			if err != nil {
				break
			}

		} else {
			cmd := ""
			if len(s.CmdArgs) > 0 {
				for _, v := range s.CmdArgs {
					cmd = fmt.Sprintf("%s %s", cmd, v)
				}
			}

			err = s.ClientSession.Run(cmd)
			if err != nil {
				break
			}
		}

		/*
			// Loop around input <-> output.
			for {
				reader := bufio.NewReader(os.Stdin)
				str, _ := reader.ReadString('\n')
				fmt.Fprint(in, str)
			}
		*/

		s.resetView()
	}

	return err
}


func (s *Ssh) getEnv() *ux.State {
	if state := s.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		s.Env = make(Environment)
		for _, item := range os.Environ() {
			if strings.HasPrefix(item, "TMPDIR=") {
				continue
			}

			sa := strings.SplitN(item, "=", 2)
			s.Env[sa[0]] = sa[1]
		}
	}

	return s.State
}
