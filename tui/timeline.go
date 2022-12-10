package tui

import (
	"fmt"

	"github.com/mrusme/gomphotherium/mast"
)

func RenderTimeline(
	timeline *mast.Timeline,
	imageCache *Images,
	width int,
	showImages bool,
	showUserImages bool,
	justifyText bool) (string, error) {
	var output string = ""
	var err error = nil

	var tootOutput string = ""
	newRenderedIndex := len(timeline.Toots)
	for i := (timeline.LastRenderedIndex + 1); i < newRenderedIndex; i++ {
		if i == timeline.LastRenderedIndex+1 {
			output += "\n"
		}
		tootOutput, err = RenderToot(&timeline.Toots[i], imageCache, width, showImages, showUserImages, justifyText)
		output = fmt.Sprintf("%s%s", output, tootOutput)
		if i != newRenderedIndex-1 {
			output += "\n"
		}
	}

	timeline.LastRenderedIndex = (newRenderedIndex - 1)
	return output, err
}
