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
    "t",
    "rt",
    "rep",
    "repall",
    "fav",
    "ufav",
    "quit",
    "exit",
    "bye",
  }
}

func CmdAutocompleter(input string) ([]string) {
  var entries []string

  if input == "" {
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
  }

  return CodeOk
}
