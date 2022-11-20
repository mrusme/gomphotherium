package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VERSION string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  `The version of gomphotherium.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gomphotherium", VERSION)
	},
}
