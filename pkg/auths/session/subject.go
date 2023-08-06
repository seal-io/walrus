package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

type Role struct {
	// ID indicates the id of the role.
	ID string `json:"id"`
	// Policies indicates the policies of the role.
	Policies types.RolePolicies `json:"policies"`
}

func (r Role) enforce(act, res, resRefer, path string) bool {
	for i := range r.Policies {
		if enforce(&r.Policies[i], act, res, resRefer, path) {
			return true
		}
	}

	return false
}

type Project struct {
	// ID indicates the id of the project.
	ID object.ID `json:"id"`
	// Name indicates the name of the project.
	Name string `json:"name"`
}

type ProjectRole struct {
	// Project indicates the project of the role.
	Project Project `json:"project"`
	// Roles indicates the roles.
	Roles []Role `json:"roles"`
}

func (r ProjectRole) enforce(projRefer, act, res, resRefer, path string) bool {
	if r.Project.ID.String() != projRefer && r.Project.Name != projRefer {
		return false
	}

	for i := range r.Roles {
		if r.Roles[i].ID == types.ProjectRoleOwner {
			return true
		}

		if r.Roles[i].enforce(act, res, resRefer, path) {
			return true
		}
	}

	return false
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

// Enforce returns true if the given conditions if allowing.
func (s Subject) Enforce(projectRefer, action, resource, resourceRefer, path string) bool {
	for i := range s.Roles {
		if s.Roles[i].ID == types.SystemRoleAdmin {
			return true
		}
	}

	if projectRefer != "" {
		for i := range s.ProjectRoles {
			if s.ProjectRoles[i].enforce(projectRefer, action, resource, resourceRefer, path) {
				return true
			}
		}

		return false
	}

	for i := range s.Roles {
		if s.Roles[i].enforce(action, resource, resourceRefer, path) {
			return true
		}
	}

	return false
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
			if p.ResourceRefers[0] == "*" {
				if slices.Contains(p.ResourceReferExcludes, resRefer) {
					// Excluded resource refer.
					return false
				}
			} else if p.ResourceRefers[0] != resRefer {
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
