package main

import (
	"os"

	"github.com/seal-io/utils/signalx"
	"github.com/seal-io/utils/stringx"

	"github.com/seal-io/walrus/cmd"
	"github.com/seal-io/walrus/pkg/server"
)

var (
	Name  = "server"
	Brief = stringx.Title(Name) + " is a Kubernetes Extension API Server to manage Walrus Kubernetes resources."
)

func main() {
	c := cmd.Init(server.NewCommand(Name, Brief))

	if err := c.ExecuteContext(signalx.Handler()); err != nil {
		os.Exit(1)
	}
}
