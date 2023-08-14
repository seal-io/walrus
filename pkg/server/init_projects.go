package server

import (
	"context"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

// createBuiltinProjects creates the built-in Project resources.
func (r *Server) createBuiltinProjects(ctx context.Context, opts initOptions) error {
	mc := opts.ModelClient

	// Get system admin user.
	adminID, err := mc.Subjects().Query().
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
		_, err = mc.Projects().Query().
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

	return mc.WithTx(ctx, func(tx *model.Tx) error {
		return tx.Projects().CreateBulk().
			Set(builtin...).
			ExecE(ctx, dao.ProjectSubjectRolesEdgeSave)
	})
}
