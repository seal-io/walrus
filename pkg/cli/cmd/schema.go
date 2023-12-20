package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/schema"
)

// NewSchemaCmd generate ui schema command.
func NewSchemaCmd() *cobra.Command {
	cfg := schema.GenerateOption{}

	// Command ui schema generate.
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate schema for template",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Configuration value from environment variables.
			viper.SetEnvPrefix("WALRUS")
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			viper.AutomaticEnv()
			common.BindFlags(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := schema.Generate(cfg)
			if err != nil {
				panic(err)
			}
		},
	}
	cfg.AddFlags(generateCmd)

	// Command schema.
	cmd := &cobra.Command{
		Use:     "schema",
		Short:   "Manage schema for templates",
		GroupID: common.GroupOther.ID,
	}
	cmd.AddCommand(
		generateCmd,
	)

	return cmd
}
