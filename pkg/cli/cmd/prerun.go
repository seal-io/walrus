package cmd

import (
	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
)

func mergeServerContext(sc *config.Config, opts *manifest.Option) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// Merge server context.
		sc.ServerContext = sc.ServerContext.Merge(config.ServerContext{
			ScopeContext: opts.ScopeContext,
		}, cmd.Flags())
	}
}
