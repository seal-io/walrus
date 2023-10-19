package auths

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
)

func RequestAccount(mc model.ClientSet, withAuthn bool) Account {
	a := Account{
		modelClient: mc,
		filters: []filter{
			authn,
			authz,
		},
	}
	if !withAuthn {
		a.filters[0] = authnSkip
	}

	return a
}

type (
	filter  func(*gin.Context, model.ClientSet, session.Subject) (session.Subject, error)
	Account struct {
		modelClient model.ClientSet
		filters     []filter
	}
)

// Filter is a gin middleware that filters the request,
// and set the subject to the context.
func (a Account) Filter(c *gin.Context) {
	var (
		sj  session.Subject
		err error
	)

	for i := range a.filters {
		sj, err = a.filters[i](c, a.modelClient, sj)
		if err != nil {
			_ = c.Error(err)
			c.Abort()

			return
		}
	}

	session.SetSubject(c, sj)

	c.Next()
}

// Authorize implements the runtime.RouteAuthorizer interface.
func (a Account) Authorize(c *gin.Context, p runtime.RouteProfile) int {
	ress := make([]session.ActionResource, len(p.Resources))
	for i := range p.Resources {
		ress[i] = session.ActionResource{
			Name:  p.Resources[i],
			Refer: c.Param(p.ResourcePathRefers[i]),
		}
	}

	sj, _ := session.GetSubject(c)
	if sj.Enforce(p.Method, ress, c.FullPath()) {
		return http.StatusOK
	}

	if sj.IsAnonymous() {
		return http.StatusUnauthorized
	}

	return http.StatusForbidden
}
