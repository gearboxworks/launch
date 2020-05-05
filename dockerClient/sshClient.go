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

func (me *Ssh) Connect() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		sshConfig := &ssh.ClientConfig{}

		var auth []ssh.AuthMethod

		// Try SSH key file first.
		var keyfile ssh.AuthMethod
		keyfile, err = readPublicKeyFile(me.ClientAuth.PublicKey)

		if err == nil && keyfile != nil {
			// Authenticate using SSH key.
			auth = []ssh.AuthMethod{keyfile}
		} else {
			// Authenticate using password
			auth = []ssh.AuthMethod{ssh.Password(me.ClientAuth.Password)}
		}

		sshConfig = &ssh.ClientConfig{
			User: me.ClientAuth.Username,
			Auth: auth,
			// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 10,
		}

		me.ClientInstance, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", me.ClientAuth.Host, me.ClientAuth.Port), sshConfig)
		if err != nil {
			break
		}

		me.ClientSession, err = me.ClientInstance.NewSession()
		defer me.ClientSession.Close()
		defer me.ClientInstance.Close()
		if err != nil {
			break
		}

		// Set IO
		me.ClientSession.Stdout = os.Stdout
		me.ClientSession.Stderr = os.Stderr
		me.ClientSession.Stdin = os.Stdin

		for k, v := range me.Env {
			err = me.ClientSession.Setenv(k, v)
			if err != nil {
				break
			}
		}

		if len(me.CmdArgs) == 0 {
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
				defer terminal.Restore(fileDescriptor, originalState)

				me.StatusLine.TermWidth, me.StatusLine.TermHeight, err = terminal.GetSize(fileDescriptor)
				if err != nil {
					break
				}

				// xterm-256color
				err = me.ClientSession.RequestPty("xterm-256color", me.StatusLine.TermHeight, me.StatusLine.TermWidth, modes)
				if err != nil {
					break
				}
			}

			go me.StatusLineUpdate()
			go me.statusLineWorker()

			// Start remote shell
			err = me.ClientSession.Shell()
			if err != nil {
				break
			}

			err = me.ClientSession.Wait()
			if err != nil {
				break
			}

		} else {
			cmd := ""
			if len(me.CmdArgs) > 0 {
				for _, v := range me.CmdArgs {
					cmd = fmt.Sprintf("%s %s", cmd, v)
				}
			}

			err = me.ClientSession.Run(cmd)
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

		me.resetView()
	}

	return err
}

func (me *Ssh) getEnv() error {
	var err error

	for range only.Once {
		me.Env = make(Environment)
		for _, item := range os.Environ() {
			if strings.HasPrefix(item, "TMPDIR=") {
				continue
			}

			s := strings.SplitN(item, "=", 2)
			me.Env[s[0]] = s[1]
		}
	}

	return err
}

// StatusLineWorker() - handles the actual updates to the status line
func (me *Ssh) StatusLineUpdate() {

	me.setView()
	// w := gob.NewEncoder(me.Session)
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, syscall.SIGWINCH)

	for me.StatusLine.TerminateFlag == false {
		// Handle terminal windows size changes properly.
		fileDescriptor := int(os.Stdin.Fd())
		width, height, _ := terminal.GetSize(fileDescriptor)
		if (me.StatusLine.TermWidth != width) || (me.StatusLine.TermHeight != height) {
			me.StatusLine.TermWidth = width
			me.StatusLine.TermHeight = height
			// me.Session.Signal(syscall.SIGWINCH)
			_ = me.ClientSession.WindowChange(height, width)
		} else {
			// Only update if we haven't seen a SIGWINCH - just to wait for things to settle.
			me.displayStatusLine()
		}

		time.Sleep(me.StatusLine.UpdateDelay)
	}

}

func (me *Ssh) SetStatusLine(text string) {

	me.StatusLine.Text = text
}

func (me *Ssh) displayStatusLine() {
	const savePos = "\033[s"
	const restorePos = "\033[u"
	bottomPos := fmt.Sprintf("\033[%d;0H", me.StatusLine.TermHeight)
	// topPos := fmt.Sprintf("\033[0;0H")

	if me.StatusLine.Enable {
		fmt.Printf("%s%s%s%s", savePos, bottomPos, me.StatusLine.Text, restorePos)
	}
}

func (me *Ssh) setView() {
	const clearScreen = "\033[H\033[2J"
	scrollFixBottom := fmt.Sprintf("\033[1;%dr", me.StatusLine.TermHeight-1)
	// scrollFixTop := fmt.Sprintf("\033[2;%dr", termHeight)

	if me.StatusLine.Enable {
		fmt.Printf(scrollFixBottom)
		fmt.Printf(clearScreen)
	}
}

func (me *Ssh) resetView() {
	const savePos = "\033[s"
	const restorePos = "\033[u"
	scrollFixBottom := fmt.Sprintf("\033[1;%dr", me.StatusLine.TermHeight)
	// scrollFixTop := fmt.Sprintf("\033[2;%dr", termHeight)

	if me.StatusLine.Enable {
		fmt.Printf(savePos)
		fmt.Printf(scrollFixBottom)
		fmt.Printf(restorePos)

		me.StatusLine.Text = ""
		for i := 0; i <= me.StatusLine.TermWidth; i++ {
			me.StatusLine.Text += " "
		}
		me.displayStatusLine()
	}
}

func stripAnsi(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)

	return re.ReplaceAllString(str, "")
}

// Example host worker. This periodically changes the me.StatusLine.Text from the host side.
// The StatusLineWorker() will update the bottom line using the me.StatusLine.Text.
func (me *Ssh) statusLineWorker() {

	yellow := color.New(color.BgBlack, color.FgHiYellow).SprintFunc()
	magenta := color.New(color.BgBlack, color.FgHiMagenta).SprintFunc()
	green := color.New(color.BgBlack, color.FgHiGreen).SprintFunc()
	//normal := color.New(color.BgWhite, color.FgHiBlack).SprintFunc()

	for me.StatusLine.TerminateFlag == false {
		//now := time.Now()
		//dateStr := normal("Date:") + " " + yellow(fmt.Sprintf("%.4d/%.2d/%.2d", now.Year(), now.Month(), now.Day()))
		//timeStr := normal("Time:") + " " + magenta(fmt.Sprintf("%.2d:%.2d:%.2d", now.Hour(), now.Minute(), now.Second()))
		statusStr := yellow("Status:") + " " + green("OK")
		infoStr := yellow("Gearbox container:") + " " + magenta(me.GearName + ":" + me.GearVersion)

		//line := fmt.Sprintf("%s	%s %s", statusStr, dateStr, timeStr)
		line := fmt.Sprintf("%s - %s", infoStr, statusStr)

		// Add spaces to ensure it's right justified.
		spaces := ""
		lineLen := len(stripAnsi(line))
		for i := 0; i < me.StatusLine.TermWidth-lineLen; i++ {
			spaces += " "
		}

		me.SetStatusLine(spaces + line) // + fmt.Sprintf("W:%d L:%d S:%d C:%d", me.StatusLine.TermWidth, len(line), len(spaces), lineLen))

		time.Sleep(time.Second * 5)
	}
}

func (ssh *Ssh) EnsureNotNil() error {
	var err error

	if ssh == nil {
		err = errors.New("unexpected error")
	}
	return err
}
