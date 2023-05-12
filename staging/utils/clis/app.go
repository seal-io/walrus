package clis

import (
	"github.com/urfave/cli/v2"

	"github.com/seal-io/seal/utils/version"
)

func AsApp(cmd *Command) *App {
	app := cli.NewApp()

	app.Name = cmd.Name
	app.HelpName = cmd.HelpName
	app.Usage = cmd.Usage
	app.UsageText = cmd.UsageText
	app.ArgsUsage = cmd.ArgsUsage
	app.Version = version.Get()
	app.Description = cmd.Description
	app.Commands = cmd.Subcommands
	app.Flags = cmd.Flags
	app.EnableBashCompletion = true
	app.HideHelp = cmd.HideHelp
	app.HideHelpCommand = cmd.HideHelpCommand
	app.HideVersion = false
	app.BashComplete = cmd.BashComplete
	app.Before = cmd.Before
	app.After = cmd.After
	app.Action = cmd.Action
	app.OnUsageError = cmd.OnUsageError
	SortApp(app)
	MutateAppEnvs(app)

	return app
}
