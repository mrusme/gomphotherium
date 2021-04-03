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
    SetLabel("@user.instance.org: ").
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
    if event.Key() == tcell.KeyCtrlN {
      _, _, w, _ := stream.Box.GetInnerRect()
      timeline.Load(mast.TimelineHome)
      output, err := RenderTimeline(&timeline, w)
      if err != nil {
        panic(err)
      }

      fmt.Fprint(stream, tview.TranslateANSI(output))
      return nil
    }
    return event
  })

  if err := app.SetRoot(grid, true).Run(); err != nil {
    panic(err)
  }
}
