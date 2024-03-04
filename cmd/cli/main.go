package main

import (
	"github.com/seal-io/utils/signalx"
	"github.com/seal-io/utils/stringx"
)

var (
	Name  = "cli"
	Brief = stringx.Title(Name) + " is a CLI to operate the Walrus."
)

func main() {
	_ = signalx.Handler()
}
