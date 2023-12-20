package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/utils/version"
)

// NewVersionCmd return version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the CLI version",
		GroupID: common.GroupOther.ID,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("walrus CLI", version.Get())
		},
	}
}
