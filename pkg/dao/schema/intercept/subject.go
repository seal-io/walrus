package intercept

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/auths/session"
)

var interceptWithSubject = instructionKey("subject")

// WithSubjectInterceptor values the given context to indicate DAO that
// must embed the contextual subject ID into the WHERE clause of the SQL query statement.
// WithSubjectInterceptor only works with those schemas using the ent.Interceptor
// created by BySubject.
func WithSubjectInterceptor(ctx context.Context) context.Context {
	return context.WithValue(ctx, interceptWithSubject, instructionValue{})
}

// BySubject returns an ent.Interceptor to intercept a SQL query statement,
// it works with WithSubjectInterceptor.
func BySubject(idColumn string) ent.Interceptor {
	return ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
		_, intercept := ctx.Value(interceptWithSubject).(instructionValue)

		if intercept {
			sj, err := session.GetSubject(ctx)
			if err != nil {
				return err
			}

			t := q.(target)
			t.WhereP(func(s *sql.Selector) {
				s.Where(sql.EQ(s.C(idColumn), sj.ID))
			})
		}

		return nil
	})
}
