package measure

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/health"
	"github.com/seal-io/seal/pkg/metric"
)

func Readyz() runtime.Handle {
	return func(c *gin.Context) {
		d, ok := health.MustValidate(c, []string{"k8sctrl", "casdoor"})
		if !ok {
			c.String(http.StatusServiceUnavailable, d)
			return
		}

		c.String(http.StatusOK, d)
	}
}

func Livez() runtime.Handle {
	return func(c *gin.Context) {
		d, ok := health.Validate(c, c.QueryArray("exclude")...)
		if !ok {
			c.String(http.StatusServiceUnavailable, d)
			return
		}

		c.String(http.StatusOK, d)
	}
}

func Metrics() runtime.HTTPHandler {
	return metric.Index(5, 30*time.Second)
}
