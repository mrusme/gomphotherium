package cli

import (
  "fmt"
  "log"
  "context"
  "bufio"
  "strings"
  "os"

  "github.com/spf13/cobra"
  "github.com/mattn/go-mastodon"
)

var authenticateCmd = &cobra.Command{
  Use:   "authenticate",
  Short: "Authenticate",
  Long: "Authenticate against a Mastodon instance.",
  Args: cobra.MinimumNArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    server := args[0]

    mastodonApp, err := mastodon.RegisterApp(
      context.Background(),
      &mastodon.AppConfig{
        Server:     server,
        ClientName: "Gomphotherium",
        Scopes:     "read write follow",
        Website:    "https://github.com/mrusme/gomphotherium",
      },
    )
    if err != nil {
      log.Fatal(err)
    }

    fmt.Printf("Please open the following URL to authenticate: %s\n\n",
      mastodonApp.AuthURI)
    fmt.Printf("Paste authentication code here and press enter: ")

    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println("An error occured while reading input. Please try again",
        err)
      os.Exit(-1)
    }
    authCode := strings.TrimSuffix(input, "\n")

    mastodonClient := mastodon.NewClient(&mastodon.Config{
      Server:       server,
      ClientID:     mastodonApp.ClientID,
      ClientSecret: mastodonApp.ClientSecret,
    })

    err = mastodonClient.AuthenticateToken(
      context.Background(),
      authCode,
      "urn:ietf:wg:oauth:2.0:oob",
    )
    if err != nil {
      log.Fatal(err)
      os.Exit(-1)
    }

    fmt.Printf("Success!\n")
    fmt.Printf(
      "Please either export the following variables to your session:\n\n")

    fmt.Printf("export GOMPHOTHERIUM_SERVER='%s'",
      server)
    fmt.Printf("export GOMPHOTHERIUM_ACCESS_TOKEN='%s'",
      mastodonClient.Config.AccessToken)

    fmt.Printf("\n... or call gomphotherium with the following flags:\n\n")

    fmt.Printf("gomphotherium --server '%s' --access-token '%s' ...\n\n",
      server,
      mastodonClient.Config.AccessToken)

    return
  },
}

func init() {
  rootCmd.AddCommand(authenticateCmd)
}
