package clis

import (
	"sort"

	"github.com/urfave/cli/v2"
)

func SortApp(app *App) {
	SortCommands(app.Commands)
	SortFlags(app.Flags)
}

func SortCommands(cmds Commands) {
	if len(cmds) == 0 {
		return
	}

	for i := 0; i < len(cmds); i++ {
		SortCommands(cmds[i].Subcommands)
		SortFlags(cmds[i].Flags)
	}
	sort.Sort(cli.CommandsByName(cmds))
}

func SortFlags(flags Flags) {
	if len(flags) == 0 {
		return
	}

	sort.Sort(cli.FlagsByName(flags))
}
