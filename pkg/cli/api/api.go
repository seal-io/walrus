package api

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/seal-io/seal/pkg/cli/config"
	"github.com/seal-io/seal/utils/strs"
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

	if root.Long == "" {
		root.Long = api.Long
	}

	group := make(map[string]*cobra.Command)

	for _, op := range api.Operations {
		// Generate sub command.
		cmd := op.Command(sc)
		cmd.Flags().AddFlagSet(root.PersistentFlags())

		// Group is empty, add sub command to root.
		if op.Group == "" {
			root.AddCommand(cmd)
			continue
		}

		// Group isn't empty, Add sub command to group.
		_, ok := group[op.Group]
		if !ok {
			// Generate group command.
			group[op.Group] = &cobra.Command{
				GroupID: op.Group,
				Use:     op.Group,
				Short:   fmt.Sprintf("Command set for %s management", op.Group),
			}
		}

		group[op.Group].AddCommand(cmd)
	}

	for i := range group {
		// Add group command to root.
		root.AddCommand(group[i])

		// Add group to root.
		gp := &cobra.Group{
			ID:    group[i].GroupID,
			Title: fmt.Sprintf("%s Commands:", strs.Decamelize(group[i].GroupID, false)),
		}
		root.AddGroup(gp)
	}
}
