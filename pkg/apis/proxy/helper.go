package proxy

import (
	"net/http"
	"net/url"
	"strings"
)

var noForwardHeaders = []string{
	"Connection",
	"Cookie",
	"Host",
}

// removeHeaders removes the headers that should not be forwarded.
func removeHeaders(request *http.Request) {
	for _, h := range noForwardHeaders {
		request.Header.Del(h)
	}
}

// getURL parses the requestPath to url.URL.
// If the requestPath does not have schema, it will append https:// to it and output a valid url.URL or error.
func getURL(requestPath string) (*url.URL, error) {
	if !strings.HasPrefix(requestPath, "http") && !strings.HasPrefix(requestPath, "https") {
		requestPath = "https://" + requestPath
	}

	return url.Parse(requestPath)
}

// isWhiteListDomain checks if the host is in the white list.
func isWhiteListDomain(host string, whiteListDomains ...string) bool {
	for _, domain := range whiteListDomains {
		if domain == host {
			return true
		}

		if strings.HasPrefix(domain, "*") && strings.HasSuffix(host, domain[1:]) {
			return true
		}
	}

	return false
}
