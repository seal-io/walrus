package session

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/slice"
)

const (
	authnSubjectGroupsKey = "authn_subject_groups"
	authnSubjectGroupKey  = "authn_subject_group"
	authnSubjectNameKey   = "authn_subject_name"
)

func StoreSubjectAuthnInfo(c *gin.Context, groups []string, name string) {
	c.Set(authnSubjectGroupsKey, groups)
	c.Set(authnSubjectGroupKey, groups[len(groups)-1])
	c.Set(authnSubjectNameKey, name)
}

const (
	authzSubjectRolesKey    = "authz_subject_roles"
	authzSubjectPoliciesKey = "authz_subject_policies"
)

func StoreSubjectAuthzInfo(c *gin.Context, roles types.SubjectRoles, policies types.RolePolicies) {
	c.Set(authzSubjectRolesKey, roles)
	c.Set(authzSubjectPoliciesKey, policies)
}

const (
	authzSubjectCurrentOperationKey = "authz_subject_current_operation"
)

func StoreSubjectCurrentOperation(c *gin.Context, operation Operation) {
	c.Set(authzSubjectCurrentOperationKey, operation)
}

func ParseSubjectKey(s string) (group, name string, err error) {
	var ss = strings.SplitN(s, "/", 2)
	if len(ss) != 2 {
		return "", "", fmt.Errorf("invalid cached subject: %s", s)
	}
	return ss[0], ss[1], nil
}

func ToSubjectKey(group, name string) string {
	return group + "/" + name
}

// LoadSubject loads the subject from the given gin.Context.
func LoadSubject(c *gin.Context) (s Subject) {
	s.Groups = c.GetStringSlice(authnSubjectGroupsKey)
	s.Group = c.GetString(authnSubjectGroupKey)
	s.Name = c.GetString(authnSubjectNameKey)
	if v, exist := c.Get(authzSubjectRolesKey); exist {
		s.Roles = v.(types.SubjectRoles)
	}
	if v, exist := c.Get(authzSubjectPoliciesKey); exist {
		s.Policies = v.(types.RolePolicies)
	}
	return
}

// LoadSubjectCurrentOperation loads the current operation of the subject from the given gin.Context.
func LoadSubjectCurrentOperation(c *gin.Context) (o Operation) {
	if v, exist := c.Get(authzSubjectCurrentOperationKey); exist {
		return v.(Operation)
	}
	return
}

type Subject struct {
	// Groups indicates all superior groups which the subject belong to,
	// includes the groups from root to the login group.
	Groups []string
	// Group indicates the group which the subject login.
	Group string
	// Name indicates the name of the subject.
	Name string
	// Roles indicates the roles of the subject.
	Roles types.SubjectRoles
	// Policies indicates the policies of the subject.
	Policies types.RolePolicies
}

// IsAnonymous return true if this subject has not been authenticated.
func (s Subject) IsAnonymous() bool {
	return s.Group == "" || s.Name == ""
}

// Key returns the key of this subject.
func (s Subject) Key() string {
	return ToSubjectKey(s.Group, s.Name)
}

// Enforce returns true if this subject has permission to access the given resource.
func (s Subject) Enforce(c *gin.Context, resource string) bool {
	var action = c.Request.Method
	var id = c.Param("id")
	var url = c.FullPath()
	for i := 0; i < len(s.Policies); i++ {
		var rp = &s.Policies[i]
		if enforce(rp, action, resource, id, url) {
			return true
		}
	}
	return false
}

func enforce(rp *types.RolePolicy, action, resource, id, url string) (allow bool) {
	// Check actions.
	switch len(rp.Actions) {
	case 0:
		return
	case 1:
		if rp.Actions[0] == "*" {
			if slice.ContainsAny[string](rp.ActionExcludes, action) {
				// Excluded action.
				return
			}
		} else if rp.Actions[0] != action {
			// Unexpected action.
			return
		}
	default:
		if !slice.ContainsAny[string](rp.Actions, action) {
			// Unexpected action.
			return
		}
	}

	// Check resources.
	switch len(rp.Resources) {
	default:
		if !slice.ContainsAny[string](rp.Resources, resource) {
			// Unexpected resource.
			return
		}
		return true
	case 1:
		if rp.Resources[0] == "*" {
			if slice.ContainsAny[string](rp.ResourceExcludes, resource) {
				// Excluded resource.
				return
			}
		} else if rp.Resources[0] != resource {
			// Unexpected resource.
			return
		}

		// Check resource ids.
		switch len(rp.ObjectIDs) {
		default:
			if !slice.ContainsAny[string](rp.ObjectIDs, id) {
				// Unexpected resource id.
				return
			}
		case 1:
			if rp.ObjectIDs[0] == "*" {
				if slice.ContainsAny[string](rp.ObjectIDExcludes, id) {
					// Excluded resource id.
					return
				}
			} else if rp.ObjectIDs[0] != id {
				// Unexpected resource id.
				return
			}
		case 0:
		}
		return true
	case 0:
	}

	// Check none resource urls.
	return slice.ContainsAny[string](rp.Paths, url)
}

// Give returns Permission of the given resource.
func (s Subject) Give(resource string) (p Permission) {
	for i := 0; i < len(s.Policies); i++ {
		var rp = &s.Policies[i]
		var pk, pv = getPermission(rp, resource)
		for k, idx := range getOperators() {
			if pk&k == 0 {
				continue
			}
			p[idx] = p[idx].merge(pv)
		}
		if pk == operatingAny {
			return
		}
	}
	return
}

func getPermission(rp *types.RolePolicy, resource string) (pk operator, pv Operation) {
	// Check resources.
	switch len(rp.Resources) {
	case 0:
		return
	case 1:
		if rp.Resources[0] == "*" {
			if slice.ContainsAny[string](rp.ResourceExcludes, resource) {
				// Excluded resource.
				return
			}
		} else if rp.Resources[0] != resource {
			// Unexpected resource.
			return
		}
	default:
		if !slice.ContainsAny[string](rp.Resources, resource) {
			// Unexpected resource.
			return
		}
	}

	// Check actions.
	switch len(rp.Actions) {
	case 0:
		return
	default:
		for i := range rp.Actions {
			switch rp.Actions[i] {
			case http.MethodGet:
				pk |= operatingGet
			case http.MethodPost:
				pk |= operatingPost
			case http.MethodPut:
				pk |= operatingPut
			case http.MethodDelete:
				pk |= operatingDelete
			case "*":
				pk |= operatingAny
			}
		}
	}

	// Check resource ids.
	switch len(rp.ObjectIDs) {
	default:
		pv.includes = rp.ObjectIDs
		return
	case 1:
		if rp.ObjectIDs[0] == "*" {
			// Checkout resource exclude ids.
			if len(rp.ObjectIDExcludes) != 0 {
				pv.excludes = rp.ObjectIDExcludes
				return
			}
		} else {
			pv.includes = rp.ObjectIDs
		}
	case 0:
	}

	// Check scope.
	pv.scope = rp.Scope
	return
}
