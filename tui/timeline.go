package tui

import (
  "fmt"

  "github.com/mrusme/gomphotherium/mast"
)

func RenderTimeline(timeline *mast.Timeline, width int) (string, error) {
  var output string = ""
  var err error = nil

  var tootOutput string = ""
  newRenderedIndex := len(timeline.Toots) - 1
  for i := newRenderedIndex; i > timeline.LastRenderedIndex; i-- {
    tootOutput, err = RenderToot(&timeline.Toots[i], width)
    output = fmt.Sprintf("%s%s\n", output, tootOutput)
  }

  timeline.LastRenderedIndex = newRenderedIndex
  return output, err
}
