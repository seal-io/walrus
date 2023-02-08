package server

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/schema"
)

func (r *Server) initSubjects(ctx context.Context, opts initOptions) error {
	var builtin = []*model.Subject{
		// group default
		{
			Kind:        "group",
			Group:       "default",
			Name:        "default",
			Description: "default/default",
			MountTo:     pointer.Bool(false),
			LoginTo:     pointer.Bool(false),
			Paths:       []string{"default"},
		},

		// user admin
		{
			Kind:        "user",
			Group:       "default",
			Name:        "admin",
			Description: "default/admin",
			MountTo:     pointer.Bool(false),
			LoginTo:     pointer.Bool(true),
			Roles: schema.SubjectRoles{
				{Domain: "system", Name: "admin"},
			},
			Paths: []string{"default", "admin"},
		},
	}

	var creates, err = dao.SubjectCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}
	for i := range creates {
		err = creates[i].
			OnConflict(
				sql.ConflictColumns(
					subject.FieldKind,
					subject.FieldGroup,
					subject.FieldName,
				),
			).
			Update(func(upsert *model.SubjectUpsert) {
				upsert.UpdateRoles()
				upsert.UpdatePaths()
			}).
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
