package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/walrus/pkg/cli/api"
	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/local"
)

// Local generate local command.
func Local() *cobra.Command {
	// Command install.
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install local Walrus",
		Run: func(cmd *cobra.Command, args []string) {
			if err := install(); err != nil {
				panic(err)
			}
		},
	}

	// Command uninstall.
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall local Walrus",
		Run: func(cmd *cobra.Command, args []string) {
			if err := uninstall(); err != nil {
				panic(err)
			}
		},
	}

	// Command local.
	localCmd := &cobra.Command{
		Use:     "local",
		Short:   "Manage local Walrus setup",
		GroupID: common.GroupOther.ID,
	}
	localCmd.AddCommand(
		installCmd,
		uninstallCmd,
	)

	return localCmd
}

// install define the function for installing local Walrus.
func install() error {
	if err := local.InstallLocalWalrus(); err != nil {
		return fmt.Errorf("failed to install local Walrus: %w", err)
	}

	fmt.Println("Checking readiness...")

	cfg := &config.Config{
		ServerContext: config.ServerContext{
			Server:   "https://localhost:7443",
			Insecure: true,
			ScopeContext: config.ScopeContext{
				Project:     "default",
				Environment: "local",
			},
		},
	}

	err := wait.PollUntilContextTimeout(
		context.Background(),
		5*time.Second,
		5*time.Minute,
		false,
		func(ctx context.Context) (done bool, err error) {
			// nolint:nilerr
			if err = cfg.CheckReachable(); err != nil {
				return false, nil
			}

			return true, nil
		},
	)
	if err != nil {
		return err
	}

	if err = api.InitOpenAPI(cfg, true); err != nil {
		return err
	}

	if err = config.SetServerContextToCache(cfg.ServerContext); err != nil {
		return err
	}

	fmt.Println("Walrus CLI is configured.")

	return nil
}

// uninstall define the function for uninstalling local Walrus.
func uninstall() error {
	if err := local.UninstallLocalWalrus(); err != nil {
		return fmt.Errorf("failed to uninstall local Walrus: %w", err)
	}

	fmt.Println("Successfully uninstalled local Walrus.")

	return nil
}
