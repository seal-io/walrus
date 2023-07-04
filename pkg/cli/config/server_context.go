package config

import (
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ServerContext contains the server config.
type ServerContext struct {
	// Server config.
	Server   string `json:"server" survey:"server"`
	Token    string `json:"token" survey:"token"`
	Insecure bool   `json:"insecure" survey:"insecure"`

	// Project config.
	ProjectID   string `json:"projectID,omitempty" survey:"project-id"`
	ProjectName string `json:"projectName,omitempty" survey:"project-name"`

	// Environment config.
	EnvironmentID   string `json:"environmentID,omitempty" survey:"environment-id"`
	EnvironmentName string `json:"environmentName,omitempty" survey:"environment-name"`
}

// OpenAPIURL generate OpenAPI url.
func (c *ServerContext) OpenAPIURL() (*url.URL, error) {
	epURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	epURL.Path, err = url.JoinPath(epURL.Path, openAPIPath)
	if err != nil {
		return nil, err
	}

	return epURL, nil
}

// InjectFields config the fields need to inject.
func (c *ServerContext) InjectFields() []string {
	return []string{"project-id", "environment-id"}
}

// Inject update the flags base on the context.
func (c *ServerContext) Inject(cmd *cobra.Command) error {
	// Skip inject while user set project name to blank.
	projName := cmd.Flags().Lookup("project-name")
	if projName != nil && projName.Changed && projName.Value.String() == "" {
		return nil
	}

	// Project id flag exist, use current project id to set it.
	fp := cmd.Flags().Lookup("project-id")
	if fp != nil && c.ProjectID != "" {
		err := fp.Value.Set(c.ProjectID)
		if err != nil {
			return err
		}
	}

	// Skip inject environment user set environment name to blank.
	fen := cmd.Flags().Lookup("environment-name")
	if fen != nil && fen.Changed && fen.Value.String() == "" {
		return nil
	}

	// Inject environment id while user doesn't set environment name.
	fe := cmd.Flags().Lookup("environment-id")
	if fe != nil && c.EnvironmentID != "" && (fen == nil || !fen.Changed) {
		err := fe.Value.Set(c.EnvironmentID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Merge generate new server context base on old context, new context and flags.
func (c *ServerContext) Merge(ns ServerContext, flags *pflag.FlagSet) ServerContext {
	merged := *c

	if ns.Server != "" {
		merged.Server = ns.Server
	}

	if ns.Token != "" {
		merged.Token = ns.Token
	}

	if flags.Changed("insecure") {
		merged.Insecure = ns.Insecure
	}

	if ns.ProjectID != "" {
		merged.ProjectID = ns.ProjectID
	}

	if ns.ProjectName != "" && flags.Changed("project-name") && ns.ProjectName != merged.ProjectName {
		merged.ProjectName = ns.ProjectName

		// Reset environment while project changed.
		merged.EnvironmentID = ""
		merged.EnvironmentName = ""
	}

	if ns.EnvironmentID != "" {
		merged.EnvironmentID = ns.EnvironmentID
	}

	if ns.EnvironmentName != "" && flags.Changed("environment-name") {
		merged.EnvironmentName = ns.EnvironmentName
	}

	return merged
}
