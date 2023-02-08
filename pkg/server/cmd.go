package server

import (
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	var cmd cli.Command
	var server = New()
	server.Flags(&cmd)
	server.Before(&cmd)
	server.Action(&cmd)
	cmd.Name = "server"
	return &cmd
}
