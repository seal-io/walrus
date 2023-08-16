package casdoor

import (
	"github.com/seal-io/walrus/utils/req"
)

const sessionCookieName = "casdoor_session_id"

// UnwrapSession returns the value of casdoor login succeeded session.
func UnwrapSession(sessions []*req.HttpCookie) string {
	for i := range sessions {
		if sessions[i] == nil ||
			string(sessions[i].Key()) != sessionCookieName {
			continue
		}

		return string(sessions[i].Value())
	}

	return ""
}

// WrapSession wraps the value as casdoor login succeeded session.
func WrapSession(value string) []*req.HttpCookie {
	// Request Cookie header.
	var s req.HttpCookie

	s.SetKey(sessionCookieName)
	s.SetValue(value)
	s.SetPath("/")
	s.SetDomain("")
	s.SetSecure(false) // Internal access.
	s.SetHTTPOnly(true)

	return []*req.HttpCookie{&s}
}
