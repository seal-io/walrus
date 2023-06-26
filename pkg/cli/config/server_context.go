package config

import (
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ServerContext contains the server config.
type ServerContext struct {
	// Server config.
	Endpoint string `json:"endpoint,omitempty"`
	Token    string `json:"token"`
	Insecure bool   `json:"insecure,omitempty"`

	// Project config.
	ProjectID   string `json:"projectID,omitempty"`
	ProjectName string `json:"projectName,omitempty"`

	// Environment config.
	EnvironmentID   string `json:"environmentID,omitempty"`
	EnvironmentName string `json:"environmentName,omitempty"`
}

// OpenAPIURL generate OpenAPI url.
func (c *ServerContext) OpenAPIURL() (*url.URL, error) {
	epURL, err := url.Parse(c.Endpoint)
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
	fp := cmd.Flags().Lookup("project-id")
	if fp != nil && c.ProjectID != "" {
		err := fp.Value.Set(c.ProjectID)
		if err != nil {
			return err
		}
	}

	fe := cmd.Flags().Lookup("environment-id")
	fen := cmd.Flags().Lookup("environment-name")

	// Inject environment id while user doesn't set environment name.
	if fe != nil && c.EnvironmentID != "" && (fen == nil || fen.Value.String() == "") {
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

	if ns.Endpoint != "" {
		merged.Endpoint = ns.Endpoint
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
