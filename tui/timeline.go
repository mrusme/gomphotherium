package tui

import (
  "log"
  "fmt"
  "context"
  "image/color"

  "github.com/grokify/html-strip-tags-go"
  "html"

  "github.com/eliukblau/pixterm/pkg/ansimage"

  ui "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"

  "github.com/mattn/go-mastodon"
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
  if err := ui.Init(); err != nil {
    log.Fatalf("failed to initialize termui: %v", err)
  }
  defer ui.Close()

  termWidth, termHeight := ui.TerminalDimensions()

  // timelines :=

  var rows []string

  timeline, err := MastodonClient.GetTimelineHome(context.Background(), nil)
  if err != nil {
    log.Fatal(err)
  }
  for i := len(timeline) - 1; i >= 0; i-- {
    toot := fmt.Sprintf("%s [%s]\n", timeline[i].Account.DisplayName, timeline[i].Account.Acct)
    toot = fmt.Sprintf("%s%s\n", toot, html.UnescapeString(strip.StripTags(timeline[i].Content)))
    for _, attachment := range timeline[i].MediaAttachments {
      pix, err := ansimage.NewScaledFromURL(attachment.PreviewURL, 30, 40, color.Transparent, ansimage.ScaleModeResize, ansimage.NoDithering)
      if err != nil {
        fmt.Println(err)
      }
      if err == nil {
        toot = fmt.Sprintf("%s%s\n", toot, pix.Render())
      }
    }

    rows = append(rows, toot)
  }


  l := widgets.NewList()
  l.Title = "Home"
  l.Rows = rows
  l.TextStyle = ui.NewStyle(ui.ColorYellow)
  l.WrapText = false
  l.SetRect(0, 5, termWidth, termHeight - 4)

  tabpane := widgets.NewTabPane("home", "local", "public", "notifications")
  tabpane.SetRect(0, 1, termWidth, 4)
  tabpane.Border = true

  renderTab := func() {
    switch tabpane.ActiveTabIndex {
    case 0:
      ui.Render(l)
    case 1:
      ui.Render(l)
    case 2:
      ui.Render(l)
    }
  }


  ui.Render(tabpane, l)

  uiEvents := ui.PollEvents()

  for {
    e := <-uiEvents
    switch e.ID {
    case "q", "<C-c>":
      return
    case "<Resize>":
      payload := e.Payload.(ui.Resize)
      termWidth = payload.Width
      termHeight = payload.Height
      ui.Clear()
      ui.Render(tabpane)
    case "h":
      tabpane.FocusLeft()
      ui.Clear()
      ui.Render(tabpane)
      renderTab()
    case "l":
      tabpane.FocusRight()
      ui.Clear()
      ui.Render(tabpane)
      renderTab()
    case "j", "<Down>":
      l.ScrollDown()
      ui.Render(l)
    case "k", "<Up>":
      l.ScrollUp()
      ui.Render(l)
    }
  }
}
