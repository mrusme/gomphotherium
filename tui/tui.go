package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/mattn/go-mastodon"
	"github.com/mrusme/gomphotherium/mast"

	"github.com/tj/go-termd"
)

type ModeType int

const (
	NormalMode  ModeType = 1
	CommandMode          = 2
)

type TUIOptions struct {
	ShowImages     bool
	AutoCompletion bool
}

type TUICore struct {
	Client  *mastodon.Client
	App     *tview.Application
	CmdLine *tview.InputField
	Profile *tview.TextView
	Stream  *tview.TextView
	Grid    *tview.Grid

	Prompt string
	Mode   ModeType

	Timeline             mast.Timeline
	RenderedTimelineType mast.TimelineType

	Options TUIOptions

	Help string
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

	tuiCore.CmdLine = tview.NewInputField().
		SetLabelColor(tcell.ColorDefault).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				cmd := tuiCore.CmdLine.GetText()
				tuiCore.CmdLine.SetText("")
				retCode, _ := mast.CmdProcessor(&tuiCore.Timeline, cmd, mast.TriggerTUI)

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
				case mast.CodeQuit:
					tuiCore.App.Stop()
				}
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
		for {
			time.Sleep(time.Second * 2)
			tuiCore.UpdateTimeline(true)

			if tuiCore.Mode == 0 {
				tuiCore.ExitCommandMode(true)
			}

			tuiCore.App.Draw()
			time.Sleep(time.Second * 58)
		}
	}()

	if err := tuiCore.App.SetRoot(tuiCore.Grid, true).Run(); err != nil {
		panic(err)
	}
}

func (tuiCore *TUICore) ShowHelp() {
	var c termd.Compiler

	help := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetDoneFunc(func(key tcell.Key) {
			tuiCore.App.SetRoot(tuiCore.Grid, true)
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
		tuiCore.CmdLine.
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
		tuiCore.CmdLine.
			SetLabel(tuiCore.Prompt)

		tuiCore.Mode = NormalMode
		return true
	}

	return false
}
