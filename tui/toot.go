package tui

import (

	// "strings"
	// "fmt"
	// "github.com/patrickmn/go-cache"
	// "time"
	// "context"

	"fmt"
	"strings"

	// "time"
	// "context"

	"html"

	//   "github.com/mattn/go-mastodon"
	//   "github.com/mrusme/gomphotherium/mast"
	// )

	strip "github.com/grokify/html-strip-tags-go"

	"github.com/mattn/go-mastodon"
	"github.com/mrusme/gomphotherium/mast"
)

func RenderToot(
	toot *mast.Toot,
	imageCache *Images,
	width int,
	showImages bool,
	showUserImages bool,
	justifyText bool) (string, error) {

	var indent Indent

	if showUserImages {
		userImage := LoadImage(imageCache, width, toot)

		indent.InitializeWithArray(len((*userImage)[0]), *userImage)
	} else {
		indentPadding := 6
		indentStrings := make([]string, 2)
		indentStrings[0] = fmt.Sprintf("[grey]%*d[-] │ ", indentPadding, toot.ID)
		indentStrings[1] = fmt.Sprintf("%*s │ ", indentPadding, " ")
		indent.InitializeWithArray(len(indentStrings[0]), *&indentStrings)
	}

	status := &toot.Status
	lines, err := RenderStatus(status, toot, imageCache, width-indent.Width, showImages, justifyText, false)
	if err == nil {
		newLines := append(*lines, "")
		lines = indent.IndentSlice(&newLines)
		output := ""
		for i, line := range *lines {
			output += line
			if i != len(*lines)-1 {
				output += "\n"
			}
		}
		return output, nil
	} else {
		return "", err
	}
}

func RenderStatus(
	status *mastodon.Status,
	toot *mast.Toot,
	imageCache *Images,
	width int,
	showImages bool,
	justifyText bool,
	isReblog bool,
) (*[]string, error) {
	var err error = nil

	var output []string = make([]string, 0)

	createdAt := status.CreatedAt

	account := status.Account.Acct
	if account == "" {
		account = status.Account.Username
	}

	inReplyToOrBoost := ""
	if status.InReplyToID != nil {
		inReplyToOrBoost = " \xe2\x87\x9f"
	} else if status.Reblog != nil {
		inReplyToOrBoost = " \xe2\x86\xab"
	}

	if !isReblog && toot.IsNotification == true {
		notification := &toot.Notification

		notificationText := ""

		notificationAccount := notification.Account.Acct
		if notificationAccount == "" {
			notificationAccount = notification.Account.Username
		}

		// https://docs.joinmastodon.org/entities/notification/#type
		switch notification.Type {
		case "follow":
			notificationText = fmt.Sprintf("[red]\xe2\x98\xba\xef\xb8\x8e %s followed you[-]",
				notificationAccount,
			)
		case "follow_request":
			notificationText = fmt.Sprintf("[blue]\xe2\x98\x95 %s requested to follow you[-]",
				notificationAccount,
			)
		case "mention":
			notificationText = fmt.Sprintf("[purple]\xe2\x86\xab %s mentioned you[-]",
				notificationAccount,
			)
		case "reblog":
			notificationText = fmt.Sprintf("[green]\xe2\x86\xbb %s boosted your toot[-]",
				notificationAccount,
			)
		case "favourite":
			notificationText = fmt.Sprintf("[yellow]\xe2\x98\x85 %s faved your toot[-]",
				notificationAccount,
			)
		case "poll":
			notificationText = fmt.Sprintf("[grey]\xe2\x9c\x8e A poll by %s has ended[-]",
				notificationAccount,
			)
		case "status":
			notificationText = fmt.Sprintf("[grey]\xe2\x9c\x8c\xef\xb8\x8e %s posted a toot[-]",
				notificationAccount,
			)
		}

		output = append(output,
			notificationText,
		)
	}

	output = append(output, fmt.Sprintf("[blue]%s[-] [grey]%s[-][purple]%s[-]",
		status.Account.DisplayName,
		account,
		inReplyToOrBoost))

	if !isReblog && status.Reblog != nil {
		var indent Indent
		indent.InitializeWithString(4, "    ")

		reblogOutput, err := RenderStatus(
			status.Reblog,
			toot,
			imageCache,
			width-indent.Width,
			showImages,
			justifyText,
			true)
		if err == nil {
			output = append(output, *indent.IndentSlice(reblogOutput)...)
		}
	} else {
		lines := WrapWithIndent(
			html.UnescapeString(strip.StripTags(status.Content)),
			width,
			justifyText,
		)

		output = append(output, *lines...)
	}

	if showImages == true {
		for _, attachment := range status.MediaAttachments {
			image := imageCache.ImageAtSize(
				attachment.PreviewURL,
				int((float64(width) * 0.75)),
				width,
				nil)

			output = append(output, *image...)
		}
	}

	finalLine := fmt.Sprintf("[purple]\xe2\x86\xab %d[-] ",
		status.RepliesCount,
	)
	finalLine += fmt.Sprintf("[green]\xe2\x86\xbb %d[-] ",
		status.ReblogsCount,
	)
	finalLine += fmt.Sprintf("[yellow]\xe2\x98\x85 %d[-] ",
		status.FavouritesCount,
	)
	finalLine += fmt.Sprintf("[grey]on %s at %s[-]",
		createdAt.Format("Jan 2"),
		createdAt.Format("15:04"),
	)
	output = append(output, finalLine)

	return &output, err
}

func LoadImage(imageCache *Images, width int, toot *mast.Toot) *[]string {
	imageWidth := 6
	if width > 100 && width <= 150 {
		imageWidth = 10
	} else if width > 150 && width <= 200 {
		imageWidth = 14
	} else if width > 200 {
		imageWidth = 18
	}

	avatarUrl := toot.Status.Account.Avatar
	if toot.Status.Reblog != nil {
		avatarUrl = toot.Status.Reblog.Account.Avatar
	}

	image := imageCache.ImageAtSize(avatarUrl, imageWidth+2, imageWidth, func(loaded *[]string) *[]string {
		splitImage := make([]string, 1)
		splitImage[0] = fmt.Sprintf(" [grey]%*d[-] │ ", imageWidth, toot.ID)
		splitImage = append(splitImage, *loaded...)
		if splitImage[len(splitImage)-1] == "" {
			splitImage = splitImage[:len(splitImage)-1]
		}
		splitImage = append(splitImage, strings.Repeat(" ", imageWidth))
		for i, row := range splitImage {
			if i != 0 {
				splitImage[i] = " " + row + " │ "
			}
		}
		return &splitImage
	})

	return image
}
