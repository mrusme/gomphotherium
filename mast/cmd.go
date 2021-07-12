package mast

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/mattn/go-mastodon"
)

type CmdReturnCode int

const (
	CodeOk                  CmdReturnCode = 0
	CodeNotOk                             = 1
	CodeCommandNotFound                   = 2
	CodeUserNotFound                      = 3
	CodeTriggerNotSupported               = 4

	CodeQuit = -1
	CodeHelp = -2
)

type CmdTrigger int

const (
	TriggerCLI = 0
	TriggerTUI = 1
)

var CmdContentRegex = regexp.MustCompile(`(?m)(( {0,1}~#| {0,1}~:)\[([^\[\]]*)\]| {0,1}~!!)`)
var CmdHandleAutoCompletionRegex = regexp.MustCompile(`(?m)(^@| @|^whois @{0,1})([^ ]+)$`)

func CmdAvailable() []string {
	return []string{
		"home",
		"local",
		"public",
		"notifications",
		"hashtag",

		"whois",

		"t",
		"toot",

		"tp",
		"tootprivate",

		"tu",
		"tootunlisted",

		"td",
		"tootdirect",

		"re",
		"reply",

		"rep",
		"replyprivate",

		"reu",
		"replyunlisted",

		"red",
		"replydirect",

		"rt",
		"retoot",
		"boost",

		"ut",
		"unretoot",
		"unboost",

		"fav",
		"unfav",

		"open",
		"share",

		// "search",

		"help",
		"?",

		"quit",
		"exit",
		"bye",
	}
}

func CmdAutocompleter(input string, knownUsers map[string]string) []string {
	var entries []string

	if input == "" {
		return entries
	}

	handles := CmdHandleAutoCompletionRegex.FindAllStringSubmatch(input, -1)
	if len(handles) > 0 {
		handle := handles[0][2]

		for _, knownUser := range knownUsers {
			if strings.HasPrefix(knownUser, handle) == true {
				entries = append(entries, input+knownUser[len(handle):])
			}
		}

		return entries
	}

	for _, cmd := range CmdAvailable() {
		if strings.HasPrefix(cmd, input) == true {
			entries = append(entries, cmd)
		}
	}

	return entries
}

func CmdProcessor(timeline *Timeline, input string, trigger CmdTrigger) (CmdReturnCode, bool) {
	split := strings.SplitN(input, " ", 2)
	cmd := split[0]

	args := ""
	if len(split) == 2 {
		args = split[1]
	}

	switch cmd {
	case "home":
		timeline.Switch(TimelineHome, nil)
		return CodeOk, true
	case "local":
		timeline.Switch(TimelineLocal, nil)
		return CodeOk, true
	case "public":
		timeline.Switch(TimelinePublic, nil)
		return CodeOk, true
	case "notifications":
		timeline.Switch(TimelineNotifications, nil)
		return CodeOk, true
	case "hashtag":
		hashtag, isLocal, err := CmdHelperGetHashtagParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		timelineOptions := TimelineOptions{
			Hashtag: hashtag,
			IsLocal: isLocal,
		}

		timeline.Switch(TimelineHashtag, &timelineOptions)
		return CodeOk, true
	case "whois":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		var account *mastodon.Account
		var err error

		accountID, accountKnown := timeline.KnownUsers[args]
		if accountKnown == true {
			account, err = timeline.User(accountID)
		}

		if accountKnown == false || err != nil {
			accounts, err := timeline.SearchUser(args, 1)
			if err != nil || len(accounts) < 1 {
				// TODO: pass info back to caller
				return CodeUserNotFound, false
			}

			account = accounts[0]
		}

		timelineOptions := TimelineOptions{
			User: *account,
		}

		timeline.Switch(TimelineUser, &timelineOptions)
		return CodeOk, true
	case "t", "toot":
		return CmdToot(timeline, args, -1, VisibilityPublic), false
	case "tp", "tootprivate":
		return CmdToot(timeline, args, -1, VisibilityPrivate), false
	case "tu", "tootunlisted":
		return CmdToot(timeline, args, -1, VisibilityUnlisted), false
	case "td", "tootdirect":
		return CmdToot(timeline, args, -1, VisibilityUnlisted), false
	case "re", "reply":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, args, err := CmdHelperGetReplyParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdToot(timeline, args, tootId, VisibilityPublic), false
	case "rep", "replyprivate":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, args, err := CmdHelperGetReplyParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdToot(timeline, args, tootId, VisibilityPrivate), false
	case "reu", "replyunlisted":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, args, err := CmdHelperGetReplyParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdToot(timeline, args, tootId, VisibilityUnlisted), false
	case "red", "replydirect":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, args, err := CmdHelperGetReplyParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdToot(timeline, args, tootId, VisibilityDirect), false
	case "rt", "retoot", "boost":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, err := CmdHelperGetBoostParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdBoost(timeline, tootId), false
	case "ut", "unretoot", "unboost":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, err := CmdHelperGetBoostParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdUnboost(timeline, tootId), false
	case "fav":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, err := CmdHelperGetFavParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdFav(timeline, tootId), false
	case "unfav":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, err := CmdHelperGetFavParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdUnfav(timeline, tootId), false
	case "open":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, err := CmdHelperGetOpenParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdOpen(timeline, tootId), false
	case "share":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		tootId, err := CmdHelperGetShareParams(args)
		if err != nil {
			return CodeNotOk, false
		}

		return CmdShare(timeline, tootId), false
	case "?", "help":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		return CodeHelp, false
	case "quit", "exit", "bye", "q":
		if trigger != TriggerTUI {
			return CodeTriggerNotSupported, false
		}

		return CodeQuit, false
	}

	return CodeCommandNotFound, false
}

