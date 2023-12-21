package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	pkgcmd "github.com/seal-io/walrus/pkg/cli/cmd"
	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/log"
)

var (
	globalConfig = &config.CommonConfig{}
	serverConfig = &config.Config{}
)

// NewRootCmd generate root command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     cliName,
		Short:   "Command line interface for Walrus",
		Example: configSetupExample,
		Args:    cobra.MinimumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if globalConfig.Debug {
				log.SetLevel(log.DebugLevel)
			}

			serverConfig.CommonConfig = *globalConfig
		},
	}

	cmd.AddGroup(
		common.GroupManagement,
		common.GroupAdvanced,
		common.GroupOther)

	cmd.AddCommand(
		NewConfigCmd(cmd),
		NewSchemaCmd(),
		NewApplyCmd(),
		NewDeleteCmd(),
		NewVersionCmd(),
	)

	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
	cmd.PersistentFlags().AddFlagSet(globalFlags())
	cmd.SetCompletionCommandGroupID(common.GroupOther.ID)
	cmd.SetHelpCommandGroupID(common.GroupOther.ID)

	return cmd
}

// NewConfigCmd generate config command.
func NewConfigCmd(root *cobra.Command) *cobra.Command {
	return pkgcmd.NewConfigCmd(serverConfig, root)
}

// NewSchemaCmd generate ui schema command.
func NewSchemaCmd() *cobra.Command {
	return pkgcmd.NewSchemaCmd()
}

// NewApplyCmd apply manifest.
func NewApplyCmd() *cobra.Command {
	cmd, err := pkgcmd.Apply(serverConfig)
	if err != nil {
		panic(err)
	}

	return cmd
}

// NewDeleteCmd apply manifest.
func NewDeleteCmd() *cobra.Command {
	cmd, err := pkgcmd.Delete(serverConfig)
	if err != nil {
		panic(err)
	}

	return cmd
}

// NewVersionCmd return cli version.
func NewVersionCmd() *cobra.Command {
	return pkgcmd.NewVersionCmd()
}

// define global flags.
func globalFlags() *pflag.FlagSet {
	gf := &pflag.FlagSet{}
	gf.BoolVarP(&globalConfig.Debug, "debug", "d", false, "Enable debug log")
	gf.BoolP("help", "h", false, "Help for this command")

	return gf
}

var configSetupExample = `
  # Setup Walrus CLI configuration
  $ walrus config setup
`

var helpTemplate = `{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
