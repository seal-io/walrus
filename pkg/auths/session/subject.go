package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/utils/slice"
)

type Role struct {
	// ID indicates the id of the role.
	ID string `json:"id"`
	// Policies indicates the policies of the role.
	Policies types.RolePolicies `json:"policies"`
}

func (r Role) enforce(res, act, rid, path string) bool {
	for i := range r.Policies {
		if enforce(&r.Policies[i], res, act, rid, path) {
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

func (r ProjectRole) match(pid object.ID, pName string) bool {
	if pid != "" && pName != "" {
		// Expect either one is set.
		return false
	}

	return r.Project.ID == pid || r.Project.Name == pName
}

func (r ProjectRole) enforce(res, act, rid, path string) bool {
	for i := range r.Roles {
		if r.Roles[i].ID == types.ProjectRoleOwner {
			return true
		}

		if r.Roles[i].enforce(res, act, rid, path) {
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
func (s Subject) Enforce(pid object.ID, pName, res, act, rid, path string) bool {
	for i := range s.Roles {
		if s.Roles[i].ID == types.SystemRoleAdmin {
			return true
		}
	}

	if pid != "" || pName != "" {
		for i := range s.ProjectRoles {
			if s.ProjectRoles[i].match(pid, pName) &&
				s.ProjectRoles[i].enforce(res, act, rid, path) {
				return true
			}
		}

		return false
	}

	for i := range s.Roles {
		if s.Roles[i].enforce(res, act, rid, path) {
			return true
		}
	}

	return false
}

const subjectIncognitoContextKey = "_subject_incognito_"

// IncognitoOn opens incognito mode,
// prevents the ent framework adopting the hooks or interceptors of mixin.Project.
func (s Subject) IncognitoOn() {
	if s.Ctx == nil {
		return
	}

	s.Ctx.Set(subjectIncognitoContextKey, true)
}

// IncognitoOff closes incognito mode,
// allows the ent framework adopting the hooks or interceptors of mixin.Project.
func (s Subject) IncognitoOff() {
	if s.Ctx == nil {
		return
	}

	s.Ctx.Set(subjectIncognitoContextKey, false)
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

var ErrIncognitoOn = errors.New("incognito subject")

// GetSubject gets the Subject from the given context.Context.
func GetSubject(ctx context.Context) (Subject, error) {
	c, err := getContext(ctx)
	if err != nil {
		return Subject{}, fmt.Errorf("failed to get context: %w", err)
	}

	if c.GetBool(subjectIncognitoContextKey) {
		return Subject{}, ErrIncognitoOn
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

func enforce(p *types.RolePolicy, res, act, rid, path string) bool {
	// Check actions.
	switch len(p.Actions) {
	default:
		if !slice.ContainsAny(p.Actions, act) {
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
		if !slice.ContainsAny(p.Resources, res) {
			// Unexpected resource.
			return false
		}

		return true
	case 1:
		if p.Resources[0] != "*" && p.Resources[0] != res {
			// Unexpected resource.
			return false
		}

		// Check resource ids.
		switch len(p.ObjectIDs) {
		default:
			if !slice.ContainsAny(p.ObjectIDs, rid) {
				// Unexpected resource id.
				return false
			}
		case 1:
			if p.ObjectIDs[0] == "*" {
				if slice.ContainsAny(p.ObjectIDExcludes, rid) {
					// Excluded resource id.
					return false
				}
			} else if p.ObjectIDs[0] != rid {
				// Unexpected resource id.
				return false
			}
		case 0:
		}

		return true
	case 0:
	}

	// Check none resource urls.
	return slice.ContainsAny(p.Paths, path)
}
