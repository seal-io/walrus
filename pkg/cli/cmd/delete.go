package cmd

import (
	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
)

func Delete(sc *config.Config) (*cobra.Command, error) {
	flags := &manifest.OperateOption{}

	use := "delete short-name"
	cmd := &cobra.Command{
		Use:     use,
		GroupID: common.GroupAdvanced.ID,
		Short:   "Delete resource using the provided file path or folder.",
		Run: func(cmd *cobra.Command, args []string) {
			err := manifest.Delete(sc, flags)
			if err != nil {
				panic(err)
			}
		},
	}

	flags.AddFlags(cmd)

	return cmd, nil
}
