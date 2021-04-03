package mast

import (
  "strings"
)

type CmdReturnCode int
const (
  CodeOk CmdReturnCode       = 0
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

func CmdProcessor(input string) (CmdReturnCode) {
  split := strings.Split(input, " ")
  switch split[0] {
  case "quit", "exit", "bye":
    return CodeQuit
  case "t", "toot":

  }

  return CodeOk
}
