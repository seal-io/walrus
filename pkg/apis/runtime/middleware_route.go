package runtime

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFound aborts the incoming request of not found route.
func NotFound() Handle {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// NoMethod aborts the incoming request of not method implementation.
func NoMethod() Handle {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}
