package mast

import (
  "fmt"
  "log"
  "context"

  "github.com/grokify/html-strip-tags-go"
  "html"

  "image/color"
  "github.com/eliukblau/pixterm/pkg/ansimage"

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

func Timeline(mastodonClient *mastodon.Client, timelineType TimelineType, width int) string {
  var output string = ""
  var timeline []*mastodon.Status
  var err error

  switch timelineType {
  case TimelineHome:
    timeline, err = mastodonClient.GetTimelineHome(context.Background(), nil)
  case TimelineLocal:
    timeline, err = mastodonClient.GetTimelinePublic(context.Background(), true, nil)
  case TimelinePublic:
    timeline, err = mastodonClient.GetTimelinePublic(context.Background(), false, nil)
  case TimelineNotifications:
    notifications, err := mastodonClient.GetNotifications(context.Background(), nil)
    if err != nil {
      log.Fatal(err) // TODO
    }

    for _, notification := range notifications {
      timeline = append(timeline, notification.Status)
    }
  }

  if err != nil {
    log.Fatal(err) // TODO
  }

  for i := len(timeline) - 1; i >= 0; i-- {
    output = fmt.Sprintf("%s%s [%s]\n", output, timeline[i].Account.DisplayName, timeline[i].Account.Acct)
    output = fmt.Sprintf("%s%s\n", output, html.UnescapeString(strip.StripTags(timeline[i].Content)))
    for _, attachment := range timeline[i].MediaAttachments {
      pix, err := ansimage.NewScaledFromURL(attachment.PreviewURL, int((float64(width) * 0.75)), width, color.Transparent, ansimage.ScaleModeResize, ansimage.NoDithering)
      if err != nil {
        fmt.Println(err)
      }
      if err == nil {
        output = fmt.Sprintf("%s%s\n", output, pix.RenderExt(false, false))
      }
    }

    output = fmt.Sprintf("%s\n", output)
  }

  return output
}
