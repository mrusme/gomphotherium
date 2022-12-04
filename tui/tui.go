package tui

import (
	"fmt"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/mattn/go-mastodon"
	"github.com/mrusme/gomphotherium/mast"

	"github.com/pookzilla/go-termd"
)

type ModeType int

const (
	NormalMode  ModeType = 1
	CommandMode          = 2
)

type TUIOptions struct {
	ShowImages     bool
	AutoCompletion bool
	JustifyText    bool
}

type TUICore struct {
	Client  *mastodon.Client
	App     *tview.Application
	CmdLine *tview.InputField
	Profile *tview.TextView
	Stream  *tview.TextView
	History *History
	Grid    *tview.Grid

	Prompt string
	Mode   ModeType

	Timeline             mast.Timeline
	RenderedTimelineType mast.TimelineType

	Options  TUIOptions
	Progress *ProgressManager
	Help     string
}

func TUI(tuiCore TUICore) {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorDefault,
		ContrastBackgroundColor:     tcell.ColorTeal,
		MoreContrastBackgroundColor: tcell.ColorTeal,
		BorderColor:                 tcell.ColorWhite,
		TitleColor:                  tcell.ColorWhite,
		GraphicsColor:               tcell.ColorWhite,
		PrimaryTextColor:            tcell.ColorDefault,
		SecondaryTextColor:          tcell.ColorBlue,
		TertiaryTextColor:           tcell.ColorGreen,
		InverseTextColor:            tcell.ColorBlack,
		ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
	}

	tuiCore.App = tview.NewApplication()

	tuiCore.Timeline = mast.NewTimeline(tuiCore.Client)
	tuiCore.RenderedTimelineType = mast.TimelineHome
	tuiCore.Timeline.Switch(mast.TimelineHome, nil)

	tuiCore.CmdLine = tview.NewInputField()
	tuiCore.Progress = NewProgress(tuiCore.CmdLine, tuiCore.App)

	tuiCore.CmdLine.
		SetLabelColor(tcell.ColorDefault).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				cmd := tuiCore.CmdLine.GetText()
				tuiCore.CmdLine.SetText("")
				tuiCore.Progress.Run(func() (mast.CmdReturnCode, bool) {

					result := mast.CmdProcessor(&tuiCore.Timeline, cmd, mast.TriggerTUI)

					retCode, err, reloadTimeline := result.Decompose()

					tuiCore.History.AddHistory(cmd, retCode, err)

					switch retCode {
					case mast.CodeOk:
						if tuiCore.Timeline.GetCurrentType() == mast.TimelineUser &&
							tuiCore.RenderedTimelineType != mast.TimelineUser {
							tuiCore.Grid.
								RemoveItem(tuiCore.Stream).
								AddItem(tuiCore.Profile, 0, 0, 1, 1, 0, 0, false).
								AddItem(tuiCore.Stream, 1, 0, 1, 1, 0, 0, false)
						} else if tuiCore.RenderedTimelineType == mast.TimelineUser &&
							tuiCore.Timeline.GetCurrentType() != mast.TimelineUser {
							tuiCore.Grid.
								RemoveItem(tuiCore.Profile).
								RemoveItem(tuiCore.Stream).
								AddItem(tuiCore.Stream, 0, 0, 2, 1, 0, 0, false)
						}
						tuiCore.UpdateTimeline(true)
					case mast.CodeHelp:
						tuiCore.ShowHelp()
					case mast.CodeHistory:
						tuiCore.ShowHistory()
					case mast.CodeQuit:
						tuiCore.App.Stop()
					}

					return retCode, reloadTimeline
				})
			}
		})

	if tuiCore.Options.AutoCompletion == true {
		tuiCore.CmdLine.SetAutocompleteFunc(func(input string) []string {
			return mast.CmdAutocompleter(input, tuiCore.Timeline.KnownUsers)
		})
	}

	tuiCore.Profile = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetScrollable(false)

	tuiCore.Stream = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true)

	tuiCore.History = NewHistory(&tuiCore)

	tuiCore.Grid = tview.NewGrid().
		SetRows(8, 0, 1).
		SetColumns(0).
		SetBorders(true).
		AddItem(tuiCore.Stream, 0, 0, 2, 1, 0, 0, false).
		AddItem(tuiCore.CmdLine, 2, 0, 1, 1, 0, 0, false)

	tuiCore.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlR:
			tuiCore.UpdateTimeline(true)
			return nil
		case tcell.KeyRune:
			eventRune := event.Rune()
			switch eventRune {
			case ':':
				if tuiCore.EnterCommandMode() == true {
					return nil
				} else {
					if len(tuiCore.CmdLine.GetText()) == 0 {
						return nil
					}
				}
			case 'u', 'd', 'b', 'f':
				if tuiCore.Mode == NormalMode {
					_, _, _, h := tuiCore.Stream.Box.GetRect()
					currentLine, _ := tuiCore.Stream.GetScrollOffset()

					var scrollLength int = 0
					if eventRune == 'u' || eventRune == 'd' {
						scrollLength = int((h / 2))
					} else if eventRune == 'b' || eventRune == 'f' {
						scrollLength = h
					}

					var scrollTo int = currentLine

					if eventRune == 'u' || eventRune == 'b' {
						scrollTo = currentLine - scrollLength
					} else if eventRune == 'd' || eventRune == 'f' {
						scrollTo = currentLine + scrollLength
					}

					if scrollTo < 0 {
						scrollTo = 0
					}

					tuiCore.Stream.ScrollTo(scrollTo, 0)
				}
			}
		case tcell.KeyEscape:
			if tuiCore.ExitCommandMode(false) == true {
				return nil
			}
		case tcell.KeyCtrlQ:
			tuiCore.App.Stop()
		}

		return event
	})

	go func() {
		first := true

		for {
			time.Sleep(time.Second * 2)
			tuiCore.Progress.Run(func() (mast.CmdReturnCode, bool) {
				result := tuiCore.UpdateTimeline(first)
				first = false
				if tuiCore.Mode == 0 {
					tuiCore.ExitCommandMode(true)
				}

				tuiCore.App.Draw()
				if result {
					return mast.CodeOk, false
				} else {
					return mast.CodeNotOk, false
				}
			})
			time.Sleep(time.Second * 58)
		}
	}()

	if err := tuiCore.App.SetRoot(tuiCore.Grid, true).Run(); err != nil {
		panic(err)
	}
}

