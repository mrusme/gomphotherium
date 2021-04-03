package tui

import (
  "fmt"

  "github.com/gdamore/tcell/v2"
  "github.com/rivo/tview"

  "github.com/mattn/go-mastodon"
  "github.com/mrusme/gomphotherium/mast"
)

func TUI(mastodonClient *mastodon.Client) {
  timeline := mast.NewTimeline(mastodonClient)
  app := tview.NewApplication()

  input := tview.NewInputField().
    SetLabelColor(tcell.ColorDefault).
    SetFieldBackgroundColor(tcell.ColorDefault).
    SetDoneFunc(func(key tcell.Key) {
      // app.Stop()
    })

  stream := tview.NewTextView().
    SetDynamicColors(true).
    SetRegions(true).
    SetWrap(true)

  grid := tview.NewGrid().
    SetRows(0, 1).
    SetColumns(0).
    SetBorders(true).
    AddItem(stream, 0, 0, 1, 1, 0, 0, false).
    AddItem(input, 1, 0, 1, 1, 0, 0, true)

  app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    switch event.Key() {
    case tcell.KeyCtrlR:
      _, _, w, _ := stream.Box.GetInnerRect()
      timeline.Load(mast.TimelineHome)
      output, err := RenderTimeline(&timeline, w)
      if err != nil {
        panic(err)
      }

      input.
        SetLabel(timeline.Account.Username + ": ").
        SetLabelColor(tcell.ColorTeal)
      app.SetFocus(input)

      fmt.Fprint(stream, tview.TranslateANSI(output))

      stream.ScrollToEnd()
      return nil
    case tcell.KeyRune:
      switch event.Rune() {
      case 'i':
        if input.Box.HasFocus() == false {
          app.SetFocus(input)
          input.SetLabelColor(tcell.ColorTeal)
          return nil
        }
      }
    case tcell.KeyEscape:
      if input.Box.HasFocus() == true {
        app.SetFocus(stream)
        input.SetLabelColor(tcell.ColorDefault)
        return nil
      }
    // case tcell.KeyPgDn:
    //   app.SetFocus(stream)
    }

    return event
  })

  if err := app.SetRoot(grid, true).Run(); err != nil {
    panic(err)
  }
}
