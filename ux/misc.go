package ux

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/rivo/tview"
	"io/ioutil"
	"path/filepath"
	"strings"
)


//noinspection GoUnusedFunction
func _Printf2(format string, args ...interface{}) {

	for range OnlyOnce {
		var w int
		var h int
		var err error

		w, h, err = GetTerminalSize()
		if err != nil {
			w = 80
			h = 24
		}

		l := widgets.NewParagraph()
		//l.Title = "[color](fg:green,bg:black)deploywp[color](fg:white,bg:black)"
		//l.TextStyle = termui.NewStyle(termui.ColorYellow)
		l.Text = "[deploywp:](fg:green)[ ](fg:white)" + fmt.Sprintf(format, args...) + "[ ](fg:white)"
		l.WrapText = false
		l.SetRect(-1, -1, w, h)
		l.Border = false
		l.BorderLeft = false
		l.BorderRight = false
		l.BorderTop = false
		l.BorderBottom = false
		l.TextStyle.Modifier = termui.ModifierBold

		fmt.Printf("l.Dx: %v\n", l.Dx())
		fmt.Printf("l.Dy: %v\n", l.Dy())

		//foo1 := l.Bounds()
		//foo2 := l.Output()
		//fmt.Printf("l: %v\n", foo1)
		//fmt.Printf("l: %v\n", foo2)

		termui.Render(l)

		fmt.Printf("\n")
	}
}

//noinspection GoUnusedFunction
func _Printf3(format string, args ...interface{}) {

	for range OnlyOnce {
		var w int
		var h int
		var err error

		w, h, err = GetTerminalSize()
		if err != nil {
			w = 80
			h = 24
		}
		fmt.Printf("w: %v\n", w)
		fmt.Printf("h: %v\n", h)

		//app := tview.NewApplication()
		table := tview.NewTextView()
		table.SetBorder(false)
		table.SetText("[deploywp:](fg:green)[ ](fg:white)" + fmt.Sprintf(format, args...) + "[ ](fg:white)")

		fmt.Printf("\n")
	}
}

//noinspection GoUnusedFunction
func _Draw2() error {
	var err error

	for range OnlyOnce {
		p := widgets.NewParagraph()
		p.Text = "Hello World!"
		p.SetRect(0, 0, 25, 5)

		termui.Render(p)

		for e := range termui.PollEvents() {
			if e.Type == termui.KeyboardEvent {
				break
			}
		}
	}

	return err
}

//noinspection GoUnusedFunction
func _Draw3() {
	l := widgets.NewList()
	l.Title = "List"
	l.Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}
	l.TextStyle = termui.NewStyle(termui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)

	termui.Render(l)

	previousKey := ""
	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>", "<Escape>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		termui.Render(l)
	}
}

//noinspection GoUnusedFunction
func _Draw4() error {
	var err error

	for range OnlyOnce {
		rootDir := "."
		root := tview.NewTreeNode(rootDir).
			SetColor(tcell.ColorRed)
		tree := tview.NewTreeView().
			SetRoot(root).
			SetCurrentNode(root)

		// A helper function which adds the files and directories of the given path
		// to the given target node.
		add := func(target *tview.TreeNode, path string) {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				panic(err)
			}

			for _, file := range files {
				node := tview.NewTreeNode(file.Name()).
					SetReference(filepath.Join(path, file.Name())).
					SetSelectable(file.IsDir())
				if file.IsDir() {
					node.SetColor(tcell.ColorGreen)
				}
				target.AddChild(node)
			}
		}

		// Add the current directory to the root node.
		add(root, rootDir)

		// If a directory was selected, open it.
		tree.SetSelectedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			if reference == nil {
				return // Selecting the root node does nothing.
			}
			children := node.GetChildren()
			if len(children) == 0 {
				// Load and show files in this directory.
				path := reference.(string)
				add(node, path)
			} else {
				// Collapse if visible, expand if collapsed.
				node.SetExpanded(!node.IsExpanded())
			}
		})

		err = tview.NewApplication().SetRoot(tree, true).EnableMouse(true).Run()
		if err != nil {
			break
		}
	}

	return err
}