func (tuiCore *TUICore) ShowHistory() {
	tuiCore.App.SetRoot(tuiCore.History.Root, true)
}

func (tuiCore *TUICore) ShowHelp() {
	var c termd.Compiler

	help := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetDoneFunc(func(key tcell.Key) {
			tuiCore.App.SetRoot(tuiCore.Grid, true)
			tuiCore.EnterCommandMode()
			return
		})

	// c.SyntaxHighlighter = termd.SyntaxTheme{
	//   "keyword": termd.Style{},
	//   "comment": termd.Style{
	//     Color: "#323232",
	//   },
	//   "literal": termd.Style{
	//     Color: "#555555",
	//   },
	//   "name": termd.Style{
	//     Color: "#777777",
	//   },
	//   "name.function": termd.Style{
	//     Color: "#444444",
	//   },
	//   "literal.string": termd.Style{
	//     Color: "#333333",
	//   },
	// }

	rendered := c.Compile(tuiCore.Help)
	fmt.Fprint(help, tview.TranslateANSI(rendered))

	tuiCore.App.SetRoot(help, true)
}

func (tuiCore *TUICore) UpdateTimeline(scrollToEnd bool) bool {
	_, _, w, _ := tuiCore.Stream.Box.GetInnerRect()

	err := tuiCore.Timeline.Load()
	if err != nil {
		// TODO: Display errors somewhere
		return false
	}

	currentTimelineType := tuiCore.Timeline.GetCurrentType()
	if tuiCore.RenderedTimelineType != currentTimelineType ||
		currentTimelineType == mast.TimelineHashtag ||
		currentTimelineType == mast.TimelineUser {
		tuiCore.Stream.Clear()
		tuiCore.RenderedTimelineType = currentTimelineType
		tuiCore.Timeline.LastRenderedIndex = -1
	}

	output, err := RenderTimeline(
		&tuiCore.Timeline,
		w,
		tuiCore.Options.ShowImages,
		tuiCore.Options.JustifyText,
	)

	if err != nil {
		// TODO: Display errors somewhere
		return false
	}

	fmt.Fprint(tuiCore.Stream, tview.TranslateANSI(output))

	if scrollToEnd == true {
		tuiCore.Stream.ScrollToEnd()
	}

	if currentTimelineType == mast.TimelineUser {
		tuiCore.Profile.Clear()

		options := tuiCore.Timeline.GetCurrentOptions()

		profileOutput, err := RenderProfile(
			&options.User,
			w,
			tuiCore.Options.ShowImages,
		)

		if err != nil {
			// TODO: Display errors somewhere
			return false
		}

		fmt.Fprint(tuiCore.Profile, tview.TranslateANSI(profileOutput))
	}

	return true
}

