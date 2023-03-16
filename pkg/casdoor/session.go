package casdoor

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/seal-io/seal/utils/req"
	"github.com/seal-io/seal/utils/strs"
	"github.com/seal-io/seal/utils/vars"
)

// MaxIdleDurationConfig holds the config of the max idle duration.
var MaxIdleDurationConfig = vars.SetOnce[time.Duration]{}

// GetSession converts external session(seal cookie) to internal session(casdoor cookie).
func GetSession(sealSessions []*http.Cookie) *req.HttpCookie {
	var dst *req.HttpCookie
	for i := range sealSessions {
		if sealSessions[i] == nil || sealSessions[i].Name != ExternalSessionCookieKey {
			continue
		}
		var value = sealSessions[i].Value
		if value == "" {
			break
		}
		// request Cookie header
		dst = &req.HttpCookie{}
		dst.SetKey(InternalSessionCookieKey)
		dst.SetValue(value)
		dst.SetPath("/")
		dst.SetDomain("")
		dst.SetSecure(false) // internal access
		dst.SetHTTPOnly(true)
		break
	}
	return dst
}

// GetToken extracts token(casdoor token) from header.
func GetToken(sealHeader http.Header) string {
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

// HoldSession holds the session of authenticated connection.
func HoldSession(w http.ResponseWriter, sessions []*req.HttpCookie) error {
	return manageSession(w, sessions, false)
}

// InterruptSession interrupts the session of authenticated connection.
func InterruptSession(w http.ResponseWriter, sessions []*req.HttpCookie) error {
	return manageSession(w, sessions, true)
}

// manageSession manages the session of authenticated connection.
func manageSession(w http.ResponseWriter, sessions []*req.HttpCookie, interrupt bool) error {
	var s = getExternalSession(sessions)
	if s == nil {
		if interrupt {
			return nil
		}
		return errors.New("cannot get external session")
	}
	if interrupt {
		s.Value = ""
		s.MaxAge = -1
		s.Expires = time.Time{}
	}
	http.SetCookie(w, s)
	return nil
}

// getExternalSession converts internal session(casdoor cookie) to external session(seal cookie).
func getExternalSession(casdoorSessions []*req.HttpCookie) *http.Cookie {
	var dst *http.Cookie
	for i := range casdoorSessions {
		if casdoorSessions[i] == nil || string(casdoorSessions[i].Key()) != InternalSessionCookieKey {
			continue
		}
		var value = string(casdoorSessions[i].Value())
		if value == "" {
			break
		}
		// response Set-Cookie header
		dst = &http.Cookie{}
		dst.Name = ExternalSessionCookieKey
		dst.Value = value
		dst.Path = "/"
		dst.Domain = ""
		dst.Secure = true
		dst.HttpOnly = true
		dst.MaxAge = int(MaxIdleDurationConfig.Get().Round(time.Second) / time.Second)
		if dst.MaxAge > 0 {
			dst.Expires = time.Now().Add(time.Duration(dst.MaxAge) * time.Second)
		}
		break
	}
	return dst
}
