package runtime

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Validator holds the operation of validation.
type Validator interface {
	// Validate returns error if the given request is invalid.
	Validate() error
}

type (
	// RouteAuthorizer holds the operation of authorization.
	RouteAuthorizer interface {
		// Authorize returns the status code of authorization result,
		// 200 if success, 401 if unauthorized, 403 if forbidden.
		Authorize(*gin.Context, RouteProfile) int
	}

	// RouteAuthorizeFunc is the function type of RouteAuthorizer.
	RouteAuthorizeFunc func(*gin.Context, RouteProfile) int
)

// Authorize implements the RouteAuthorizer interface.
func (fn RouteAuthorizeFunc) Authorize(c *gin.Context, p RouteProfile) int {
	if fn == nil {
		return http.StatusOK
	}

	return fn(c, p)
}

// WithResourceAuthorizer if a RouterOption to configure the authorizer for the routes of IResourceHandler.
func WithResourceAuthorizer(authorizer RouteAuthorizer) RouterOption {
	return routerOption(func(r *Router) {
		r.authorizer = authorizer
	})
}
