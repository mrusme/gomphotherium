package tui

import (
  "fmt"
  // "time"

  "github.com/grokify/html-strip-tags-go"
  "html"

  // "image/color"
  // "github.com/eliukblau/pixterm/pkg/ansimage"
  // "context"

  // "github.com/mattn/go-mastodon"
  "github.com/mrusme/gomphotherium/mast"
)

func RenderToot(toot *mast.Toot, width int) (string, error) {
  var output string = ""
  var err error = nil

  status := &toot.Status

  createdAt := status.CreatedAt

  account := status.Account.Acct
  if account == "" {
    account = status.Account.Username
  }

  inReplyTo := ""
  inReplyToLen := 0
  if status.InReplyToID != nil {
    inReplyTo = " \xe2\x87\x9f"
    inReplyToLen = 1
  }

  output = fmt.Sprintf("%s[blue]%s[-] [grey]%s[-][magenta]%s[-][grey]%*d[-]\n", output, status.Account.DisplayName, account, inReplyTo, (width - len(string(toot.ID)) - len(status.Account.DisplayName) - len(account) - inReplyToLen), toot.ID)
  output = fmt.Sprintf("%s%s\n", output, html.UnescapeString(strip.StripTags(status.Content)))

  // for _, attachment := range status.MediaAttachments {
  //   pix, err := ansimage.NewScaledFromURL(attachment.PreviewURL, int((float64(width) * 0.75)), width, color.Transparent, ansimage.ScaleModeResize, ansimage.NoDithering)
  //   if err == nil {
  //     output = fmt.Sprintf("%s\n%s\n", output, pix.RenderExt(false, false))
  //   }
  // }

  output = fmt.Sprintf("%s[magenta]\xe2\x86\xab %d[-] [green]\xe2\x86\xbb %d[-] [yellow]\xe2\x98\x85 %d[-] [grey]on %s at %s[-]\n", output, status.RepliesCount, status.ReblogsCount, status.FavouritesCount, createdAt.Format("Jan 2"), createdAt.Format("15:04"))

  output = fmt.Sprintf("%s\n", output)
  return output, err
}
