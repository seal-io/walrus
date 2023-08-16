package auths

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/role"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/model/subjectrolerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
)

func authz(c *gin.Context, mc model.ClientSet, s session.Subject) (session.Subject, error) {
	var err error

	if s.IsAnonymous() {
		// Query anonymity role.
		s.Roles, err = getRoles(c, mc, "system/anonymity")
		return s, err
	}

	// Query assigned roles.
	q := mc.Subjects().Query()

	if len(s.Groups) == 0 {
		q.Where(
			subject.Kind(types.SubjectKindUser),
			subject.Domain(s.Domain),
			subject.Name(s.Name))
	} else {
		ps := make([]predicate.Subject, 0, len(s.Groups)+1)
		ps = append(ps, subject.And(
			subject.Kind(types.SubjectKindUser),
			subject.Domain(s.Domain),
			subject.Name(s.Name)))

		for i := range s.Groups {
			ps = append(ps, subject.And(
				subject.Kind(types.SubjectKindGroup),
				subject.Domain(s.Domain),
				subject.Name(s.Groups[i])))
		}

		q.Where(subject.Or(ps...))
	}

	var rids []string

	err = q.Clone().QueryRoles().
		Where(subjectrolerelationship.ProjectIDIsNil()).
		Select(subjectrolerelationship.FieldRoleID).
		Scan(c, &rids)
	if err != nil {
		return s, err
	}

	rids = append(rids, "system/user") // Inject user role.

	s.Roles, err = getRoles(c, mc, rids...)
	if err != nil {
		return s, err
	}

	// Query project roles.
	rs, err := q.QueryRoles().
		Where(subjectrolerelationship.ProjectIDNotNil()).
		Select(
			subjectrolerelationship.FieldRoleID,
			subjectrolerelationship.FieldProjectID).
		WithRole(func(rq *model.RoleQuery) {
			rq.Order(model.Desc(role.FieldCreateTime)).
				Unique(false).
				Select(
					role.FieldID,
					role.FieldPolicies)
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Order(model.Desc(project.FieldCreateTime)).
				Unique(false).
				Select(
					project.FieldID,
					project.FieldName)
		}).
		All(c)
	if err != nil {
		return s, err
	}

	type pri struct {
		p  *model.Project
		rs []*model.Role
	}

	var (
		pris []pri
		pi   = make(map[*model.Project]int)
	)

	for i := range rs {
		e := &rs[i].Edges
		if _, exist := pi[e.Project]; !exist {
			pi[e.Project] = len(pris)
			pris = append(pris, pri{p: e.Project})
		}

		pris[pi[e.Project]].rs = append(pris[pi[e.Project]].rs, e.Role)
	}

	s.ProjectRoles = make([]session.ProjectRole, len(pris))

	for i := range pris {
		prs := make([]session.Role, len(pris[i].rs))
		for j := range pris[i].rs {
			prs[j] = session.Role{
				ID:       pris[i].rs[j].ID,
				Policies: pris[i].rs[j].Policies,
			}
		}

		s.ProjectRoles[i] = session.ProjectRole{
			Project: session.Project{
				ID:   pris[i].p.ID,
				Name: pris[i].p.Name,
			},
			Roles: prs,
		}
	}

	return s, err
}

func getRoles(ctx context.Context, mc model.ClientSet, ids ...string) (r []session.Role, err error) {
	rs, err := mc.Roles().Query().
		Order(model.Desc(role.FieldCreateTime)).
		Where(role.IDIn(ids...)).
		Select(
			role.FieldID,
			role.FieldPolicies).
		All(ctx)
	if err != nil {
		return
	}

	r = make([]session.Role, len(rs))
	for i := range rs {
		r[i] = session.Role{
			ID:       rs[i].ID,
			Policies: rs[i].Policies,
		}
	}

	return
}
