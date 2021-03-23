package cli

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/mattn/go-mastodon"
)

var server string
var accessToken string
// var clientID string
// var clientSecret string

// var MastodonApp *mastodon.Application
var MastodonClient *mastodon.Client

var rootCmd = &cobra.Command{
  Use:   "gomphotherium",
  Short: "Command line Mastodon client",
  Long:  `A command line client for Mastodon.`,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Printf("%+v\n", err)
    os.Exit(-1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  rootCmd.PersistentFlags().StringVar(&server, "server", LookupStrEnv("GOMPHOTHERIUM_SERVER",""), "Mastodon server")
  rootCmd.PersistentFlags().StringVar(&accessToken, "access-token", LookupStrEnv("GOMPHOTHERIUM_ACCESS_TOKEN",""), "Mastodon access token")
}

func initConfig() {
  MastodonClient = mastodon.NewClient(&mastodon.Config{
    Server:       server,
    AccessToken:  accessToken,
  })
}
