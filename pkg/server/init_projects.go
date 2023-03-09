package server

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
)

var defaultProject = &model.Project{
	Name:        "Default",
	Description: "Default project",
}

// initDefaultProject creates a default project.
func (r *Server) initDefaultProject(ctx context.Context, opts initOptions) error {
	if count, err := opts.ModelClient.Projects().Query().Count(ctx); err != nil {
		return err
	} else if count > 0 {
		// initialized
		return nil
	}

	var creates, err = dao.ProjectCreates(opts.ModelClient, defaultProject)
	if err != nil {
		return err
	}

	if _, err = creates[0].Save(ctx); err != nil {
		return fmt.Errorf("failed to create the default project: %w", err)
	}

	return nil
}
