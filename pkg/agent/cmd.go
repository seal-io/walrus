package agent

import (
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	var cmd cli.Command
	var agent = New()
	agent.Flags(&cmd)
	agent.Before(&cmd)
	agent.Action(&cmd)
	cmd.Name = "agent"
	return &cmd
}
