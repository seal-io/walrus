package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/seal-io/seal/pkg/cli/api"
	"github.com/seal-io/seal/pkg/cli/ask"
	"github.com/seal-io/seal/pkg/cli/config"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

const (
	flagNameServer      = "server"
	flagNameToken       = "token"
	flagNameInsecure    = "insecure"
	flagNameProject     = "project"
	flagNameEnvironment = "environment"
)

var (
	globalConfig = &config.CommonConfig{}
	serverConfig = &config.Config{}
)

// NewRootCmd generate root command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     cliName,
		Long:    "A client for Seal to manage resources",
		Version: cliVersion,
		Example: configSetupExample,
		Args:    cobra.MinimumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel(log.InfoLevel)
			if globalConfig.Debug {
				log.SetLevel(log.DebugLevel)
			}

			serverConfig.CommonConfig = *globalConfig
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
	cmd.AddCommand(
		NewConfigCmd(),
	)
	cmd.PersistentFlags().AddFlagSet(globalFlags())

	return cmd
}

// NewConfigCmd generate config command.
func NewConfigCmd() *cobra.Command {
	// Command config setup.
	cfg := config.ServerContext{}
	setupCmdFlags := pflag.NewFlagSet("config setup", pflag.ExitOnError)
	setupCmdFlags.StringVarP(&cfg.Server, flagNameServer, "s", "", "Server address, format: scheme://host:port")
	setupCmdFlags.StringVarP(&cfg.Token, flagNameToken, "", "", "Auth token to communicate to server")
	setupCmdFlags.BoolVarP(&cfg.Insecure, flagNameInsecure, "", false, "Disable SSL verification")
	setupCmdFlags.StringVarP(&cfg.Project, flagNameProject, "p", "", "Project for default use")
	setupCmdFlags.StringVarP(&cfg.Environment, flagNameEnvironment, "e", "", "Environment for default use")

	// Command config setup.
	setupCmd := &cobra.Command{
		Use:   "setup short-name",
		Short: "Connect Seal server and setup cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Configuration value from environment variables.
			viper.SetEnvPrefix("SEAL")
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			viper.AutomaticEnv()
			bindFlags(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			isInputByFlags := inputByFlags(cmd)

			// When the user does not provide any flags, interactive configuration is provided.
			if !isInputByFlags {
				qs := questions()
				err := survey.Ask(
					qs,
					&cfg,
					survey.WithHideCharacter('*'),
				)
				if err != nil {
					panic(err)
				}
			}

			set := cmd.Flags()
			err := setup(cfg, set, isInputByFlags)
			if err != nil {
				panic(err)
			}
		},
	}
	setupCmd.Flags().AddFlagSet(setupCmdFlags)

	// Command config sync.
	syncCmd := &cobra.Command{
		Use:   "sync short-name",
		Short: "Sync cli action to the latest",
		Run: func(cmd *cobra.Command, args []string) {
			err := sync()
			if err != nil {
				panic(err)
			}
		},
	}

	// Command config current context.
	currentContextCmd := &cobra.Command{
		Use:   "current-context short-name",
		Short: "Get current context",
		Run: func(cmd *cobra.Command, args []string) {
			currentContext()
		},
	}

	// Command config.
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
	}
	configCmd.AddCommand(
		setupCmd,
		syncCmd,
		currentContextCmd,
	)

	return configCmd
}

// globalFlags define global flags.
func globalFlags() *pflag.FlagSet {
	gf := &pflag.FlagSet{}
	gf.StringVarP(&globalConfig.Format, "output", "o", "table", "Output format [table, json, yaml]")
	gf.BoolVarP(&globalConfig.Debug, "debug", "d", false, "Enable debug log")
	gf.BoolP("help", "h", false, "Help for this command")
	gf.BoolP("version", "v", false, "Version for CLI")

	return gf
}

