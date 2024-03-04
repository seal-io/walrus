package httpx

import (
	"net"
	"net/http"
	"strings"
)

type LoopbackAccessHandlerFunc func(http.ResponseWriter, *http.Request)

func (f LoopbackAccessHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isLoopBackAccessing := func(r *http.Request) bool {
		host := r.Host
		if host == "127.0.0.1" || host == "localhost" || host == "::1" {
			ip, _, _ := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
			return ip == "::1" || host == ip
		}

		return false
	}

	if !isLoopBackAccessing(r) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f(w, r)
}
