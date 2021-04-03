package cli

import (
  // "fmt"

  "github.com/spf13/cobra"

  // "github.com/mrusme/gomphotherium/mast"
)

var timelineCmd = &cobra.Command{
  Use:   "timeline",
  Short: "Display timeline",
  Long: "Display different timelines.",
  Run: func(cmd *cobra.Command, args []string) {
    // fmt.Printf(mast.Timeline(MastodonClient, mast.TimelineHome, 40))
    return
  },
}

func init() {
  rootCmd.AddCommand(timelineCmd)
}
