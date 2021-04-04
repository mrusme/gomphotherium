package mast

import (
  "strings"
  "strconv"
  "regexp"
  "log"
  "errors"
)

type CmdReturnCode int
const (
  CodeOk CmdReturnCode       = 0
  CodeNotOk                  = 1
  CodeQuit                   = -1
)

var CmdContentRegex = regexp.MustCompile(`(?m)(( {0,1}~#| {0,1}~:)\[([^\[\]]*)\]| {0,1}~!!)`)

func CmdAvailable() ([]string) {
  return []string{
    "home",
    "local",
    "public",
    "notifications",

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

    "fav",
    "ufav",

    "open",
    "share",

    "whois",

    "search",

    "help",
    "?",

    "quit",
    "exit",
    "bye",
  }
}

func CmdAutocompleter(input string, knownUsers []string) ([]string) {
  var entries []string

  if input == "" {
    return entries
  }

  if input[len(input)-1:] == "@" {
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
    timeline.Switch(TimelineHome)
    return CodeOk
  case "local":
    timeline.Switch(TimelineLocal)
    return CodeOk
  case "public":
    timeline.Switch(TimelinePublic)
    return CodeOk
  case "notifications":
    timeline.Switch(TimelineNotifications)
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
  case "quit", "exit", "bye":
    return CodeQuit
  }

  return CodeOk
}

func CmdHelperGetReplyParams(args string) (int, string, error) {
  splitArgs := strings.SplitN(args, " ", 2)

  if len(splitArgs) < 2 {
    return -1, args, errors.New("Toot ID missing!")
  }

  tootId, err := strconv.Atoi(splitArgs[0])
  if err != nil {
    return -1, args, errors.New("Toot ID invalid!")
  }

  newArgs := splitArgs[1]

  return tootId, newArgs, nil
}

func CmdToot(timeline *Timeline, content string, inReplyTo int, visibility string) (CmdReturnCode) {
  var status string = ""
  var spoiler string = ""
  var sensitive bool = false
  var filesToUpload []string

  // this is a ~#[sample] ~:[string] with ~!! special words

  for _, token := range CmdContentRegex.FindAllStringSubmatch(content, -1) {
    if token[0] == "~!!" {
      sensitive = true
      continue
    } else if len(token) == 4 {
      switch token[2] {
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

  log.Fatalln(status)
  _, err := timeline.Toot(&status, inReplyTo, filesToUpload, &visibility, sensitive, &spoiler)
  if err != nil {
    return CodeNotOk
  }

  return CodeOk
}
