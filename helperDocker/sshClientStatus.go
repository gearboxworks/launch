package helperDocker

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"time"
)


type StatusLine struct {
	Text          string
	Enable        bool
	UpdateDelay   time.Duration
	TermWidth     int
	TermHeight    int
	TerminateFlag bool
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
