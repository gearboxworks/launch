package dockerClient

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"launch/only"
	"net"
	"os"
	"regexp"
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

	//
	Debug      bool
}
type SshClientArgs Ssh

type Environment map[string]string

type StatusLine struct {
	Text          string
	Enable        bool
	UpdateDelay   time.Duration
	TermWidth     int
	TermHeight    int
	TerminateFlag bool
}

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

	// Query VB to see if it exists.
	// If not return nil.

	return sshClient
}

func (s *Ssh) Connect() error {
	var err error

	for range only.Once {
		err = s.EnsureNotNil()
		if err != nil {
			break
		}

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

func (s *Ssh) getEnv() error {
	var err error

	for range only.Once {
		s.Env = make(Environment)
		for _, item := range os.Environ() {
			if strings.HasPrefix(item, "TMPDIR=") {
				continue
			}

			sa := strings.SplitN(item, "=", 2)
			s.Env[sa[0]] = sa[1]
		}
	}

	return err
}

// StatusLineWorker() - handles the actual updates to the status line
func (s *Ssh) StatusLineUpdate() {

	s.setView()
	// w := gob.NewEncoder(s.Session)
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, syscall.SIGWINCH)

	for s.StatusLine.TerminateFlag == false {
		// Handle terminal windows size changes properly.
		fileDescriptor := int(os.Stdin.Fd())
		width, height, _ := terminal.GetSize(fileDescriptor)
		if (s.StatusLine.TermWidth != width) || (s.StatusLine.TermHeight != height) {
			s.StatusLine.TermWidth = width
			s.StatusLine.TermHeight = height
			// s.Session.Signal(syscall.SIGWINCH)
			_ = s.ClientSession.WindowChange(height, width)
		} else {
			// Only update if we haven't seen a SIGWINCH - just to wait for things to settle.
			s.displayStatusLine()
		}

		time.Sleep(s.StatusLine.UpdateDelay)
	}

}

func (s *Ssh) SetStatusLine(text string) {

	s.StatusLine.Text = text
}

func (s *Ssh) displayStatusLine() {
	const savePos = "\033[s"
	const restorePos = "\033[u"
	bottomPos := fmt.Sprintf("\033[%d;0H", s.StatusLine.TermHeight)
	// topPos := fmt.Sprintf("\033[0;0H")

	if s.StatusLine.Enable {
		fmt.Printf("%s%s%s%s", savePos, bottomPos, s.StatusLine.Text, restorePos)
	}
}

func (s *Ssh) setView() {
	const clearScreen = "\033[H\033[2J"
	scrollFixBottom := fmt.Sprintf("\033[1;%dr", s.StatusLine.TermHeight-1)
	// scrollFixTop := fmt.Sprintf("\033[2;%dr", termHeight)

	if s.StatusLine.Enable {
		fmt.Printf(scrollFixBottom)
		fmt.Printf(clearScreen)
	}
}

func (s *Ssh) resetView() {
	const savePos = "\033[s"
	const restorePos = "\033[u"
	scrollFixBottom := fmt.Sprintf("\033[1;%dr", s.StatusLine.TermHeight)
	// scrollFixTop := fmt.Sprintf("\033[2;%dr", termHeight)

	if s.StatusLine.Enable {
		fmt.Printf(savePos)
		fmt.Printf(scrollFixBottom)
		fmt.Printf(restorePos)

		s.StatusLine.Text = ""
		for i := 0; i <= s.StatusLine.TermWidth; i++ {
			s.StatusLine.Text += " "
		}
		s.displayStatusLine()
	}
}

func stripAnsi(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)

	return re.ReplaceAllString(str, "")
}

// Example host worker. This periodically changes the me.StatusLine.Text from the host side.
// The StatusLineWorker() will update the bottom line using the me.StatusLine.Text.
func (s *Ssh) statusLineWorker() {

	yellow := color.New(color.BgBlack, color.FgHiYellow).SprintFunc()
	magenta := color.New(color.BgBlack, color.FgHiMagenta).SprintFunc()
	green := color.New(color.BgBlack, color.FgHiGreen).SprintFunc()
	//normal := color.New(color.BgWhite, color.FgHiBlack).SprintFunc()

	for s.StatusLine.TerminateFlag == false {
		//now := time.Now()
		//dateStr := normal("Date:") + " " + yellow(fmt.Sprintf("%.4d/%.2d/%.2d", now.Year(), now.Month(), now.Day()))
		//timeStr := normal("Time:") + " " + magenta(fmt.Sprintf("%.2d:%.2d:%.2d", now.Hour(), now.Minute(), now.Second()))
		statusStr := yellow("Status:") + " " + green("OK")
		infoStr := yellow("Gearbox container:") + " " + magenta(s.GearName + ":" + s.GearVersion)

		//line := fmt.Sprintf("%s	%s %s", statusStr, dateStr, timeStr)
		line := fmt.Sprintf("%s - %s", infoStr, statusStr)

		// Add spaces to ensure it's right justified.
		spaces := ""
		lineLen := len(stripAnsi(line))
		for i := 0; i < s.StatusLine.TermWidth-lineLen; i++ {
			spaces += " "
		}

		s.SetStatusLine(spaces + line) // + fmt.Sprintf("W:%d L:%d S:%d C:%d", s.StatusLine.TermWidth, len(line), len(spaces), lineLen))

		time.Sleep(time.Second * 5)
	}
}

func (s *Ssh) EnsureNotNil() error {
	var err error

	if s == nil {
		err = errors.New("unexpected error")
	}
	return err
}