func CmdHelperGetHashtagParams(args string) (string, bool, error) {
	splitArgs := strings.SplitN(args, " ", 2)

	if len(splitArgs) < 2 {
		return args, false, nil
	}

	hashtag := splitArgs[0]
	isLocal := false

	if strings.ToLower(splitArgs[1]) == "local" {
		isLocal = true
	}

	return hashtag, isLocal, nil
}

func CmdHelperGetTootIDFromString(s string) (int, error) {
	tootId, err := strconv.Atoi(s)
	if err != nil {
		return -1, errors.New("Toot ID invalid!")
	}

	return tootId, nil
}

func CmdHelperGetReplyParams(args string) (int, string, error) {
	splitArgs := strings.SplitN(args, " ", 2)

	if len(splitArgs) < 2 {
		return -1, args, errors.New("Toot ID missing!")
	}

	tootId, err := CmdHelperGetTootIDFromString(splitArgs[0])
	if err != nil {
		return tootId, args, err
	}

	newArgs := splitArgs[1]

	return tootId, newArgs, nil
}

func CmdHelperGetBoostParams(args string) (int, error) {
	return CmdHelperGetTootIDFromString(args)
}

func CmdHelperGetFavParams(args string) (int, error) {
	return CmdHelperGetTootIDFromString(args)
}

func CmdHelperGetOpenParams(args string) (int, error) {
	return CmdHelperGetTootIDFromString(args)
}

func CmdHelperGetShareParams(args string) (int, error) {
	return CmdHelperGetTootIDFromString(args)
}

func CmdToot(
	timeline *Timeline,
	content string,
	inReplyTo int,
	visibility string) CmdReturnCode {
	var status string = ""
	var spoiler string = ""
	var sensitive bool = false
	var filesToUpload []string

	tokens := CmdContentRegex.FindAllStringSubmatch(content, -1)
	for _, token := range tokens {
		if len(token[0]) >= 3 && token[0][(len(token[0])-3):] == "~!!" {
			sensitive = true
			continue
		} else if len(token) == 4 && len(token[2]) >= 2 {
			switch token[2][(len(tokens[0][2]) - 2):] {
			case "~#":
				spoiler = token[3]
				continue
			case "~:":
				filesToUpload = append(filesToUpload, token[3])
				continue
			}
		}
	}

	status = CmdContentRegex.ReplaceAllString(content, "")

	_, err := timeline.Toot(
		&status,
		inReplyTo,
		filesToUpload,
		&visibility,
		sensitive,
		&spoiler,
	)
	if err != nil {
		return CodeNotOk
	}

	return CodeOk
}

func CmdBoost(timeline *Timeline, tootID int) CmdReturnCode {
	_, err := timeline.Boost(tootID, true)
	if err != nil {
		return CodeNotOk
	}

	return CodeOk
}

func CmdUnboost(timeline *Timeline, tootID int) CmdReturnCode {
	_, err := timeline.Boost(tootID, false)
	if err != nil {
		return CodeNotOk
	}

	return CodeOk
}

func CmdFav(timeline *Timeline, tootID int) CmdReturnCode {
	_, err := timeline.Fav(tootID, true)
	if err != nil {
		return CodeNotOk
	}

	return CodeOk
}

func CmdUnfav(timeline *Timeline, tootID int) CmdReturnCode {
	_, err := timeline.Fav(tootID, false)
	if err != nil {
		return CodeNotOk
	}

	return CodeOk
}

func CmdOpen(timeline *Timeline, tootID int) CmdReturnCode {
	var cmd *exec.Cmd

	url := timeline.Toots[tootID].Status.URL

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		// errors.New("Platform not supported!")
		return CodeNotOk
	}

	cmd.Env = append(os.Environ())
	err := cmd.Start()
	if err != nil {
		return CodeNotOk
	}

	return CodeOk
}

func CmdShare(timeline *Timeline, tootID int) CmdReturnCode {
	url := timeline.Toots[tootID].Status.URL

	err := clipboard.WriteAll(url)
	if err != nil {
		return CodeNotOk
	}

	return CodeOk
}
