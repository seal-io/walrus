package server

import (
	"context"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

func (r *Server) initProjects(ctx context.Context, opts initOptions) error {
	// Get system admin user.
	adminID, err := opts.ModelClient.Subjects().Query().
		Where(
			subject.Kind(types.SubjectKindUser),
			subject.Domain(types.SubjectDomainBuiltin),
			subject.Name("admin")).
		OnlyID(ctx)
	if err != nil {
		return err
	}

	candidate := []*model.Project{
		// Default project.
		{
			Name:        "default",
			Description: "Default project",
			Edges: model.ProjectEdges{
				SubjectRoles: []*model.SubjectRoleRelationship{
					{
						SubjectID: adminID,
						RoleID:    types.ProjectRoleOwner,
					},
				},
			},
		},
	}

	// Filter out the existed same name project.
	builtin := make([]*model.Project, 0, len(candidate))

	for i := range candidate {
		_, err = opts.ModelClient.Projects().Query().
			Where(project.Name(candidate[i].Name)).
			FirstID(ctx)
		if err == nil {
			continue
		}

		if !model.IsNotFound(err) {
			return err
		}

		builtin = append(builtin, candidate[i])
	}

	if len(builtin) == 0 {
		return nil
	}

	creates, err := dao.ProjectCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}

	for i := range creates {
		// Do nothing if the project has been created.
		err = creates[i].Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
