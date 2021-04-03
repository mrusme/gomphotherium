package tui

import (
  "fmt"

  "github.com/gdamore/tcell/v2"
  "github.com/rivo/tview"

  "github.com/mattn/go-mastodon"
  "github.com/mrusme/gomphotherium/mast"
)

type TimelineType int
const (
  TimelineHome TimelineType = 0
  TimelineLocal = 1
  TimelinePublic = 2
  TimelineNotifications = 3
  TimelineEnd = 4
)

func Timeline(MastodonClient *mastodon.Client) {
    app := tview.NewApplication()

    logoBox := tview.NewTextView().
    SetDynamicColors(true).
    SetRegions(true).
    SetWrap(true)

    box := tview.NewFlex().
    SetDirection(tview.FlexRow).
    AddItem(logoBox, 0, 1, true)

    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
      if event.Key() == tcell.KeyCtrlN {
        _, _, w, _ := logoBox.Box.GetInnerRect()
        output := mast.Timeline(MastodonClient, w)
        fmt.Fprint(logoBox, tview.TranslateANSI(output))
        return nil
      }
      return event
    })

    if err := app.SetRoot(box, true).Run(); err != nil {
      panic(err)
    }
}
