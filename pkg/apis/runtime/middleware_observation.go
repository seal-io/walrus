package runtime

import (
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/utils/log"
)

// SkipLoggingPaths is a RouterOption to ignore logging for the given paths.
func SkipLoggingPaths(paths ...string) RouterOption {
	return routerOption(func(r *Router) {
		for i := range paths {
			skipLoggingPath(paths[i])
		}
	})
}

var (
	pathSkipLogging       = sets.New[string]()
	pathPrefixSkipLogging = sets.New[string]()
)

func skipLoggingPath(p string) {
	pathSkipLogging.Insert(p)

	lastIdx := strings.LastIndex(p, "/") + 1
	if lastIdx <= 0 {
		return
	}

	lastSeg := p[lastIdx:]
	if !strings.HasPrefix(lastSeg, "*") {
		return
	}

	pathPrefixSkipLogging.Insert(p[:lastIdx])
}

// observing is a gin middleware,
// it measures the http request/response with logging and monitoring.
func observing(c *gin.Context) {
	logger := log.WithName("api")

	// Validate to skip logging or not.
	skipLogging := !logger.Enabled(log.DebugLevel)

	reqPath := c.FullPath()
	if reqPath == "" {
		reqPath = c.Request.URL.Path
	}

	if pathSkipLogging.Has(reqPath) {
		skipLogging = true
	} else if i := strings.LastIndex(reqPath, "/") + 1; i > 0 {
		if pathPrefixSkipLogging.Has(reqPath[:i]) {
			skipLogging = true
		}
	}

	reqProto := c.Request.Proto
	reqMethod := c.Request.Method

	switch {
	case IsUnidiStreamRequest(c):
		reqMethod = "US"
	case IsBidiStreamRequest(c):
		reqMethod = "BS"
	}

	// Record inflight request.
	_statsCollector.requestInflight.
		WithLabelValues(reqProto, reqPath, reqMethod).
		Inc()

	defer func() {
		_statsCollector.requestInflight.
			WithLabelValues(reqProto, reqPath, reqMethod).
			Dec()
	}()

	start := time.Now()

	c.Next()

	reqLatency := time.Since(start)

	reqSize := c.Request.ContentLength
	if v := c.GetInt64("request_size"); v != 0 {
		reqSize = v
	}

	respStatus := strconv.Itoa(c.Writer.Status())
	if v := c.GetInt("response_status"); v != 0 {
		respStatus = strconv.Itoa(v)
	}

	var respSize int64
	if c.Writer.Written() {
		respSize = int64(c.Writer.Size())
	} else if v := c.GetInt64("response_size"); v != 0 {
		respSize = v
	}

	// Record request latency.
	_statsCollector.requestDurations.
		WithLabelValues(reqProto, reqPath, reqMethod, respStatus).
		Observe(reqLatency.Seconds())

	// Record request time.
	_statsCollector.requestCounter.
		WithLabelValues(reqProto, reqPath, reqMethod, respStatus).
		Inc()

	// Record request size.
	_statsCollector.requestSizes.
		WithLabelValues(reqProto, reqPath, reqMethod).
		Observe(float64(reqSize))

	// Record response size.
	_statsCollector.responseSizes.
		WithLabelValues(reqProto, reqPath, reqMethod).
		Observe(float64(respSize))

	if !skipLogging {
		// Complete logging info.
		reqSize := humanize.IBytes(uint64(reqSize))
		respSize := humanize.IBytes(uint64(respSize))

		reqLatency := reqLatency
		if reqLatency > time.Minute {
			reqLatency -= reqLatency % time.Second
		}

		reqClientIP := c.ClientIP()

		reqPath := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			reqPath = reqPath + "?" + raw
		}

		logger.Debugf("%s | %8s | %10s | %10s | %13v | %15s | %-7s %s",
			respStatus,
			reqProto,
			reqSize,
			respSize,
			reqLatency,
			reqClientIP,
			reqMethod,
			reqPath,
		)
	}
}
