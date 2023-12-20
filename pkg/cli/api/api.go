package api

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/strs"
)

// API represents an abstracted API description, include details used to build CLI commands.
type API struct {
	Version    string      `json:"version"`
	Short      string      `json:"short"`
	Long       string      `json:"long,omitempty"`
	Operations []Operation `json:"operations,omitempty"`
}

// GenerateCommand generate command from api and add to root command.
func (api *API) GenerateCommand(sc *config.Config, root *cobra.Command) {
	if root.Short == "" {
		root.Short = api.Short
	}

	set := make(map[string]*cobra.Command)

	for _, op := range api.Operations {
		if op.CmdIgnore {
			continue
		}

		// Generate sub command.
		cmd := op.Command(sc)
		cmd.Flags().AddFlagSet(root.PersistentFlags())

		// Group is empty, add sub command to root.
		if op.Group == "" {
			root.AddCommand(cmd)
			continue
		}

		// Group isn't empty, Add sub command to command set.
		cmdSet := strings.ToLower(strs.Singularize(op.Group))
		_, ok := set[cmdSet]

		if !ok {
			// Generate command set.
			set[cmdSet] = &cobra.Command{
				Use:     cmdSet,
				GroupID: common.GroupManagement.ID,
				Short:   fmt.Sprintf("Manage %s", strs.Decamelize(op.Group, true)),
			}

			// Add command set to root.
			root.AddCommand(set[cmdSet])
		}

		// Add sub command to command set.
		set[cmdSet].AddCommand(cmd)
	}
}

func (api *API) GetOperation(group, name string) *Operation {
	for _, op := range api.Operations {
		if strings.ToLower(op.Group) == group && op.Name == name {
			return &op
		}
	}

	return nil
}
