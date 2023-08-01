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
	// ResourceRouteAuthorizer holds the operation of authorization.
	ResourceRouteAuthorizer interface {
		// Authorize returns the status code of authorization result,
		// 200 if success, 401 if unauthorized, 403 if forbidden.
		Authorize(*gin.Context, ResourceRouteProfile) int
	}

	// ResourceAuthorizeFunc is the function type of ResourceRouteAuthorizer.
	ResourceAuthorizeFunc func(*gin.Context, ResourceRouteProfile) int
)

// Authorize implements the ResourceRouteAuthorizer interface.
func (fn ResourceAuthorizeFunc) Authorize(c *gin.Context, p ResourceRouteProfile) int {
	if fn == nil {
		return http.StatusOK
	}

	return fn(c, p)
}

// WithResourceAuthorizer if a RouterOption to configure the resourceRouteAuthorizer for the routes of IResourceHandler.
func WithResourceAuthorizer(authorizer ResourceRouteAuthorizer) RouterOption {
	return routerOption(func(r *Router) {
		r.resourceRouteAuthorizer = authorizer
	})
}
