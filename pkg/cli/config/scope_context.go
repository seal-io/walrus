package config

import "github.com/spf13/cobra"

// Flags for scope context.
const (
	FlagNameProject     = "project"
	FlagNameEnvironment = "environment"
)

type ScopeContext struct {
	// Project name.
	Project string `json:"project,omitempty" survey:"project"`

	// Environment name.
	Environment string `json:"environment,omitempty" survey:"environment"`
}

func (c *ScopeContext) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&c.Project, FlagNameProject, "p", "", "Project for default use")
	cmd.Flags().StringVarP(&c.Environment, FlagNameEnvironment, "e", "", "Environment for default use")
}

func (c *ScopeContext) FlagsData() map[string]any {
	return map[string]any{
		"project":     &c.Project,
		"environment": &c.Environment,
	}
}
