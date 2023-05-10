package runtime

import (
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/utils/log"
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

// Logging is a gin middleware,
// which is the same as gin.Logger but uses a unified logging tool.
func Logging(ignorePaths ...string) Handle {
	logger := log.WithName("api")
	if !logger.Enabled(log.DebugLevel) {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	skipPaths := sets.Set[string]{}
	skipPrefixPaths := sets.Set[string]{}
	for i := range ignorePaths {
		skipPaths.Insert(ignorePaths[i])
		lastIdx := strings.LastIndex(ignorePaths[i], "/") + 1
		if lastIdx <= 0 {
			continue
		}
		lastSeg := ignorePaths[i][lastIdx:]
		if !strings.HasPrefix(lastSeg, "*") {
			continue
		}
		skipPrefixPaths.Insert(ignorePaths[i][:lastIdx])
	}
	return func(c *gin.Context) {
		path := pointer.String(c.FullPath())
		if *path == "" {
			path = &c.Request.URL.Path
		}
		if skipPaths.Has(*path) {
			c.Next()
			return
		}
		if lastIdx := strings.LastIndex(*path, "/") + 1; lastIdx > 0 {
			if skipPrefixPaths.Has((*path)[:lastIdx]) {
				c.Next()
				return
			}
		}

		// Start timer.
		start := time.Now()

		c.Next()

		reqLatency := time.Since(start)
		if reqLatency > time.Minute {
			reqLatency -= reqLatency % time.Second
		}
		respStatus := c.Writer.Status()
		respSize := "0 B"
		if c.Writer.Written() {
			respSize = humanize.IBytes(uint64(c.Writer.Size()))
		}
		reqClientIP := c.ClientIP()
		reqMethod := c.Request.Method
		switch {
		case IsUnidiStreamRequest(c):
			reqMethod = "US"
		case IsBidiStreamRequest(c):
			reqMethod = "BS"
		}
		reqPath := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			reqPath = reqPath + "?" + raw
		}
		var reqProto = c.Request.Proto
		logger.Debugf("%d | %8s | %10s | %13v | %15s | %-7s %s",
			respStatus,
			reqProto,
			respSize,
			reqLatency,
			reqClientIP,
			reqMethod, reqPath)
	}
}
