package cmd

import (
	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
)

func Apply(sc *config.Config) (*cobra.Command, error) {
	flags := &manifest.OperateOption{}

	use := "apply short-name"
	cmd := &cobra.Command{
		Use:     use,
		GroupID: common.GroupAdvanced.ID,
		Short:   "Apply a configuration to a resource using the provided file path or folder.",
		Run: func(cmd *cobra.Command, args []string) {
			err := manifest.Apply(sc, flags)
			if err != nil {
				panic(err)
			}
		},
	}

	flags.AddFlags(cmd)

	return cmd, nil
}
