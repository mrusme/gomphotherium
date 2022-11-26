package tui

import (
	"fmt"
	"image/color"

	// "time"
	// "context"

	"html"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mattn/go-runewidth"

	"github.com/mattn/go-mastodon"
	"github.com/mrusme/gomphotherium/mast"
)

func RenderToot(
	toot *mast.Toot,
	width int,
	showImages bool,
	justifyText bool) (string, error) {
	status := &toot.Status
	return RenderStatus(status, toot, width, showImages, justifyText, false)
}

func RenderStatus(
	status *mastodon.Status,
	toot *mast.Toot,
	width int,
	showImages bool,
	justifyText bool,
	isReblog bool,
) (string, error) {
	var output string = ""
	var err error = nil

	var indent string = ""
	if isReblog {
		indent = "    "
	}

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

	idPadding :=
		width -
			len(fmt.Sprint(toot.ID)) -
			runewidth.StringWidth(status.Account.DisplayName) -
			len(account) -
			// https://github.com/mattn/go-runewidth/issues/36
			runewidth.StringWidth(inReplyToOrBoost)

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

		output = fmt.Sprintf("%s%s\n",
			output,
			notificationText,
		)
	}

	if isReblog {
		output = fmt.Sprintf("%s%s[blue]%s[-] [grey]%s[-][purple]%s[-]\n",
			output,
			indent,
			status.Account.DisplayName,
			account,
			inReplyToOrBoost)
	} else {
		output = fmt.Sprintf("%s%s[blue]%s[-] [grey]%s[-][purple]%s[-][grey]%*d[-]\n",
			output,
			indent,
			status.Account.DisplayName,
			account,
			inReplyToOrBoost,
			idPadding,
			toot.ID)
	}

	if !isReblog && status.Reblog != nil {
		reblogOutput, err := RenderStatus(
			status.Reblog,
			toot,
			width,
			showImages,
			justifyText,
			true)
		if err == nil {
			output = fmt.Sprintf("%s%s", output, reblogOutput)
		}
	} else {
		var wrappedContent string = WrapWithIndent(
			html.UnescapeString(strip.StripTags(status.Content)),
			width-len(indent),
			indent,
			justifyText,
		)

		output = fmt.Sprintf("%s%s\n",
			output,
			wrappedContent)
	}

	if showImages == true {
		for _, attachment := range status.MediaAttachments {
			pix, err := ansimage.NewScaledFromURL(
				attachment.PreviewURL,
				int((float64(width) * 0.75)),
				width,
				color.Transparent,
				ansimage.ScaleModeResize,
				ansimage.NoDithering,
			)
			if err == nil {
				output = fmt.Sprintf("%s\n%s\n", output, pix.RenderExt(false, false))
			}
		}
	}

	output = fmt.Sprintf("%s%s[purple]\xe2\x86\xab %d[-] ",
		output,
		indent,
		status.RepliesCount,
	)
	output = fmt.Sprintf("%s[green]\xe2\x86\xbb %d[-] ",
		output,
		status.ReblogsCount,
	)
	output = fmt.Sprintf("%s[yellow]\xe2\x98\x85 %d[-] ",
		output,
		status.FavouritesCount,
	)
	output = fmt.Sprintf("%s[grey]on %s at %s[-]\n",
		output,
		createdAt.Format("Jan 2"),
		createdAt.Format("15:04"),
	)

	output = fmt.Sprintf("%s\n", output)
	return output, err
}