func (tuiCore *TUICore) EnterCommandMode() bool {
	if tuiCore.CmdLine.Box.HasFocus() == false {
		tuiCore.App.SetFocus(tuiCore.CmdLine)
		tuiCore.CmdLine.SetLabelColor(tcell.ColorTeal)

		tuiCore.Prompt = tuiCore.Timeline.Account.Username + ": "
		tuiCore.Progress.
			SetLabel(tuiCore.Prompt)

		tuiCore.Mode = CommandMode
		return true
	}

	return false
}

func (tuiCore *TUICore) ExitCommandMode(force bool) bool {
	if tuiCore.CmdLine.Box.HasFocus() == true || force == true {
		tuiCore.App.SetFocus(tuiCore.Stream)
		tuiCore.CmdLine.SetLabelColor(tcell.ColorDefault)

		tuiCore.Prompt = tuiCore.Timeline.Account.Username + "  "
		tuiCore.Progress.
			SetLabel(tuiCore.Prompt)

		tuiCore.Mode = NormalMode
		return true
	}

	return false
}

var spinners = [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// var spinners = [...]string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"}
// var spinners = [...]string{"◴", "◷", "◶", "◵"}
// var spinners = [...]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

type ProgressManager struct {
	inputField       *tview.InputField
	app              *tview.Application
	label            string
	prefix           string
	labelColor       *tcell.Color
	labelColorString string
	ProgressMutex    sync.Mutex
}

func NewProgress(field *tview.InputField, app *tview.Application) *ProgressManager {
	i := &ProgressManager{
		inputField: field,
		app:        app,
		prefix:     " ",
		labelColor: nil,
	}

	return i
}

func (i *ProgressManager) updateColor() {
	newColor, _, _ := i.inputField.GetLabelStyle().Decompose()
	if i.labelColor == nil || newColor != *i.labelColor {
		i.labelColor = &newColor
		i.labelColorString = "default"
		for key, value := range tcell.ColorNames {
			if value == newColor {
				i.labelColorString = key
				break
			}
		}
	}
}

func (i *ProgressManager) Run(action func() (mast.CmdReturnCode, bool)) {
	i.ProgressMutex.Lock()
	cmdResult := make(chan string, 1)
	i.updateColor()

	go func() {
		cmd, _ := action()
		if cmd == mast.CodeOk {
			cmdResult <- "[green]✓"
		} else if cmd == mast.CodeNotOk {
			cmdResult <- "[red]𐄂"
		} else {
			cmdResult <- " "
		}
		close(cmdResult)
		i.ProgressMutex.Unlock()
	}()

	go func() {
		var spin int

		for {
			select {
			case result, ok := <-cmdResult:
				if ok {
					i.prefix = result
					i.app.QueueUpdateDraw(func() {
						i.updateLabel()
					})
					return
				}

			case <-time.After(100 * time.Millisecond):
				i.app.QueueUpdateDraw(func() {
					i.prefix = "[yellow]" + spinners[spin%len(spinners)]
					i.updateLabel()
				})
				spin++
			}
		}
	}()
}

func (i *ProgressManager) SetLabel(label string) {
	i.label = label
	i.prefix = " "
	i.updateColor()
	i.updateLabel()
}

func (i *ProgressManager) updateLabel() {
	i.inputField.SetLabel(fmt.Sprintf("%s[%s]%s", i.prefix, i.labelColorString, i.label))
}

type History struct {
	Root  *tview.Pages
	Table *tview.Table
	modal *tview.Modal
}

func NewHistory(app *TUICore) *History {
	history := &History{}

	history.Table = tview.NewTable().
		SetSelectable(true, true).
		SetSelectedStyle(tcell.Style{}.Reverse(true)).
		SetCell(0, 0,
			tview.NewTableCell("Time").
				SetAttributes(tcell.AttrBold).
				SetSelectable(false)).
		SetCell(0, 1,
			tview.NewTableCell("Command").
				SetAttributes(tcell.AttrBold).
				SetSelectable(false)).
		SetCell(0, 2,
			tview.NewTableCell("Status").
				SetAttributes(tcell.AttrBold).
				SetAlign(tview.AlignCenter).
				SetSelectable(false)).
		SetCell(0, 3,
			tview.NewTableCell("Error").
				SetAttributes(tcell.AttrBold).
				SetSelectable(false).
				SetExpansion(1)).
		SetBorders(true).
		SetSelectedFunc(func(y, x int) {
			cell := history.Table.GetCell(y, x)
			history.modal.SetText(cell.Text)
			history.modal.SetFocus(0)
			history.Root.SendToFront("modal")
			history.Root.ShowPage("modal")
		}).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.App.SetRoot(app.Grid, true)
				app.EnterCommandMode()
			}
		})

	history.modal = tview.NewModal().
		SetButtonTextColor(tcell.ColorWhite). // this shouldnt be required but tview seems to have a bug where button background doesnt visibly update if PrimaryTextColor is ColorDefault
		AddButtons([]string{"Close", "Copy to Clipboard"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Copy to Clipboard" {
				clipboard.WriteAll(history.Table.GetCell(history.Table.GetSelection()).Text)
			} else {
				history.Root.HidePage("modal")
				history.Root.SendToBack("modal")
			}
		})

	help := tview.NewTextView().
		SetDynamicColors(true).
		SetText(" ESC to return, Enter to view cell")

	grid := tview.NewGrid().
		SetRows(0, 1).
		AddItem(history.Table, 0, 0, 1, 2, 0, 0, true).
		AddItem(help, 1, 0, 1, 1, 0, 0, false)

	history.Root = tview.NewPages().
		AddPage("background", grid, true, true).
		AddPage("modal", history.modal, true, false)

	return history
}

