package server

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
)

var defaultEnvironment = &model.Environment{
	Name:        "Default",
	Description: "Default environment",
}

// initDefaultEnvironment creates a default environment.
func (r *Server) initDefaultEnvironment(ctx context.Context, opts initOptions) error {
	if count, err := opts.ModelClient.Environments().Query().Count(ctx); err != nil {
		return err
	} else if count > 0 {
		// initialized
		return nil
	}

	var creates, err = dao.EnvironmentCreates(opts.ModelClient, defaultEnvironment)
	if err != nil {
		return err
	}

	if _, err = creates[0].Save(ctx); err != nil {
		return fmt.Errorf("failed to create the default environment: %w", err)
	}

	return nil
}
