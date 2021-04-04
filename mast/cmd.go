package mast

import (
  "strings"
  "strconv"
  "regexp"
  "errors"
  "os/exec"
  "runtime"

  "github.com/atotto/clipboard"
)

type CmdReturnCode int
const (
  CodeOk CmdReturnCode       = 0
  CodeNotOk                  = 1
  CodeCommandNotFound        = 2
  CodeUserNotFound           = 3

  CodeQuit                   = -1
  CodeHelp                   = -2
)

var CmdContentRegex =
  regexp.MustCompile(`(?m)(( {0,1}~#| {0,1}~:)\[([^\[\]]*)\]| {0,1}~!!)`)

func CmdAvailable() ([]string) {
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

func CmdAutocompleter(input string, knownUsers map[string]string) ([]string) {
  var entries []string

  if input == "" {
    return entries
  }

  if input[len(input)-1:] == "@" ||
     input == "whois " {
    for _, knownUser := range knownUsers {
      line := input + knownUser
      lineFound := false

      for _, entry := range entries {
        if entry == line {
          lineFound = true
          break;
        }
      }

      if lineFound == false {
        entries = append(entries, line)
      }
    }

    return entries
  }

  for _, cmd := range CmdAvailable() {
    if strings.HasPrefix(cmd, input) == true {
      entries = append(entries, cmd + " ")
    }
  }

  return entries
}

func CmdProcessor(timeline *Timeline, input string) (CmdReturnCode) {
  split := strings.SplitN(input, " ", 2)
  cmd := split[0]
  args := split[1]

  switch cmd {
  case "home":
    timeline.Switch(TimelineHome, nil)
    return CodeOk
  case "local":
    timeline.Switch(TimelineLocal, nil)
    return CodeOk
  case "public":
    timeline.Switch(TimelinePublic, nil)
    return CodeOk
  case "notifications":
    timeline.Switch(TimelineNotifications, nil)
    return CodeOk
  case "hashtag":
    hashtag, isLocal, err := CmdHelperGetHashtagParams(args)
    if err != nil {
      return CodeNotOk
    }

    timelineOptions := TimelineOptions{
      Hashtag: hashtag,
      IsLocal: isLocal,
    }

    timeline.Switch(TimelineHashtag, &timelineOptions)
    return CodeOk
  case "whois":
    accounts, err := timeline.SearchUser(args, 1)
    if err != nil || len(accounts) < 1 {
      // TODO: pass info back to caller
      return CodeUserNotFound
    }

    timelineOptions := TimelineOptions{
      User: *accounts[0],
    }

    timeline.Switch(TimelineUser, &timelineOptions)
    return CodeOk
  case "t", "toot":
    return CmdToot(timeline, args, -1, VisibilityPublic)
  case "tp", "tootprivate":
    return CmdToot(timeline, args, -1, VisibilityPrivate)
  case "tu", "tootunlisted":
    return CmdToot(timeline, args, -1, VisibilityUnlisted)
  case "td", "tootdirect":
    return CmdToot(timeline, args, -1, VisibilityUnlisted)
  case "re", "reply":
    tootId, args, err := CmdHelperGetReplyParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdToot(timeline, args, tootId, VisibilityPublic)
  case "rep", "replyprivate":
    tootId, args, err := CmdHelperGetReplyParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdToot(timeline, args, tootId, VisibilityPrivate)
  case "reu", "replyunlisted":
    tootId, args, err := CmdHelperGetReplyParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdToot(timeline, args, tootId, VisibilityUnlisted)
  case "red", "replydirect":
    tootId, args, err := CmdHelperGetReplyParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdToot(timeline, args, tootId, VisibilityDirect)
  case "rt", "retoot", "boost":
    tootId, err := CmdHelperGetBoostParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdBoost(timeline, tootId)
  case "ut", "unretoot", "unboost":
    tootId, err := CmdHelperGetBoostParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdUnboost(timeline, tootId)
  case "fav":
    tootId, err := CmdHelperGetFavParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdFav(timeline, tootId)
  case "unfav":
    tootId, err := CmdHelperGetFavParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdUnfav(timeline, tootId)
  case "open":
    tootId, err := CmdHelperGetOpenParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdOpen(timeline, tootId)
  case "share":
    tootId, err := CmdHelperGetShareParams(args)
    if err != nil {
      return CodeNotOk
    }

    return CmdShare(timeline, tootId)
  case "?", "help":
    return CodeHelp
  case "quit", "exit", "bye":
    return CodeQuit
  }

  return CodeCommandNotFound
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
  visibility string) (CmdReturnCode) {
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
      switch token[2][(len(tokens[0][2])-2):] {
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

func CmdBoost(timeline *Timeline, tootID int) (CmdReturnCode) {
  _, err := timeline.Boost(tootID, true)
  if err != nil {
    return CodeNotOk
  }

  return CodeOk
}

func CmdUnboost(timeline *Timeline, tootID int) (CmdReturnCode) {
  _, err := timeline.Boost(tootID, false)
  if err != nil {
    return CodeNotOk
  }

  return CodeOk
}

func CmdFav(timeline *Timeline, tootID int) (CmdReturnCode) {
  _, err := timeline.Fav(tootID, true)
  if err != nil {
    return CodeNotOk
  }

  return CodeOk
}

func CmdUnfav(timeline *Timeline, tootID int) (CmdReturnCode) {
  _, err := timeline.Fav(tootID, false)
  if err != nil {
    return CodeNotOk
  }

  return CodeOk
}

func CmdOpen(timeline *Timeline, tootID int) (CmdReturnCode) {
  var err error

  url := timeline.Toots[tootID].Status.URL

  switch runtime.GOOS {
  case "darwin":
    err = exec.Command("open", url).Start()
  case "linux":
    err = exec.Command("xdg-open", url).Start()
  case "windows":
    err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
  default:
    err = errors.New("Platform not supported!")
  }

  if err != nil {
    return CodeNotOk
  }

  return CodeOk
}

func CmdShare(timeline *Timeline, tootID int) (CmdReturnCode) {
  url := timeline.Toots[tootID].Status.URL

  err := clipboard.WriteAll(url)
  if err != nil {
    return CodeNotOk
  }

  return CodeOk
}
