package cli

import (
	"github.com/mrusme/gomphotherium/tui"
	"github.com/spf13/cobra"
)

var flagAutocompletion bool

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch TUI",
	Long:  "Launch TUI.",
	Run: func(cmd *cobra.Command, args []string) {
		tuiCore := tui.TUICore{
			Client: MastodonClient,
			Options: tui.TUIOptions{
				ShowImages:     flagShowImages,
				ShowUserImages: flagShowUserImages,
				TempDir:        tempDir,
				AutoCompletion: flagAutocompletion,
				JustifyText:    flagJustifyText,
			},
			Help: help,
		}
		tui.TUI(tuiCore)
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	tuiCmd.Flags().BoolVarP(
		&flagAutocompletion,
		"auto-completion",
		"a",
		true,
		"Auto-completion on command input",
	)
}
