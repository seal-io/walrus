package cache

import (
	"bytes"
	"strings"
)

const (
	activeSuffix   = ":active"
	inactiveSuffix = ":inactive"
)

const sessionKeyPrefix = "session:"

// StoreSessionSubject stores the subject via the given session as the given status.
func StoreSessionSubject(sessionValue, subject string, active bool) {
	if sessionValue == "" {
		return
	}
	suffix := inactiveSuffix
	if active {
		suffix = activeSuffix
	}
	_ = cacher.Set(sessionKeyPrefix+sessionValue, []byte(subject+suffix))
}

// LoadSessionSubject checks the given session is active,
// if the session is active, returns the subject,
// if the session is inactive, returns a none nil subject,
// if the session is not recorded, returns a nil subject.
func LoadSessionSubject(sessionValue string) (*string, bool) {
	sessionValue = strings.TrimSpace(sessionValue)
	if sessionValue != "" {
		bs, _ := cacher.Get(sessionKeyPrefix + sessionValue)
		if len(bs) < 7 {
			return nil, false
		}
		if bytes.Equal(bs[len(bs)-7:], []byte(activeSuffix)) {
			if len(bs[:len(bs)-7]) == 0 {
				return nil, false
			}
			s := string(bs[:len(bs)-7])
			return &s, true
		}
	}
	var s string
	return &s, false
}

// CleanSessionSubject cleans the subject of the given session.
func CleanSessionSubject(sessionValue string) {
	if sessionValue == "" {
		return
	}
	_ = cacher.Delete(sessionKeyPrefix + sessionValue)
}

// CleanSessionSubjects cleans all subjects about session.
func CleanSessionSubjects() {
	it := cacher.Iterator()
	for it.SetNext() {
		e, err := it.Value()
		if err != nil {
			break
		}
		key := e.Key()
		if strings.HasPrefix(key, sessionKeyPrefix) {
			_ = cacher.Delete(key)
		}
	}
}

const tokenKeyPrefix = "token:"

// StoreTokenSubject stores the subject via the given token as the given status.
func StoreTokenSubject(tokenValue, subject string, active bool) {
	if tokenValue == "" {
		return
	}
	suffix := inactiveSuffix
	if active {
		suffix = activeSuffix
	}
	_ = cacher.Set(tokenKeyPrefix+tokenValue, []byte(subject+suffix))
}

// LoadTokenSubject loads the subject via the given token,
// if the token is active, returns the subject,
// if the token is inactive, returns a none nil subject,
// if the token is not recorded, returns a nil subject.
func LoadTokenSubject(tokenValue string) (*string, bool) {
	tokenValue = strings.TrimSpace(tokenValue)
	if tokenValue != "" {
		bs, _ := cacher.Get(tokenKeyPrefix + tokenValue)
		if len(bs) < 7 {
			return nil, false
		}
		if bytes.Equal(bs[len(bs)-7:], []byte(activeSuffix)) {
			if len(bs[:len(bs)-7]) == 0 {
				return nil, false
			}
			s := string(bs[:len(bs)-7])
			return &s, true
		}
	}
	var s string
	return &s, false
}

// CleanTokenSubject cleans the subject of the given token.
func CleanTokenSubject(tokenValue string) {
	if tokenValue == "" {
		return
	}
	_ = cacher.Delete(tokenKeyPrefix + tokenValue)
}

// CleanTokenSubjects cleans all subjects about token.
func CleanTokenSubjects() {
	it := cacher.Iterator()
	for it.SetNext() {
		e, err := it.Value()
		if err != nil {
			break
		}
		key := e.Key()
		if strings.HasPrefix(key, tokenKeyPrefix) {
			_ = cacher.Delete(key)
		}
	}
}
