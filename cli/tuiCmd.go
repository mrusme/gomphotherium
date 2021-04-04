package cli

import (
  "github.com/spf13/cobra"
  "github.com/mrusme/gomphotherium/tui"
)

var tuiCmd = &cobra.Command{
  Use:   "tui",
  Short: "Launch TUI",
  Long: "Launch TUI.",
  Run: func(cmd *cobra.Command, args []string) {
    tuiCore := tui.TUICore{
      Client: MastodonClient,
      Options: tui.TUIOptions{
        ShowImages: flagShowImages,
      },
      Help: help,
    }
    tui.TUI(tuiCore)
  },
}

func init() {
  rootCmd.AddCommand(tuiCmd)
}
