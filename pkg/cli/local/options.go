package local

import "github.com/spf13/cobra"

type InstallOptions struct {
	Env []string // List of environment variable to set in the local Walrus.
}

func (o *InstallOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringSliceVarP(&o.Env, "env", "e", []string{}, "Set environment variables for the local Walrus")
}
