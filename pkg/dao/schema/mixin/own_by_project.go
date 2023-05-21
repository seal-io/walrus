package mixin

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

type OwnByProject struct {
	schema

	optional bool
}

func (p OwnByProject) AsOptional() OwnByProject {
	p.optional = true
	return p
}

func (p OwnByProject) Fields() []ent.Field {
	// Keep the json tag in camel case.
	f := oid.Field("projectID").
		Immutable()

	if p.optional {
		f.Comment("ID of the project to which the resource belongs, empty means using for global level.").
			Optional()
	} else {
		f.Comment("ID of the project to which the resource belongs.").
			NotEmpty()
	}

	return []ent.Field{f}
}

func (p OwnByProject) Hooks() []ent.Hook {
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
				p.injectWith(sj, t)
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
				p.filterWith(sj, m.(target))
			}

			return n.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		injector,
		filter,
	}
}

func (p OwnByProject) Interceptors() []ent.Interceptor {
	type target interface {
		WhereP(...func(*sql.Selector))
	}

	// Filters out the entities not belong to owner during querying.
	filter := ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
		sj, err := session.GetSubject(ctx)
		if err == nil {
			p.filterWith(sj, q.(target))
		}

		return nil
	})

	return []ent.Interceptor{
		filter,
	}
}

func (p OwnByProject) injectWith(sj session.Subject, t interface{ SetProjectID(oid.ID) }) {
	// Inject `projectID` query value if found.
	pid := sj.Ctx.Query("projectID")
	if pid != "" {
		t.SetProjectID(oid.ID(pid))
	}
}

func (p OwnByProject) filterWith(sj session.Subject, t interface{ WhereP(...func(*sql.Selector)) }) {
	// Filter with `projectID` query value if found.
	pid := sj.Ctx.Query("projectID")
	if pid != "" {
		t.WhereP(func(ss *sql.Selector) {
			ss.Where(sql.EQ(ss.C("project_id"), pid))
		})

		return
	}

	if sj.IsAdmin() {
		return
	}

	pids := make([]any, len(sj.ProjectRoles))
	for i := range sj.ProjectRoles {
		pids[i] = sj.ProjectRoles[i].Project.ID
	}

	// Only affect the projects that the session subject related.
	if len(pids) != 0 {
		t.WhereP(func(ss *sql.Selector) {
			ss.Where(sql.In(ss.C("project_id"), pids...))
		})

		return
	}

	t.WhereP(func(ss *sql.Selector) {
		ss.Where(sql.False())
	})
}
