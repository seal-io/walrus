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
		NewLoginCmd(cmd),
		NewContextCmd(),
		NewSchemaCmd(),
		NewApplyCmd(),
		NewDeleteCmd(),
		NewVersionCmd(),
		NewLocalCmd(),
		NewPreviewCmd(),
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

// NewLoginCmd generate login command.
func NewLoginCmd(root *cobra.Command) *cobra.Command {
	return pkgcmd.Login(serverConfig, root)
}

// NewContextCmd generate context command.
func NewContextCmd() *cobra.Command {
	return pkgcmd.Context(serverConfig)
}

// NewLocalCmd generate local command.
func NewLocalCmd() *cobra.Command {
	return pkgcmd.Local()
}

// NewSchemaCmd generate ui schema command.
func NewSchemaCmd() *cobra.Command {
	return pkgcmd.Schema()
}

// NewApplyCmd generate apply manifest command.
func NewApplyCmd() *cobra.Command {
	cmd, err := pkgcmd.Apply(serverConfig)
	if err != nil {
		panic(err)
	}

	return cmd
}

// NewDeleteCmd generate delete manifest command.
func NewDeleteCmd() *cobra.Command {
	cmd, err := pkgcmd.Delete(serverConfig)
	if err != nil {
		panic(err)
	}

	return cmd
}

// NewPreviewCmd generate preview/preview apply command.
func NewPreviewCmd() *cobra.Command {
	cmd, err := pkgcmd.Preview(serverConfig)
	if err != nil {
		panic(err)
	}

	return cmd
}

// NewVersionCmd return cli version.
func NewVersionCmd() *cobra.Command {
	return pkgcmd.Version(serverConfig)
}

// define global flags.
func globalFlags() *pflag.FlagSet {
	gf := &pflag.FlagSet{}
	gf.BoolVarP(&globalConfig.Debug, common.DebugFlag, "d", false, "Enable debug log")
	gf.BoolP(common.HelpFlag, "h", false, "Help for this command")

	return gf
}

var configSetupExample = `
  # Login to configure the Walrus CLI settings.
  $ walrus login
`

var helpTemplate = `{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
