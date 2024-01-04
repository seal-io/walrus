package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/api"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
)

func setupServerContextFunc(sc *config.Config, opts *manifest.Option) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// Validate and init.
		err := validateAndInitOpenAPI(sc)
		if err != nil {
			panic(err)
		}

		// Merge server context.
		sc.ServerContext = sc.ServerContext.Merge(config.ServerContext{
			ScopeContext: opts.ScopeContext,
		}, cmd.Flags())
	}
}

func validateAndInitOpenAPI(sc *config.Config) error {
	err := sc.ValidateAndSetup()
	if err != nil {
		return fmt.Errorf(`cli configuration is invalid: %w. You can configure cli by running "walrus login"`, err)
	}

	shouldUpdate, err := api.CompareVersion(sc)
	if err != nil {
		return err
	}

	if shouldUpdate {
		return api.InitOpenAPI(sc, true)
	}

	return nil
}
