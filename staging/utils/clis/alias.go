package clis

import "github.com/urfave/cli/v2"

type (
	App     = cli.App
	Command = cli.Command
	Flag    = cli.Flag

	Commands = []*Command
	Flags    = []Flag
)
