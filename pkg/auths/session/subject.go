package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type Role struct {
	// ID indicates the id of the role.
	ID string `json:"id"`
	// Policies indicates the policies of the role.
	Policies types.RolePolicies `json:"policies"`
}

type Resource struct {
	// ID indicates the id of the resource.
	ID object.ID `json:"id"`
	// Name indicates the name of the resource.
	Name string `json:"name"`
}

func (r Resource) Match(refer string) bool {
	return string(r.ID) == refer || r.Name == refer
}

type Project struct {
	Resource `json:",inline"`

	// ReadOnlyEnvironments indicates the read-only environment list of the project.
	ReadOnlyEnvironments []Resource `json:"readOnlyEnvironments"`
	// ReadOnlyConnectors indicates the read-only connector list of the project.
	ReadOnlyConnectors []Resource `json:"readOnlyConnectors"`
}

type ProjectRole struct {
	// Project indicates the project of the role.
	Project Project `json:"project"`
	// Roles indicates the roles.
	Roles []Role `json:"roles"`
}

type Subject struct {
	// Ctx holds the request gin.Context.
	Ctx *gin.Context `json:"-"`
	// ID indicates the id of the subject.
	ID object.ID `json:"id"`
	// Domain indicates the domain of the subject.
	Domain string `json:"domain"`
	// Groups indicates all superior groups which the subject belong to.
	Groups []string `json:"groups"`
	// Name indicates the name of the subject.
	Name string `json:"name"`
	// Roles indicates the roles of the subject.
	Roles []Role `json:"roles"`
	// ProjectRoles indicates the project roles of the subject.
	ProjectRoles []ProjectRole `json:"projectRoles"`
	// ApplicableEnvironmentTypes indicates which environment type are the subject to apply.
	ApplicableEnvironmentTypes []string `json:"applicableEnvironmentTypes"`
}

// IsAnonymous returns true if this subject has not been authenticated.
func (s Subject) IsAnonymous() bool {
	return s.Name == ""
}

// IsAdmin returns true if this subject has a system/admin role.
func (s Subject) IsAdmin() bool {
	for i := range s.Roles {
		if s.Roles[i].ID == types.SystemRoleAdmin {
			return true
		}
	}

	return false
}

// IsApplicableEnvironmentType returns true if the given environment type is applicable.
func (s Subject) IsApplicableEnvironmentType(et string) bool {
	return slices.Contains(s.ApplicableEnvironmentTypes, et)
}

type ActionResource struct {
	// Name indicates the name of resource.
	Name string
	// Refer indicates the refer of resource, either ID or name.
	Refer string
}

// Enforce returns true if the given conditions if allowing.
func (s Subject) Enforce(action string, resources []ActionResource, path string) bool {
	const allow, reject = true, false

	var (
		checkProjectRoles           bool
		targetResource, targetRefer string
	)

	if sz := len(resources); sz != 0 {
		checkProjectRoles = resources[0].Name == "projects" && resources[0].Refer != ""
		targetResource, targetRefer = resources[sz-1].Name, resources[sz-1].Refer
	}

	for i := range s.Roles {
		r := &s.Roles[i]

		if r.ID == types.SystemRoleAdmin {
			return allow
		}

		if checkProjectRoles {
			continue
		}

		if !enforces(r.Policies, action, targetResource, targetRefer, path) {
			continue
		}

		return allow
	}

	if !checkProjectRoles {
		return reject // Reject if failed in the validation of all system roles.
	}

	projectRefer := resources[0].Refer
	subResources := resources[1:]

	for i := range s.ProjectRoles {
		pr := &s.ProjectRoles[i]
		p := &pr.Project

		if !p.Match(projectRefer) {
			continue
		}

		for j := range pr.Roles {
			if pr.Roles[j].ID != types.ProjectRoleOwner &&
				!enforces(pr.Roles[j].Policies, action, targetResource, targetRefer, path) {
				continue
			}

			if action != http.MethodGet && len(subResources) > 0 {
				var rejectList []Resource

				switch subResources[0].Name {
				case "environments":
					rejectList = p.ReadOnlyEnvironments
				case "connectors":
					rejectList = p.ReadOnlyConnectors
				}

				for k := range rejectList {
					if rejectList[k].Match(subResources[0].Refer) {
						return reject // Reject if the subresource is readonly.
					}
				}
			}

			return allow
		}

		break
	}

	return reject // Reject if failed in the validation of all project roles.
}

const subjectContextKey = "_subject_"

// SetSubject sets the Subject into the given context.Context.
func SetSubject(ctx context.Context, subject Subject) {
	c, err := getContext(ctx)
	if err != nil {
		return
	}

	subject.Ctx = c
	c.Set(subjectContextKey, subject)
}

// GetSubject gets the Subject from the given context.Context.
func GetSubject(ctx context.Context) (Subject, error) {
	c, err := getContext(ctx)
	if err != nil {
		return Subject{}, fmt.Errorf("failed to get context: %w", err)
	}

	v, ok := c.Get(subjectContextKey)
	if !ok {
		return Subject{}, fmt.Errorf("failed to get subject: %w", err)
	}

	s, ok := v.(Subject)
	if !ok {
		return Subject{}, errors.New("failed to unwrap subject")
	}

	return s, nil
}

// MustGetSubject is similar to GetSubject,
// but panic if error raising.
func MustGetSubject(ctx context.Context) Subject {
	s, err := GetSubject(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to load subject: %w", err))
	}

	return s
}

func getContext(ctx context.Context) (*gin.Context, error) {
	if ctx == nil {
		return nil, errors.New("nil ctx")
	}

	c, ok := ctx.(*gin.Context)
	if ok {
		return c, nil
	}

	c, ok = ctx.Value(gin.ContextKey).(*gin.Context)
	if ok {
		return c, nil
	}

	return nil, errors.New("not gin context")
}

func enforces(ps types.RolePolicies, act, res, resRefer, path string) bool {
	for i := range ps {
		if enforce(&ps[i], act, res, resRefer, path) {
			return true
		}
	}

	return false
}

func enforce(p *types.RolePolicy, act, res, resRefer, path string) bool {
	// Check actions.
	switch len(p.Actions) {
	default:
		if !slices.Contains(p.Actions, act) {
			// Unexpected action.
			return false
		}
	case 1:
		if p.Actions[0] != "*" && p.Actions[0] != act {
			// Unexpected action.
			return false
		}
	case 0:
		// Unexpected action.
		return false
	}

	// Check resources.
	switch len(p.Resources) {
	default:
		if !slices.Contains(p.Resources, res) {
			// Unexpected resource.
			return false
		}

		return true
	case 1:
		if p.Resources[0] != "*" && p.Resources[0] != res {
			// Unexpected resource.
			return false
		}

		// Check resource refers.
		switch len(p.ResourceRefers) {
		default:
			if !slices.Contains(p.ResourceRefers, resRefer) {
				// Unexpected resource refer.
				return false
			}
		case 1:
			if p.ResourceRefers[0] != "*" && p.ResourceRefers[0] != resRefer {
				// Unexpected resource refer.
				return false
			}
		case 0:
		}

		return true
	case 0:
	}

	// Check none resource urls.
	return slices.Contains(p.Paths, path)
}
