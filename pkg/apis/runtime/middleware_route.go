package runtime

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// noMethod is a gin middleware,
// it aborts the incoming request of not method implementation.
func noMethod(c *gin.Context) {
	c.AbortWithStatus(http.StatusMethodNotAllowed)
}

// noRoute is a gin middleware,
// it aborts the incoming request of not found route.
func noRoute(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotFound)
}
