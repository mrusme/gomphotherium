package mast

import (
  "strings"
)

type CmdReturnCode int
const (
  CodeOk CmdReturnCode       = 0
  CodeNotOk                  = 1
  CodeQuit                   = -1
)


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

    "rt",
    "retoot",
    "boost",

    "rep",
    "reply",
    "repall",
    "replyall",

    "fav",
    "ufav",

    "open",
    "share",

    "whois",

    "search",

    "help",

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
  case "t", "toot":
    return CmdToot(timeline, args, VisibilityPublic)
  case "tp", "tootprivate":
    return CmdToot(timeline, args, VisibilityPrivate)
  case "tu", "tootunlisted":
    return CmdToot(timeline, args, VisibilityUnlisted)
  case "td", "tootdirect":
    return CmdToot(timeline, args, VisibilityUnlisted)
  case "quit", "exit", "bye":
    return CodeQuit
  }

  return CodeOk
}

func CmdToot(timeline *Timeline, content string, visibility string) (CmdReturnCode) {
  var status string = ""
  var spoiler string = ""
  var sensitive bool = false

  splitSensitive := strings.SplitN(content, "~~!", 2)
  if len(splitSensitive) == 2 {
    sensitive = true
  }

  splitSpoiler := strings.SplitN(splitSensitive[0], "~~:", 2)
  if len(splitSpoiler) == 2 {
    spoiler = splitSpoiler[1]
  }

  status = splitSpoiler[0]

  _, err := timeline.Toot(&status, -1, nil, &visibility, sensitive, &spoiler)
  if err != nil {
    return CodeNotOk
  }

  return CodeOk
}
