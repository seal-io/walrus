package runtime

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Only is a gin middleware,
// which is used for judging the incoming request,
// aborts with 403 if not match.
func Only(match func(*gin.Context) bool) Handle {
	if match == nil {
		return next()
	}
	return func(c *gin.Context) {
		if !match(c) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

// OnlyLocalIP judges the incoming request whether is from localhost,
// aborts with 403 if not match.
func OnlyLocalIP() Handle {
	var isLocalIP = func(c *gin.Context) bool {
		var host = c.Request.Host
		if host == "127.0.0.1" || host == "localhost" || host == "::1" {
			var ip = c.RemoteIP()
			return ip == "::1" || host == ip
		}
		return false
	}
	return Only(isLocalIP)
}

// If is a gin middleware,
// which is used for judging the incoming request,
// execute given handle if matched.
func If(match func(*gin.Context) bool, then Handle) Handle {
	if match == nil || then == nil {
		return next()
	}
	return func(c *gin.Context) {
		if match(c) {
			then(c)
		}
	}
}

// Per is a gin middleware,
// which is used for providing new handler for different incoming request.
func Per(hashRequest func(*gin.Context) string, provideHandler func() Handle) Handle {
	if hashRequest == nil || provideHandler == nil {
		return next()
	}
	var m sync.Map
	return func(c *gin.Context) {
		var k = hashRequest(c)
		var h Handle
		var v, ok = m.LoadOrStore(k, nil)
		if !ok {
			h = provideHandler()
			m.Store(k, h)
		} else {
			h = v.(Handle)
		}
		h(c)
	}
}

// PerIP provides new handler according to incoming request IP.
func PerIP(provideHandler func() Handle) Handle {
	var hashRequestByIP = func(c *gin.Context) string {
		return c.ClientIP()
	}
	return Per(hashRequestByIP, provideHandler)
}

func next() Handle {
	return func(c *gin.Context) { c.Next() }
}
