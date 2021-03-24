package cli

import (
  "fmt"
  "log"
  "context"
  "image/color"

  "github.com/spf13/cobra"

  "github.com/grokify/html-strip-tags-go"
  "html"

  "github.com/eliukblau/pixterm/pkg/ansimage"
)

var timelineCmd = &cobra.Command{
  Use:   "timeline",
  Short: "Display timeline",
  Long: "Display different timelines.",
  Run: func(cmd *cobra.Command, args []string) {
    timeline, err := MastodonClient.GetTimelineHome(context.Background(), nil)
    if err != nil {
      log.Fatal(err)
    }
    for i := len(timeline) - 1; i >= 0; i-- {
      fmt.Printf("%s [%s]\n", timeline[i].Account.DisplayName, timeline[i].Account.Acct)
      fmt.Println(html.UnescapeString(strip.StripTags(timeline[i].Content)))
      for _, attachment := range timeline[i].MediaAttachments {
        pix, err := ansimage.NewScaledFromURL(attachment.PreviewURL, 30, 40, color.Transparent, ansimage.ScaleModeResize, ansimage.NoDithering)
        if err != nil {
          fmt.Println(err)
        }
        if err == nil {
          fmt.Printf("%s\n", pix.Render())
        }
      }

      fmt.Printf("\n")
    }

    return
  },
}

func init() {
  rootCmd.AddCommand(timelineCmd)
}
