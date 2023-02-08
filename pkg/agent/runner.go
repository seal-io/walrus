package agent

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/seal-io/seal/utils/clis"
)

type Agent struct {
	Logger clis.Logger
}

func New() *Agent {
	return &Agent{}
}

func (r *Agent) Flags(cmd *cli.Command) {
	var flags = [...]cli.Flag{}
	for i := range flags {
		cmd.Flags = append(cmd.Flags, flags[i])
	}
}

func (r *Agent) Before(cmd *cli.Command) {
	r.Logger.Before(cmd)
}

func (r *Agent) Action(cmd *cli.Command) {
	cmd.Action = func(c *cli.Context) error {
		return r.Run(c.Context)
	}
}

func (r *Agent) Run(ctx context.Context) error {
	return nil
}
