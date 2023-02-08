package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
)

func Livez() runtime.Handle {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	}
}
