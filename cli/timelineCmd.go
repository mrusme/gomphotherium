package cli

import (
  "fmt"

  "github.com/spf13/cobra"
  "github.com/rivo/tview"

  "github.com/mrusme/gomphotherium/mast"
  "github.com/mrusme/gomphotherium/tui"
)

var timelineCmd = &cobra.Command{
  Use:   "timeline",
  Short: "Display timeline",
  Long: "Display different timelines.",
  Run: func(cmd *cobra.Command, args []string) {
    timeline := mast.NewTimeline(MastodonClient)
    timeline.Switch(mast.TimelineHome, nil)
    timeline.Load()
    output, err := tui.RenderTimeline(&timeline, 72, flagShowImages)
    if err != nil {
      panic(err)
    }

    fmt.Printf(tview.TranslateANSI(output))
    return
  },
}

func init() {
  rootCmd.AddCommand(timelineCmd)
}
