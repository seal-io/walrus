package intercept

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/auths/session"
)

var interceptWithProject = instructionKey("project")

// WithProjectInterceptor values the given context to indicate DAO that
// must embed the owned projects of the contextual subject into the WHERE clause of the SQL query statement.
// WithProjectInterceptor only works with those schemas using the ent.Interceptor
// created by ByProject.
func WithProjectInterceptor(ctx context.Context) context.Context {
	return context.WithValue(ctx, interceptWithProject, instructionValue{})
}

// ByProject returns an ent.Interceptor to intercept a SQL query statement,
// it works with WithProjectInterceptor.
func ByProject(idColumn string) ent.Interceptor {
	return ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
		_, intercept := ctx.Value(interceptWithProject).(instructionValue)

		if intercept {
			sj, err := session.GetSubject(ctx)
			if err != nil {
				return err
			}

			if sj.IsAdmin() {
				return nil
			}

			projIDs := make([]any, len(sj.ProjectRoles))
			for i := range sj.ProjectRoles {
				projIDs[i] = sj.ProjectRoles[i].Project.ID
			}

			t := q.(target)
			t.WhereP(func(s *sql.Selector) {
				s.Where(sql.In(s.C(idColumn), projIDs...))
			})
		}

		return nil
	})
}
