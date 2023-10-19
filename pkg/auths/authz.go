package auths

import (
	"context"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
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
		s.Roles, s.ApplicableEnvironmentTypes, err = getRolesAndApplicableEnvironmentTypes(c, mc, "system/anonymity")
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

	rids = append(rids, types.SystemRoleUser) // Inject user role.

	s.Roles, s.ApplicableEnvironmentTypes, err = getRolesAndApplicableEnvironmentTypes(c, mc, rids...)
	if err != nil {
		return s, err
	}

	if s.IsAdmin() {
		return s, nil // Return directly if the subject is admin.
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
					project.FieldName).
				WithEnvironments(func(eq *model.EnvironmentQuery) {
					eq.Order(model.Desc(environment.FieldCreateTime)).
						Unique(false).
						Where(environment.TypeNotIn(s.ApplicableEnvironmentTypes...)).
						Select(
							environment.FieldID,
							environment.FieldName)
				}).
				WithConnectors(func(cq *model.ConnectorQuery) {
					cq.Order(model.Desc(connector.FieldCreateTime)).
						Unique(false).
						Where(connector.ApplicableEnvironmentTypeNotIn(s.ApplicableEnvironmentTypes...)).
						Select(
							connector.FieldID,
							connector.FieldName)
				})
		}).
		All(c)
	if err != nil {
		return s, err
	}

	type pri struct {
		proj    *model.Project
		roles   []*model.Role
		roEnvs  []*model.Environment
		roConns []*model.Connector
	}

	var (
		pris []pri
		pi   = make(map[*model.Project]int)
	)

	for i := range rs {
		e := &rs[i].Edges
		if _, exist := pi[e.Project]; !exist {
			pi[e.Project] = len(pris)
			pris = append(pris, pri{
				proj:    e.Project,
				roEnvs:  e.Project.Edges.Environments,
				roConns: e.Project.Edges.Connectors,
			})
		}

		pris[pi[e.Project]].roles = append(pris[pi[e.Project]].roles, e.Role)
	}

	s.ProjectRoles = make([]session.ProjectRole, len(pris))

	for i := range pris {
		proj := session.Project{
			Resource: session.Resource{
				ID:   pris[i].proj.ID,
				Name: pris[i].proj.Name,
			},
		}

		proj.ReadOnlyEnvironments = make([]session.Resource, len(pris[i].roEnvs))
		for j := range pris[i].roEnvs {
			proj.ReadOnlyEnvironments[j] = session.Resource{
				ID:   pris[i].roEnvs[j].ID,
				Name: pris[i].roEnvs[j].Name,
			}
		}

		proj.ReadOnlyConnectors = make([]session.Resource, len(pris[i].roConns))
		for j := range pris[i].roConns {
			proj.ReadOnlyConnectors[j] = session.Resource{
				ID:   pris[i].roConns[j].ID,
				Name: pris[i].roConns[j].Name,
			}
		}

		projRoles := make([]session.Role, len(pris[i].roles))
		for j := range pris[i].roles {
			projRoles[j] = session.Role{
				ID:       pris[i].roles[j].ID,
				Policies: pris[i].roles[j].Policies,
			}
		}

		s.ProjectRoles[i] = session.ProjectRole{
			Project: proj,
			Roles:   projRoles,
		}
	}

	return s, nil
}

func getRolesAndApplicableEnvironmentTypes(
	ctx context.Context,
	mc model.ClientSet,
	ids ...string,
) ([]session.Role, []string, error) {
	rs, err := mc.Roles().Query().
		Order(model.Desc(role.FieldCreateTime)).
		Where(role.IDIn(ids...)).
		Select(
			role.FieldID,
			role.FieldPolicies,
			role.FieldApplicableEnvironmentTypes).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	var (
		srs = make([]session.Role, len(rs))
		us  = sets.NewString()
	)

	for i := range rs {
		srs[i] = session.Role{
			ID:       rs[i].ID,
			Policies: rs[i].Policies,
		}

		us.Insert(rs[i].ApplicableEnvironmentTypes...)
	}

	return srs, us.List(), nil
}
