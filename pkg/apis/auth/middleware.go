package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/auth/session"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/schema"
)

// Auth is a gin middleware,
// which is used for authenticating and authorizing the incoming request.
func Auth(enableAuthn bool, modelClient model.ClientSet) runtime.Handle {
	type authFn func(*gin.Context, model.ClientSet) error
	var authFns = []authFn{authn, authz}
	if !enableAuthn {
		authFns = []authFn{noAuth}
	}

	return func(c *gin.Context) {
		for i := range authFns {
			var err = authFns[i](c, modelClient)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func noAuth(c *gin.Context, _ model.ClientSet) error {
	var roles = schema.SubjectRoles{
		{
			Domain: "system",
			Name:   "admin",
		},
	}
	var policies = schema.RolePolicies{
		{
			Actions:   schema.RolePolicyFields("*"),
			Resources: schema.RolePolicyFields("*"),
		},
	}

	session.StoreSubjectAuthnInfo(c, "default", "admin")
	session.StoreSubjectAuthzInfo(c, roles, policies)
	return nil
}
