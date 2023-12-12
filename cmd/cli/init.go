package main

import (
	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/config"
)

// root represent the root command.
var root *cobra.Command

// Init define init steps.
func Init() error {
	sc, err := config.InitConfig()
	if err != nil {
		return err
	}

	serverConfig.ServerContext = *sc

	root = NewRootCmd()

	return nil
}
