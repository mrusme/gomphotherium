package tui

import (
  "fmt"
  // "time"

  "github.com/grokify/html-strip-tags-go"
  "html"

  // "image/color"
  // "github.com/eliukblau/pixterm/pkg/ansimage"
  // "context"

  "github.com/mattn/go-mastodon"
  // "github.com/mrusme/gomphotherium/mast"
)

func RenderToot(toot *mastodon.Status, width int) (string, error) {
  var output string = ""
  var err error = nil

  createdAt := toot.CreatedAt

  account := toot.Account.Acct
  if account == "" {
    account = toot.Account.Username
  }

  inReplyTo := ""
  if toot.InReplyToID != nil {
    inReplyTo = " [magenta]\xe2\x87\x9f[-]"
  }

  output = fmt.Sprintf("%s[blue]%s[-] [grey]%s%s[-]\n", output, toot.Account.DisplayName, account, inReplyTo)
  output = fmt.Sprintf("%s%s\n", output, html.UnescapeString(strip.StripTags(toot.Content)))

  // for _, attachment := range toot.MediaAttachments {
  //   pix, err := ansimage.NewScaledFromURL(attachment.PreviewURL, int((float64(width) * 0.75)), width, color.Transparent, ansimage.ScaleModeResize, ansimage.NoDithering)
  //   if err == nil {
  //     output = fmt.Sprintf("%s\n%s\n", output, pix.RenderExt(false, false))
  //   }
  // }

  output = fmt.Sprintf("%s[magenta]\xe2\x86\xab %d[-] [green]\xe2\x86\xbb %d[-] [yellow]\xe2\x98\x85 %d[-] [grey]#%d on %s at %s[-]\n", output, toot.RepliesCount, toot.ReblogsCount, toot.FavouritesCount, createdAt.Format("Jan 2"), createdAt.Format("15:04"))

  output = fmt.Sprintf("%s\n", output)
  return output, err
}
