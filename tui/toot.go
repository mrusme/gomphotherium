package tui

import (
  "fmt"

  "github.com/grokify/html-strip-tags-go"
  "html"

  "image/color"
  "github.com/eliukblau/pixterm/pkg/ansimage"
  // "context"

  "github.com/mattn/go-mastodon"
  // "github.com/mrusme/gomphotherium/mast"
)

func RenderToot(toot *mastodon.Status, width int) (string, error) {
  var output string = ""
  var err error = nil

  output = fmt.Sprintf("%s%s [%s]\n", output, toot.Account.DisplayName, toot.Account.Acct)
  output = fmt.Sprintf("%s%s\n", output, html.UnescapeString(strip.StripTags(toot.Content)))

  for _, attachment := range toot.MediaAttachments {
    pix, err := ansimage.NewScaledFromURL(attachment.PreviewURL, int((float64(width) * 0.75)), width, color.Transparent, ansimage.ScaleModeResize, ansimage.NoDithering)
    if err == nil {
      output = fmt.Sprintf("%s%s\n", output, pix.RenderExt(false, false))
    }
  }

  output = fmt.Sprintf("%s\n", output)
  return output, err
}
