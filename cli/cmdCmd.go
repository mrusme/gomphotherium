package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"

	"github.com/mrusme/gomphotherium/mast"
	"github.com/mrusme/gomphotherium/tui"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Run command",
	Long:  "Run command directly from the command line.",
	Run: func(cmd *cobra.Command, args []string) {
		timeline := mast.NewTimeline(MastodonClient)
		result := mast.CmdProcessor(
			&timeline,
			strings.Join(args, " "),
			mast.TriggerCLI,
		)

		cmdReturn, _, loadTimeline := result.Decompose()
		switch cmdReturn {
		case mast.CodeOk:
			if loadTimeline == true {
				timeline.Load()
				output, err := tui.RenderTimeline(&timeline, 72, flagShowImages, flagJustifyText)
				if err != nil {
					panic(err)
				}

				fmt.Printf(tview.TranslateANSI(output))
			}
		case mast.CodeTriggerNotSupported:
			fmt.Printf("Command not supported from CLI!\n")
			os.Exit(-1)
		default:
			fmt.Printf("%v\n", cmdReturn)
			os.Exit(-1)
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(cmdCmd)
}
