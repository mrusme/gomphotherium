package main

import (
	_ "embed"

	"github.com/mrusme/gomphotherium/cli"
)

//go:embed README.md
var EmbeddedHelp string

func main() {
	cli.Execute(EmbeddedHelp)
}
