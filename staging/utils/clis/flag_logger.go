package clis

import (
	"github.com/urfave/cli/v2"

	"github.com/seal-io/seal/utils/log"
)

type Logger struct{}

func (Logger) Flags(cmd *Command) {
	cmd.Flags = append(cmd.Flags,
		&cli.BoolFlag{
			Name:  "log-debug",
			Usage: "Use debugging log.",
		},
		&cli.BoolFlag{
			Name:  "log-json",
			Usage: "Log in JSON format.",
		},
		&cli.BoolFlag{
			Name:  "log-stdout",
			Usage: "Log to stdout.",
		},
		&cli.Uint64Flag{
			Name:  "log-verbosity",
			Usage: "Log verbosity level.",
		},
	)
}

func (Logger) Before(cmd *Command) {
	pb := cmd.Before
	cmd.Before = func(ctx *cli.Context) error {
		l := log.NewWrappedZapperAsLogger(ctx.Bool("log-json"), !ctx.Bool("log-debug"), ctx.Bool("log-stdout"))
		log.SetLogger(l)
		log.SetVerbosity(ctx.Uint64("log-verbosity"))
		if pb != nil {
			return pb(ctx)
		}
		return nil
	}
}
