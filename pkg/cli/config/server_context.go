package config

import (
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var InjectFields = []string{"project", "environment"}

// Flags for server context.
const (
	FlagNameServer   = "server"
	FlagNameToken    = "token"
	FlagNameInsecure = "insecure"
)

// ServerContext contains the server config.
type ServerContext struct {
	// Scope context.
	ScopeContext

	// Server config.
	Server   string `json:"server" survey:"server"`
	Token    string `json:"token" survey:"token"`
	Insecure bool   `json:"insecure" survey:"insecure"`
}

func (c *ServerContext) AddFlags(cmd *cobra.Command) {
	c.ScopeContext.AddFlags(cmd)

	cmd.Flags().StringVarP(&c.Server, FlagNameServer, "s", "", "Server address, format: scheme://host:port")
	cmd.Flags().StringVarP(&c.Token, FlagNameToken, "", "", "Auth token to communicate to server")
	cmd.Flags().BoolVarP(&c.Insecure, FlagNameInsecure, "", false, "Disable SSL verification")
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
	return InjectFields
}

// Inject update the flags base on the context.
func (c *ServerContext) Inject(cmd *cobra.Command) error {
	// Skip inject while user set project name to blank.
	projName := cmd.Flags().Lookup("project")
	if projName != nil && projName.Changed && projName.Value.String() == "" {
		return nil
	}

	// Project flag exist, use current project to set it.
	fp := cmd.Flags().Lookup("project")
	if fp != nil && !fp.Changed && c.Project != "" {
		err := fp.Value.Set(c.Project)
		if err != nil {
			return err
		}
	}

	// Skip inject environment while user set environment name to blank.
	fen := cmd.Flags().Lookup("environment")
	if fen != nil && fen.Changed && fen.Value.String() == "" {
		return nil
	}

	// Inject environment while user doesn't set environment name.
	fe := cmd.Flags().Lookup("environment")
	if fen != nil && !fen.Changed && c.Environment != "" {
		err := fe.Value.Set(c.Environment)
		if err != nil {
			return err
		}
	}

	return nil
}

// InjectURI update the uri base on the context.
func (c *ServerContext) InjectURI(uri, name string) string {
	const (
		projectPlaceholder     = "{project}"
		environmentPlaceholder = "{environment}"
	)

	switch name {
	case "project":
		// Inject project name.
		if c.Project != "" && !strings.HasSuffix(uri, projectPlaceholder) {
			uri = strings.Replace(uri, projectPlaceholder, c.Project, 1)
		}
	case "environment":
		// Inject environment name.
		if c.Environment != "" && !strings.HasSuffix(uri, environmentPlaceholder) {
			uri = strings.Replace(uri, environmentPlaceholder, c.Environment, 1)
		}
	}

	return uri
}

// ContextExisted check whether the value already existed in the context.
func (c *ServerContext) ContextExisted(name string) bool {
	switch {
	case name == "project" && c.Project != "":
		return true
	case name == "environment" && c.Environment != "":
		return true
	default:
		return false
	}
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

	if ns.Project != "" && flags.Changed("project") && ns.Project != merged.Project {
		merged.Project = ns.Project

		// Reset environment while project changed.
		merged.Environment = ""
	}

	if ns.Environment != "" && flags.Changed("environment") {
		merged.Environment = ns.Environment
	}

	return merged
}
