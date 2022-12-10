package cli

import (
	"fmt"
	"os"

	"github.com/mattn/go-mastodon"
	"github.com/spf13/cobra"
)

var help string
var server string
var accessToken string
var tempDir string
var flagShowImages bool
var flagShowUserImages bool
var flagJustifyText bool

// var clientID string
// var clientSecret string

// var MastodonApp *mastodon.Application
var MastodonClient *mastodon.Client

var rootCmd = &cobra.Command{
	Use:   "gomphotherium",
	Short: "Command line Mastodon client",
	Long:  `A command line client for Mastodon.`,
}

func Execute(embeddedHelp string) {
	help = embeddedHelp
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&server,
		"server",
		LookupStrEnv("GOMPHOTHERIUM_SERVER", ""),
		"Mastodon server",
	)
	rootCmd.PersistentFlags().StringVar(
		&accessToken,
		"access-token",
		LookupStrEnv("GOMPHOTHERIUM_ACCESS_TOKEN", ""),
		"Mastodon access token",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&flagShowImages,
		"show-images",
		"i",
		true,
		"Show images in timeline",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&flagShowUserImages,
		"show-user-images",
		"u",
		true,
		"Show user images",
	)
	rootCmd.PersistentFlags().StringVar(
		&tempDir,
		"temp-dir",
		"",
		"Temporary director for image cache.  Default is the system default temporary directory",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&flagJustifyText,
		"justify-text",
		"j",
		true,
		"Justify text in timeline",
	)

	if tempDir == "" {
		tempDir = os.TempDir()
	}
}

func initConfig() {
	MastodonClient = mastodon.NewClient(&mastodon.Config{
		Server:      server,
		AccessToken: accessToken,
	})
}
