package runtime

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// If is a gin middleware,
// which is used for judging the incoming request.
func If(has func(r *http.Request) bool) Handle {
	if has == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {
		if !has(c.Request.Clone(c.Request.Context())) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

// IfLocalIP is a gin middleware,
// which is used for judging the incoming request whether
func IfLocalIP() Handle {
	return If(IsLocalIP)
}

// IsLocalIP returns true if the given request is connected from localhost.
func IsLocalIP(r *http.Request) bool {
	var host = r.Host
	if host == "127.0.0.1" || host == "localhost" || host == "::1" {
		ip, _, _ := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
		return ip == "::1" || host == ip
	}
	return false
}
