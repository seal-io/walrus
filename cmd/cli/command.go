package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/seal-io/seal/pkg/cli/api"
	"github.com/seal-io/seal/pkg/cli/config"
	"github.com/seal-io/seal/utils/log"
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
	cmd.AddGroup(NewCmdGroups()...)
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
	setupCmdFlags.StringVarP(&cfg.Endpoint, "server", "s", "", "Server address, format: scheme://host:port")
	setupCmdFlags.StringVarP(&cfg.ProjectName, "project-name", "p", "default", "Project for default use")
	setupCmdFlags.BoolVarP(&cfg.Insecure, "insecure", "", false, "Disable SSL verification")
	setupCmdFlags.StringVarP(&cfg.Token, "token", "", "", "Auth token to communicate to server")
	setupCmdFlags.StringVarP(&cfg.EnvironmentName, "environment-name", "e", "", "Environment for default use")

	// Command config setup.
	setupCmd := &cobra.Command{
		Use:   "setup short-name",
		Short: "Connect Seal server and setup cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			viper.SetEnvPrefix("SEAL")
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			viper.AutomaticEnv()
			bindFlags(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			set := cmd.Flags()
			err := setup(cfg, set)
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
		GroupID: "config",
		Use:     "config",
		Short:   "Command set for manage CLI configuration",
	}
	configCmd.AddCommand(
		setupCmd,
		syncCmd,
		currentContextCmd,
	)

	return configCmd
}

// NewCmdGroups generate command group.
func NewCmdGroups() []*cobra.Group {
	configGroup := &cobra.Group{ID: "config", Title: "config commands:"}

	return []*cobra.Group{
		configGroup,
	}
}

// globalFlags define global flags.
func globalFlags() *pflag.FlagSet {
	gf := &pflag.FlagSet{}
	gf.StringVarP(&globalConfig.Format, "output", "o", "table", "Output format [table, json, yaml]")
	gf.BoolVarP(&globalConfig.Debug, "debug", "d", false, "Enable debug log")

	return gf
}

// setup define the function for config setup command.
func setup(sc config.ServerContext, flags *pflag.FlagSet) error {
	origin := &config.ServerContext{}
	if serverConfig.Endpoint != "" {
		origin = &serverConfig.ServerContext
	}

	merged := config.Config{
		ServerContext: origin.Merge(sc, flags),
	}

	err := merged.ValidateAndSetup()
	if err != nil {
		return err
	}

	serverConfig.ServerContext = merged.ServerContext

	return setServerContextToCache(serverConfig.ServerContext)
}

// setup define the function for config sync command.
func sync() error {
	err := serverConfig.ValidateAndSetup()
	if err != nil {
		return err
	}

	err = load(serverConfig, root, true)

	return err
}

// setup define the function for config current-context command.
func currentContext() {
	if serverConfig.ProjectName != "" {
		name := serverConfig.ProjectName
		if name != "" {
			fmt.Println("Current Project: " + name)
		}

		env := serverConfig.EnvironmentName
		if env != "" {
			fmt.Println("Current Environment:" + env)
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
  # Setup seal cli
  $ seal config setup --server [Seal_Server_URL] --project-name [Project_Name] --token [Token]
`

var helpTemplate = `{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

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
