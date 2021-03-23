package cli

import (
  "fmt"
  "log"
  "context"

  "github.com/spf13/cobra"
  // "github.com/mattn/go-mastodon"
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
      fmt.Println(timeline[i])
    }

    return
  },
}

func init() {
  rootCmd.AddCommand(timelineCmd)
}
