package tui

import (
  "fmt"
  "time"

  "github.com/gdamore/tcell/v2"
  "github.com/rivo/tview"

  "github.com/mattn/go-mastodon"
  "github.com/mrusme/gomphotherium/mast"
)

type ModeType int
const (
  NormalMode ModeType        = 1
  InsertMode                 = 2
)

type TUICore struct {
  App                        *tview.Application
  CmdLine                    *tview.InputField
  Stream                     *tview.TextView
  Grid                       *tview.Grid

  Prompt                     string
  Mode                       ModeType

  Timeline                   mast.Timeline
}

func TUI(tuiCore TUICore, mastodonClient *mastodon.Client) {
  tuiCore.Timeline = mast.NewTimeline(mastodonClient)
  tuiCore.App = tview.NewApplication()

  tuiCore.CmdLine = tview.NewInputField().
    SetLabelColor(tcell.ColorDefault).
    SetFieldBackgroundColor(tcell.ColorDefault).
    SetAutocompleteFunc(func(input string) ([]string) {
      return mast.CmdAutocompleter(input, tuiCore.Timeline.KnownUsers)
    }).
    SetDoneFunc(func(key tcell.Key) {
      if key == tcell.KeyEnter {
        cmd := tuiCore.CmdLine.GetText()
        tuiCore.CmdLine.SetText("")
        retCode := mast.CmdProcessor(cmd)

        if retCode == mast.CodeQuit {
          tuiCore.App.Stop();
        }
      }
    })

  tuiCore.Stream = tview.NewTextView().
    SetDynamicColors(true).
    SetRegions(true).
    SetWrap(true)

  tuiCore.Grid = tview.NewGrid().
    SetRows(0, 1).
    SetColumns(0).
    SetBorders(true).
    AddItem(tuiCore.Stream, 0, 0, 1, 1, 0, 0, false).
    AddItem(tuiCore.CmdLine, 1, 0, 1, 1, 0, 0, true)

  tuiCore.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    switch event.Key() {
    case tcell.KeyCtrlR:
      tuiCore.UpdateTimeline(true)
      return nil
    case tcell.KeyRune:
      switch event.Rune() {
      case 'i':
        if tuiCore.EnterInsertMode() == true {
          return nil
        }
      }
    case tcell.KeyEscape:
      if tuiCore.ExitInsertMode(false) == true {
        return nil
      }
    }

    return event
  })

  go func() {
    for {
      time.Sleep(time.Second * 2)
      tuiCore.UpdateTimeline(true)

      if tuiCore.Mode == 0 {
        tuiCore.ExitInsertMode(true)
      }

      tuiCore.App.Draw()
      time.Sleep(time.Second * 58)
    }
  }()

  if err := tuiCore.App.SetRoot(tuiCore.Grid, true).Run(); err != nil {
    panic(err)
  }
}

func (tuiCore *TUICore) UpdateTimeline(scrollToEnd bool) bool {
  _, _, w, _ := tuiCore.Stream.Box.GetInnerRect()

  err := tuiCore.Timeline.Load(mast.TimelineHome)
  if err != nil {
    // TODO: Display errors somewhere
    return false
  }

  output, err := RenderTimeline(&tuiCore.Timeline, w)

  if err != nil {
    // TODO: Display errors somewhere
    return false
  }

  fmt.Fprint(tuiCore.Stream, tview.TranslateANSI(output))

  if scrollToEnd == true {
    tuiCore.Stream.ScrollToEnd()
  }

  return true
}

func (tuiCore *TUICore) EnterInsertMode() bool {
  if tuiCore.CmdLine.Box.HasFocus() == false {
    tuiCore.App.SetFocus(tuiCore.CmdLine)
    tuiCore.CmdLine.SetLabelColor(tcell.ColorTeal)

    tuiCore.Prompt = tuiCore.Timeline.Account.Username + ": "
    tuiCore.CmdLine.
      SetLabel(tuiCore.Prompt)

    tuiCore.Mode = InsertMode
    return true
  }

  return false
}

func (tuiCore *TUICore) ExitInsertMode(force bool) bool {
  if tuiCore.CmdLine.Box.HasFocus() == true || force == true {
    tuiCore.App.SetFocus(tuiCore.Stream)
    tuiCore.CmdLine.SetLabelColor(tcell.ColorDefault)

    tuiCore.Prompt = tuiCore.Timeline.Account.Username
    tuiCore.CmdLine.
      SetLabel(tuiCore.Prompt)

    tuiCore.Mode = NormalMode
    return true
  }

  return false
}
