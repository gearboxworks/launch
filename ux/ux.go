// Other possibilities:
//
// https://github.com/nsf/termbox-go
// https://github.com/jroimartin/gocui
// https://github.com/marcusolsson/tui-go
// https://github.com/rivo/tview
// https://github.com/gizak/termui
// https://github.com/logrusorgru/aurora
package ux

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)


type UxGetter interface {
	Open() error
	Close()
	PrintfWhite(format string, args ...interface{})
	PrintfCyan(format string, args ...interface{})
	PrintfYellow(format string, args ...interface{})
	PrintfRed(format string, args ...interface{})
	PrintfGreen(format string, args ...interface{})
	PrintfBlue(format string, args ...interface{})
	PrintfMagenta(format string, args ...interface{})
	SprintfWhite(format string, args ...interface{}) string
	SprintfCyan(format string, args ...interface{}) string
	SprintfYellow(format string, args ...interface{}) string
	SprintfRed(format string, args ...interface{}) string
	SprintfGreen(format string, args ...interface{}) string
	SprintfBlue(format string, args ...interface{}) string
	SprintfMagenta(format string, args ...interface{}) string
	Sprintf(format string, args ...interface{}) string
	Printf(format string, args ...interface{})
	SprintfOk(format string, args ...interface{}) string
	PrintfOk(format string, args ...interface{})
	SprintfWarning(format string, args ...interface{}) string
	PrintfWarning(format string, args ...interface{})
	SprintfError(format string, args ...interface{}) string
	PrintfError(format string, args ...interface{})
	SprintError(err error) string
	PrintError(err error)
	GetTerminalSize() (int, int, error)
}

type Ux struct {
}


//noinspection GoUnusedGlobalVariable
var Color aurora.Aurora
var _defined bool
var _name string


func Open(name string) error {
	var err error

	for range OnlyOnce {
		Color = aurora.NewAurora(true)
		_name = name
		if name == "" {
			name = "Gearbox: "
		}

		//err = termui.Init();
		//if err != nil {
		//	fmt.Printf("failed to initialize termui: %v", err)
		//	break
		//}

		_defined = true
	}

	return err
}


func Close() {
	if _defined {
		//termui.Close()
	}
}


func PrintfWhite(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s", aurora.BrightWhite(inline), aurora.Reset(""))
}
func PrintfCyan(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s", aurora.BrightCyan(inline), aurora.Reset(""))
}
func PrintfYellow(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s", aurora.BrightYellow(inline), aurora.Reset(""))
}
func PrintfRed(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s", aurora.BrightRed(inline), aurora.Reset(""))
}
func PrintfGreen(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s", aurora.BrightGreen(inline), aurora.Reset(""))
}
func PrintfBlue(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s", aurora.BrightBlue(inline), aurora.Reset(""))
}
func PrintfMagenta(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s", aurora.BrightMagenta(inline), aurora.Reset(""))
}

func PrintflnWhite(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", aurora.BrightWhite(inline), aurora.Reset(""))
}
func PrintflnCyan(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", aurora.BrightCyan(inline), aurora.Reset(""))
}
func PrintflnYellow(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", aurora.BrightYellow(inline), aurora.Reset(""))
}
func PrintflnRed(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", aurora.BrightRed(inline), aurora.Reset(""))
}
func PrintflnGreen(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", aurora.BrightGreen(inline), aurora.Reset(""))
}
func PrintflnBlue(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", aurora.BrightBlue(inline), aurora.Reset(""))
}
func PrintflnMagenta(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", aurora.BrightMagenta(inline), aurora.Reset(""))
}


func SprintfWhite(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s", aurora.BrightWhite(inline), aurora.Reset(""))
}
func SprintfCyan(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s", aurora.BrightCyan(inline), aurora.Reset(""))
}
func SprintfYellow(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s", aurora.BrightYellow(inline), aurora.Reset(""))
}
func SprintfRed(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s", aurora.BrightRed(inline), aurora.Reset(""))
}
func SprintfGreen(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s", aurora.BrightGreen(inline), aurora.Reset(""))
}
func SprintfBlue(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s", aurora.BrightBlue(inline), aurora.Reset(""))
}
func SprintfMagenta(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s", aurora.BrightMagenta(inline), aurora.Reset(""))
}


func Sprintf(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s%s", aurora.BrightCyan(_name).Bold(), inline, aurora.Reset(""))
}
func Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, Sprintf(format, args...))
}


//func SprintfNormal(format string, args ...interface{}) string {
//	return fmt.Sprintf(format, args...)
//}
func SprintfNormal(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return Sprintf("%s", aurora.BrightBlue(inline))
}
//func PrintfNormal(format string, args ...interface{}) {
//	_, _ = fmt.Fprintf(os.Stdout, fmt.Sprintf(format, args...))
//}
func PrintfNormal(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfNormal(format, args...))
}
func PrintflnNormal(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfNormal(format + "\n", args...))
}


func SprintfOk(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return Sprintf("%s", aurora.BrightGreen(inline))
}
func PrintfOk(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfOk(format, args...))
}
func PrintflnOk(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfOk(format + "\n", args...))
}


func SprintfDebug(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
func PrintfDebug(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf(format + "\n", args...))
}


func SprintfWarning(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return Sprintf("%s", aurora.BrightYellow(inline))
}
func PrintfWarning(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfWarning(format, args...))
}
func PrintflnWarning(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfWarning(format + "\n", args...))
}


func SprintfError(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return Sprintf("%s", aurora.BrightRed(inline))
}
func PrintfError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, SprintfError(format, args...))
}
func PrintflnError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, SprintfError(format + "\n", args...))
}


func SprintError(err error) string {
	var s string

	for range OnlyOnce {
		if err == nil {
			break
		}

		s = Sprintf("%s%s\n", aurora.BrightRed("ERROR: ").Framed(), aurora.BrightRed(err).Framed().SlowBlink().BgBrightWhite())
	}

	return s
}
func PrintError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, SprintError(err))
}


func GetTerminalSize() (int, int, error) {
	var width int
	var height int
	var err error

	fileDescriptor := int(os.Stdin.Fd())
	if terminal.IsTerminal(fileDescriptor) {
		width, height, err = terminal.GetSize(fileDescriptor)
	} else {
		err = errors.New("not a terminal")
	}

	return width, height, err
}
