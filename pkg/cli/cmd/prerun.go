package cmd

import (
	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/api"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
)

func setupServerContextFunc(sc *config.Config, opts *manifest.Option) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// Context unchanged.
		if sc.ServerContext.Project == opts.Project && sc.ServerContext.Environment == opts.Environment {
			if !sc.Reachable {
				panic(`remote server is unreachable. You can configure cli by running "walrus login"`)
			}

			return
		}

		// Context changed.
		originalReachable := sc.Reachable

		sc.ServerContext = sc.ServerContext.Merge(config.ServerContext{
			ScopeContext: opts.ScopeContext,
		}, cmd.Flags())

		// Validate.
		err := sc.ValidateAndSetup()
		if err != nil {
			panic(err)
		}

		// Load openAPI.
		if !originalReachable {
			err = api.InitOpenAPI(sc, false)
			if err != nil {
				panic(err)
			}
		}
	}
}