func (history *History) AddHistory(cmd string, code mast.CmdReturnCode, err error) {
	var color tcell.Color
	var codeString string

	switch code {
	case mast.CodeOk:
		codeString = "Success"
		color = tcell.ColorGreen
	case mast.CodeNotOk:
		codeString = "Failure"
		color = tcell.ColorRed
	case mast.CodeCommandNotFound:
		codeString = "Unknown Command"
		color = tcell.ColorRed
		break
	case mast.CodeUserNotFound:
		codeString = "User not found"
		color = tcell.ColorRed
		break
	default:
		return
	}

	dateTime := time.Now()
	dateTimeString := fmt.Sprint(dateTime.Format("01-02-2006 15:04:05"))

	history.Table.InsertRow(1)
	history.Table.SetCell(1, 0,
		tview.NewTableCell(dateTimeString))
	history.Table.SetCell(1, 1,
		tview.NewTableCell(cmd).
			SetMaxWidth(25))
	history.Table.SetCell(1, 2,
		tview.NewTableCell(codeString).
			SetAlign(tview.AlignCenter).
			SetTextColor(color))

	if err != nil {
		history.Table.SetCell(1, 3, tview.NewTableCell(err.Error()))
	} else {
		history.Table.SetCell(1, 3, tview.NewTableCell(""))
	}

	dataRows := history.Table.GetRowCount() - 1
	max := 100

	if dataRows > max {
		history.Table.RemoveRow(dataRows)
	}
}
