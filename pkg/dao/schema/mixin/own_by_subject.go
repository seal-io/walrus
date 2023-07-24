package mixin

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

func OwnBySubject() ownBySubject {
	return ownBySubject{}
}

type ownBySubject struct {
	mixin.Schema
}

func (i ownBySubject) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("subject_id").
			Comment("ID of the subject to belong.").
			NotEmpty().
			Immutable(),
	}
}

func (i ownBySubject) Hooks() []ent.Hook {
	type target interface {
		SetSubjectID(object.ID)
		SubjectID() (object.ID, bool)
		WhereP(...func(*sql.Selector))
	}

	// Injects the owner to entity during creating.
	injector := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate) {
				return n.Mutate(ctx, m)
			}

			sj, err := session.GetSubject(ctx)
			if err == nil {
				t := m.(target)
				t.SetSubjectID(sj.ID)
			}

			return n.Mutate(ctx, m)
		})
	}

	// Filters out the entities not belong to owner during updating and deleting.
	filter := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpDelete | ent.OpDeleteOne | ent.OpUpdate | ent.OpUpdateOne) {
				return n.Mutate(ctx, m)
			}

			sj, err := session.GetSubject(ctx)
			if err == nil {
				t := m.(target)
				t.WhereP(func(s *sql.Selector) {
					s.Where(sql.EQ(s.C("subject_id"), sj.ID))
				})
			}

			return n.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		injector,
		filter,
	}
}

func (i ownBySubject) Interceptors() []ent.Interceptor {
	type target interface {
		WhereP(...func(*sql.Selector))
	}

	// Filters out the entities not belong to owner during querying.
	filter := ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
		sj, err := session.GetSubject(ctx)
		if err == nil {
			t := q.(target)
			t.WhereP(func(s *sql.Selector) {
				s.Where(sql.EQ(s.C("subject_id"), sj.ID))
			})
		}

		return nil
	})

	return []ent.Interceptor{
		filter,
	}
}
