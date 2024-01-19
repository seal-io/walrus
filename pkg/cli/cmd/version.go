package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/seal-io/walrus/pkg/cli/api"
	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
)

// Version return version command.
func Version(sc *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the CLI and server version information",
		GroupID: common.GroupOther.ID,
		Run: func(cmd *cobra.Command, args []string) {
			info, err := api.GetVersion(sc)
			if err != nil {
				panic(err)
			}

			b, err := yaml.Marshal(info)
			if err != nil {
				panic(err)
			}

			fmt.Print(string(b))
		},
	}
}
