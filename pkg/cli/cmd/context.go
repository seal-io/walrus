package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/strs"
)

// Context generate context command.
func Context(serverConfig *config.Config) *cobra.Command {
	cfgCtx := config.ScopeContext{}

	// Command switch context.
	switchCmd := &cobra.Command{
		Use:   "switch",
		Short: "Switch CLI context",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Configuration value from environment variables.
			viper.SetEnvPrefix("WALRUS")
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			viper.AutomaticEnv()
			common.BindFlags(cmd)

			// Validate and init.
			err := validateAndInitOpenAPI(serverConfig)
			if err != nil {
				panic(err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			isInputByFlags := common.InputByFlags(cmd)

			// When the user does not provide any flags, interactive configuration is provided.
			if !isInputByFlags {
				qs := contextSwitchQuestions(serverConfig)
				err := survey.Ask(
					qs,
					&cfgCtx,
					survey.WithHideCharacter('*'),
				)
				if err != nil {
					panic(err)
				}
			}

			set := cmd.Flags()
			err := setup(
				config.ServerContext{
					ScopeContext: cfgCtx,
					Server:       serverConfig.Server,
					Token:        serverConfig.Token,
					Insecure:     serverConfig.Insecure,
				},
				serverConfig,
				set,
				isInputByFlags)
			if err != nil {
				panic(err)
			}

			fmt.Println("Switched context successfully.")
			contextCurrent(serverConfig)
		},
	}

	cfgCtx.AddFlags(switchCmd)

	// Command context current.
	currentCmd := &cobra.Command{
		Use:   "current",
		Short: "Get current context",
		Run: func(cmd *cobra.Command, args []string) {
			contextCurrent(serverConfig)
		},
	}

	// Command context switch.
	contextCmd := &cobra.Command{
		Use:     "context",
		Short:   "Manage CLI context",
		GroupID: common.GroupOther.ID,
	}

	contextCmd.AddCommand(
		currentCmd,
		switchCmd,
	)

	return contextCmd
}

// contextSwitchQuestions is the interactive prompt questions use to switch context.
func contextSwitchQuestions(serverConfig *config.Config) []*survey.Question {
	proj := serverConfig.Project
	if proj == "" {
		proj = "default"
	}

	qs := []*survey.Question{
		{
			Name: config.FlagNameProject,
			Prompt: &survey.Input{
				Message: strs.Question(config.FlagNameProject),
				Default: proj,
			},
		},
		{
			Name: config.FlagNameEnvironment,
			Prompt: &survey.Input{
				Message: strs.Question(config.FlagNameEnvironment),
				Default: serverConfig.Environment,
			},
		},
	}

	return qs
}

// contextCurrent define the function for context current command.
func contextCurrent(serverConfig *config.Config) {
	toBeConfigured := "to be configured"

	project := serverConfig.Project
	if project == "" {
		project = toBeConfigured
	}

	env := serverConfig.Environment
	if env == "" {
		env = toBeConfigured
	}

	fmt.Println("Current Project: " + project)
	fmt.Println("Current Environment: " + env)
}
