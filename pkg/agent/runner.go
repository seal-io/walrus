package agent

import (
	"context"
	stdlog "log"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"k8s.io/klog"
	klogv2 "k8s.io/klog/v2"

	"github.com/seal-io/seal/utils/clis"
	"github.com/seal-io/seal/utils/log"
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
	// compatible with other loggers.
	var logger = log.GetLogger()
	stdlog.SetOutput(logger)
	logrus.SetOutput(logger)
	klog.SetOutput(logger)
	klogv2.SetLogger(log.AsLogr(logger))
}

func (r *Agent) Action(cmd *cli.Command) {
	cmd.Action = func(c *cli.Context) error {
		return r.Run(c.Context)
	}
}

func (r *Agent) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
