package server

import (
	"context"

	"github.com/seal-io/walrus/pkg/templates"
)

// setupEnvVariables sets up the environment variables for the server.
func (r *Server) setupEnvVariables(ctx context.Context, opts initOptions) error {
	return templates.SetGitCAEnvVar(ctx, opts.ModelClient)
}