//noinspection GoUnusedFunction
func _Draw5() error {
	var err error

	for range OnlyOnce {
		app := tview.NewApplication()
		table := tview.NewTable().
			SetBorders(true)
		lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
		cols, rows := 10, 40
		word := 0
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				color := tcell.ColorWhite
				if c < 1 || r < 1 {
					color = tcell.ColorYellow
				}
				table.SetCell(r, c,
					tview.NewTableCell(lorem[word]).
						SetTextColor(color).
						SetAlign(tview.AlignCenter))
				word = (word + 1) % len(lorem)
			}
		}
		table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.Stop()
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
			table.GetCell(row, column).SetTextColor(tcell.ColorRed)
			table.SetSelectable(false, false)
		})

		err = app.SetRoot(table, true).EnableMouse(true).Run()
		if err != nil {
			break
		}
	}

	return err
}


// 	"github.com/marcusolsson/tui-go"
//func Draw() {
//
//	for range OnlyOnce {
//		t := tui.NewTheme()
//		normal := tui.Style{Bg: tui.ColorWhite, Fg: tui.ColorBlack}
//		t.SetStyle("normal", normal)
//
//		// A simple label.
//		okay := tui.NewLabel("Everything is fine.")
//
//		// A list with some items selected.
//		l := tui.NewList()
//		l.SetFocused(true)
//		l.AddItems("First row", "Second row", "Third row", "Fourth row", "Fifth row", "Sixth row")
//		l.SetSelected(0)
//
//		t.SetStyle("list.item", tui.Style{Bg: tui.ColorCyan, Fg: tui.ColorMagenta})
//		t.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorRed, Fg: tui.ColorWhite})
//
//		// The style name is appended to the widget name to support coloring of
//		// individual labels.
//		warning := tui.NewLabel("WARNING: This is a warning")
//		warning.SetStyleName("warning")
//		t.SetStyle("label.warning", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorYellow})
//
//		fatal := tui.NewLabel("FATAL: Cats and dogs are now living together.")
//		fatal.SetStyleName("fatal")
//		t.SetStyle("label.fatal", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed})
//
//		// Styles inherit properties of the parent widget by default;
//		// setting a property overrides only that property.
//		message1 := tui.NewLabel("This is an ")
//		emphasis := tui.NewLabel("important")
//		message2 := tui.NewLabel(" message from our sponsors.")
//		message := &StyledBox{
//			Style: "bsod",
//			Box:   tui.NewHBox(message1, emphasis, message2, tui.NewSpacer()),
//		}
//
//		emphasis.SetStyleName("emphasis")
//		t.SetStyle("label.emphasis", tui.Style{Bold: tui.DecorationOn, Underline: tui.DecorationOn, Bg: tui.ColorRed})
//		t.SetStyle("bsod", tui.Style{Bg: tui.ColorCyan, Fg: tui.ColorWhite})
//
//		// Another unstyled label.
//		okay2 := tui.NewLabel("Everything is still fine.")
//
//		root := tui.NewVBox(okay, l, warning, fatal, message, okay2)
//
//		ui, err := tui.New(root)
//		if err != nil {
//			break
//		}
//
//		ui.SetTheme(t)
//		ui.SetKeybinding("Esc", func() { ui.Quit() })
//
//		if err := ui.Run(); err != nil {
//			break
//		}
//	}
//}
//
//// StyledBox is a Box with an overriden Draw method.
//// Embedding a Widget within another allows overriding of some behaviors.
//type StyledBox struct {
//	Style string
//	*tui.Box
//}
//
//// Draw decorates the Draw call to the widget with a style.
//func (s *StyledBox) Draw(p *tui.Painter) {
//	p.WithStyle(s.Style, func(p *tui.Painter) {
//		s.Box.Draw(p)
//	})
//}
