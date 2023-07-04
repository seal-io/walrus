package mixin

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

func OwnByProject() ownByProject {
	return ownByProject{}
}

type ownByProject struct {
	mixin.Schema

	optional bool
}

func (i ownByProject) Optional() ownByProject {
	i.optional = true
	return i
}

func (i ownByProject) Fields() []ent.Field {
	f := oid.Field("projectID").
		Immutable()

	if i.optional {
		f.Comment("ID of the project to belong, empty means for all projects.").
			Optional()
	} else {
		f.Comment("ID of the project to belong.").
			NotEmpty()
	}

	return []ent.Field{f}
}

func (i ownByProject) Hooks() []ent.Hook {
	type target interface {
		SetProjectID(oid.ID)
		ProjectID() (oid.ID, bool)
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
				i.injectWith(sj, t)
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
				i.filterWith(sj, m.(target), false)
			}

			return n.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		injector,
		filter,
	}
}

func (i ownByProject) Interceptors() []ent.Interceptor {
	type target interface {
		WhereP(...func(*sql.Selector))
	}

	// Filters out the entities not belong to owner during querying.
	filter := ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
		sj, err := session.GetSubject(ctx)
		if err == nil {
			i.filterWith(sj, q.(target), true)
		}

		return nil
	})

	return []ent.Interceptor{
		filter,
	}
}

func (i ownByProject) injectWith(sj session.Subject, t interface{ SetProjectID(oid.ID) }) {
	// Inject `projectID` query value if found.
	pid := sj.Ctx.Query("projectID")
	if pid != "" {
		t.SetProjectID(oid.ID(pid))
	}
}

func (i ownByProject) filterWith(sj session.Subject, t interface{ WhereP(...func(*sql.Selector)) }, readonly bool) {
	// Filter with `projectID` query value if found.
	pid := sj.Ctx.Query("projectID")
	if pid != "" {
		sltFunc := func(ss *sql.Selector) {
			ss.Where(sql.EQ(ss.C("project_id"), pid))
		}

		if i.optional && readonly {
			// Query both project and global scope for optionally own-by-project resources.
			sltFunc = func(ss *sql.Selector) {
				ss.Where(
					sql.Or(
						sql.EQ(ss.C("project_id"), pid),
						sql.IsNull(ss.C("project_id")),
					),
				)
			}
		}

		t.WhereP(sltFunc)

		return
	}

	if sj.IsAdmin() {
		return
	}

	pids := make([]any, len(sj.ProjectRoles))
	for i := range sj.ProjectRoles {
		pids[i] = sj.ProjectRoles[i].Project.ID
	}

	sltFunc := func(ss *sql.Selector) {
		ss.Where(sql.In(ss.C("project_id"), pids...))
	}

	if i.optional {
		// Query both project and global scope for optionally own-by-project resources.
		sltFunc = func(ss *sql.Selector) {
			ss.Where(sql.Or(
				sql.In(ss.C("project_id"), pids...),
				sql.IsNull(ss.C("project_id"))))
		}
	}

	t.WhereP(sltFunc)
}
