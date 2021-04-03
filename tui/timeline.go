package tui

import (
  "fmt"

  "github.com/mrusme/gomphotherium/mast"
)

func RenderTimeline(timeline *mast.Timeline, width int) (string, error) {
  var output string = ""
  var err error = nil

  var tootOutput string = ""
  for i := len(timeline.Toots) - 1; i >= 0; i-- {
    tootOutput, err = RenderToot(timeline.Toots[i], width)
    output = fmt.Sprintf("%s%s\n", output, tootOutput)
  }

  return output, err
}
