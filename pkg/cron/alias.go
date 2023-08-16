package cron

import (
	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/utils/cron"
)

type (
	// Expr holds the definition of cron expression.
	Expr = cron.Expr

	// Task defines the interface to hold the job executing main logic.
	Task = cron.Task
)

// AwaitedExpr returns an Expr and runs in the next round.
func AwaitedExpr(raw string) Expr {
	return cron.AwaitedExpr(raw)
}

// ImmediateExpr returns an Expr and runs in the next round.
func ImmediateExpr(raw string) Expr {
	return cron.ImmediateExpr(raw)
}

// CurrentJobs return the currents running jobs.
func CurrentJobs() map[string]cron.Job {
	var (
		rj = cron.Jobs()
		cj = make(map[string]cron.Job, len(cron.Jobs()))
	)

	for n := range js {
		for _, j := range rj {
			if slices.Contains(j.Tags(), n) {
				cj[n] = j
			}
		}
	}

	return cj
}
