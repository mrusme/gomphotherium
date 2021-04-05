package tui

import (
  "fmt"
  "math"
  // "time"
  // "context"

  "github.com/mattn/go-runewidth"
  "github.com/grokify/html-strip-tags-go"
  "html"

  // "image/color"
  // "github.com/eliukblau/pixterm/pkg/ansimage"

  "github.com/mattn/go-mastodon"
  // "github.com/mrusme/gomphotherium/mast"
)

func RenderProfile(
  profile *mastodon.Account,
  width int,
  showImages bool) (string, error) {
  var output string = ""
  var err error = nil

  account := profile.Acct
  if account == "" {
    account = profile.Username
  }

  bot := ""
  if profile.Bot == true {
    bot = "\xF0\x9F\xA4\x96"
  }

  output = fmt.Sprintf("%s[blue]%s[-] [grey]%s[-] [red]%s[-]\n%s\n\n",
    output,
    profile.DisplayName,
    account,
    bot,
    runewidth.Truncate(
      html.UnescapeString(strip.StripTags(profile.Note)),
      width,
      "...",
    ),
  )

  halfwidth := int(math.Floor(float64(width)))

  fieldsNumber := len(profile.Fields)
  if fieldsNumber > 4 {
    fieldsNumber = 4
  }

  for i := 0; i < fieldsNumber; i++ {
    field := profile.Fields[i]
    output = fmt.Sprintf("%s[grey]%s:[-] [purple]%s[-]\n",
      output,
      runewidth.Truncate(
        field.Name,
        halfwidth,
        "...",
      ),
      runewidth.Truncate(
        html.UnescapeString(strip.StripTags(field.Value)),
        halfwidth,
        "...",
      ),
    )
  }

  for i := fieldsNumber; i < 4; i++ {
    output = fmt.Sprintf("%s\n", output)
  }

  output = fmt.Sprintf("%s[blue]%d[-] [grey]toots[-] ",
    output,
    profile.StatusesCount,
  )

  output = fmt.Sprintf("%s[blue]%d[-] [grey]followers[-] ",
    output,
    profile.FollowersCount,
  )

  output = fmt.Sprintf("%s[blue]%d[-] [grey]following[-] ",
    output,
    profile.FollowingCount,
  )

  return output, err
}
