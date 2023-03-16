package casdoor

import (
	"net/http"
	"strings"

	"github.com/seal-io/seal/utils/req"
	"github.com/seal-io/seal/utils/strs"
)

func GetInternalSession(sealSessions []*http.Cookie) *req.HttpCookie {
	var dst *req.HttpCookie
	for i := range sealSessions {
		if sealSessions[i] == nil || sealSessions[i].Name != ExternalSessionCookieKey {
			continue
		}
		dst = &req.HttpCookie{}
		dst.SetKey(InternalSessionCookieKey)
		dst.SetValue(sealSessions[i].Value)
		dst.SetMaxAge(sealSessions[i].MaxAge)
		dst.SetPath("/")
		dst.SetDomain("")
		dst.SetSecure(false) // internal access
		dst.SetHTTPOnly(true)
	}
	return dst
}

func GetExternalSession(casdoorSessions []*req.HttpCookie) *http.Cookie {
	var dst *http.Cookie
	for i := range casdoorSessions {
		if casdoorSessions[i] == nil || string(casdoorSessions[i].Key()) != InternalSessionCookieKey {
			continue
		}
		dst = &http.Cookie{}
		dst.Name = ExternalSessionCookieKey
		dst.Value = string(casdoorSessions[i].Value())
		dst.MaxAge = casdoorSessions[i].MaxAge()
		dst.Path = "/"
		dst.Domain = ""
		dst.Secure = false // TODO
		dst.HttpOnly = true
	}
	return dst
}

func GetInternalToken(sealHeader http.Header) string {
	// get basic auth
	authorization := sealHeader.Get("Authorization")
	if strings.Contains(authorization, "Basic") {
		basicAuth := strings.TrimSpace(strings.TrimPrefix(authorization, "Basic "))
		// decode basic auth
		data, err := strs.DecodeBase64(basicAuth)
		if err != nil {
			return ""
		}
		// get token
		splits := strings.SplitN(string(data), ":", 2)
		if len(splits) != 2 {
			return ""
		}
		return splits[1]
	}

	return strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
}