// setup define the function for config setup command.
func setup(sc config.ServerContext, flags *pflag.FlagSet, inputByFlags bool) error {
	merged := config.Config{
		ServerContext: sc,
	}

	if inputByFlags {
		origin := &config.ServerContext{}
		if serverConfig.Server != "" {
			origin = &serverConfig.ServerContext
		}

		merged = config.Config{
			ServerContext: origin.Merge(sc, flags),
		}
	} else {
		// Ignore empty.
		if merged.Project == `""` {
			merged.Project = ""
		}

		if merged.Environment == `""` {
			merged.Environment = ""
		}
	}

	err := merged.ValidateAndSetup()
	if err != nil {
		return err
	}

	serverConfig.ServerContext = merged.ServerContext

	return setServerContextToCache(serverConfig.ServerContext)
}

// sync define the function for config sync command.
func sync() error {
	err := serverConfig.ValidateAndSetup()
	if err != nil {
		return err
	}

	err = load(serverConfig, root, true)

	return err
}

// currentContext define the function for config current-context command.
func currentContext() {
	if serverConfig.Project != "" {
		name := serverConfig.Project
		if name != "" {
			fmt.Println("Current Project: " + name)
		}

		env := serverConfig.Environment
		if env != "" {
			fmt.Println("Current Environment: " + env)
		}
	}
}

// load OpenAPI from cache or remote and setup command.
func load(sc *config.Config, root *cobra.Command, skipCache bool) error {
	start := time.Now()
	defer func() {
		log.Debugf("API loading took %s", time.Since(start))
	}()

	// Load from cache while existed.
	if !skipCache {
		log.Debug("Load from cache")

		api := getAPIFromCache()
		if api != nil {
			api.GenerateCommand(sc, root)
			return nil
		}
	}

	// Load from remote.
	log.Debug("Load from remote")

	ep, err := sc.OpenAPIURL()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, ep.String(), nil)
	if err != nil {
		return err
	}

	resp, err := sc.DoRequest(req)
	defer func() { _ = resp.Body.Close() }()

	if err != nil {
		return err
	}

	api, err := api.LoadOpenAPI(resp)
	if err != nil {
		return err
	}

	api.GenerateCommand(sc, root)

	err = setAPIToCache(api)

	return err
}

var configSetupExample = `
  # Setup Seal CLI configuration
  $ seal config setup
`

var helpTemplate = `{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

// bindFlags bind the environment with flags.
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "")

		// Use viper value when the flag is not set and environment variable has a value.
		if !f.Changed && viper.IsSet(configName) {
			val := viper.Get(configName)

			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				panic(err)
			}
		}
	})
}

// inputByFlags check whether user set flags.
func inputByFlags(cmd *cobra.Command) bool {
	var change bool

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			change = true
		}
	})

	return change
}

// questions is the interactive prompt questions use to config CLI.
func questions() []*survey.Question {
	hiddenPassword := func(val string) string {
		return fmt.Sprintf("****%s", strs.LastContent(val, 4))
	}

	ap := &ask.Password{
		Message:        strs.Question(flagNameToken),
		Default:        serverConfig.Token,
		DefaultDisplay: hiddenPassword(serverConfig.Token),
	}

	proj := serverConfig.Project
	if proj == "" {
		proj = "default"
	}

	qs := []*survey.Question{
		{
			Name: flagNameServer,
			Prompt: &survey.Input{
				Message: strs.Question(flagNameServer),
				Default: serverConfig.Server,
			},
			Validate: survey.Required,
		},
		{
			Name:     flagNameToken,
			Prompt:   ap,
			Validate: ap.Required,
		},
		{
			Name: flagNameInsecure,
			Prompt: &survey.Confirm{
				Message: strs.Question(flagNameInsecure),
				Default: serverConfig.Insecure,
			},
		},
		{
			Name: flagNameProject,
			Prompt: &survey.Input{
				Message: strs.Question(flagNameProject),
				Default: proj,
			},
		},
		{
			Name: flagNameEnvironment,
			Prompt: &survey.Input{
				Message: strs.Question(flagNameEnvironment),
				Default: serverConfig.Environment,
			},
		},
	}

	return qs
}
