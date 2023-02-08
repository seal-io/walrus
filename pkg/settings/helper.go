package settings

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type urlSchemeChecker func(scheme string) (allowed bool)

func anySchemeUrl(scheme string) bool {
	return true
}

func httpSchemeUrlOnly(scheme string) bool {
	switch strings.ToLower(scheme) {
	case "http", "https":
		return true
	}
	return false
}

func sockSchemeUrlOnly(scheme string) bool {
	switch strings.ToLower(scheme) {
	case "socks4", "socks5":
		return true
	}
	return false
}

// parseUrl parses the given string as *url.URL.
func parseUrl(str string, check urlSchemeChecker) (*url.URL, error) {
	var v, err = url.Parse(str)
	if err != nil {
		return nil, fmt.Errorf("%s is illegal URL format: %w", str, err)
	}
	var scheme = v.Scheme
	if check != nil {
		if !check(scheme) {
			return nil, fmt.Errorf("invalid scheme: %s", scheme)
		}
	}
	var port = v.Port()
	if port != "" {
		p, err := strconv.ParseUint(port, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing given port: %w", err)
		}
		if p > 65535 {
			return nil, fmt.Errorf("error given port %d: exceeded upper limit", p)
		}
	}
	return v, nil
}
