package cache

import (
	"bytes"
	"context"
	"strings"

	"github.com/seal-io/seal/utils/cache"
)

const (
	activeSuffix   = ":active"
	inactiveSuffix = ":inactive"
)

const sessionKeyPrefix = "session:"

// StoreSessionSubject stores the subject via the given session as the given status.
func StoreSessionSubject(ctx context.Context, sessionValue, subject string, active bool) {
	if sessionValue == "" {
		return
	}
	var suffix = inactiveSuffix
	if active {
		suffix = activeSuffix
	}
	_ = cacher.Set(ctx, sessionKeyPrefix+sessionValue, []byte(subject+suffix))
}

// LoadSessionSubject checks the given session is active,
// if the session is active, returns the subject,
// if the session is inactive, returns a none nil subject,
// if the session is not recorded, returns a nil subject.
func LoadSessionSubject(ctx context.Context, sessionValue string) (*string, bool) {
	sessionValue = strings.TrimSpace(sessionValue)
	if sessionValue != "" {
		var bs, _ = cacher.Get(ctx, sessionKeyPrefix+sessionValue)
		if len(bs) < 7 {
			return nil, false
		}
		if bytes.Equal(bs[len(bs)-7:], []byte(activeSuffix)) {
			if len(bs[:len(bs)-7]) == 0 {
				return nil, false
			}
			var s = string(bs[:len(bs)-7])
			return &s, true
		}
	}
	var s string
	return &s, false
}

// CleanSessionSubject cleans the subject of the given session.
func CleanSessionSubject(ctx context.Context, sessionValue string) {
	if sessionValue == "" {
		return
	}
	_ = cacher.Delete(ctx, sessionKeyPrefix+sessionValue)
}

// CleanSessionSubjects cleans all subjects about session.
func CleanSessionSubjects(ctx context.Context) {
	_ = cacher.Iterate(ctx, cache.HasPrefix(sessionKeyPrefix),
		func(ctx context.Context, e cache.Entry) (bool, error) {
			_ = cacher.Delete(ctx, e.Key())
			return true, nil
		})
}

const tokenKeyPrefix = "token:"

// StoreTokenSubject stores the subject via the given token as the given status.
func StoreTokenSubject(ctx context.Context, tokenValue, subject string, active bool) {
	if tokenValue == "" {
		return
	}
	var suffix = inactiveSuffix
	if active {
		suffix = activeSuffix
	}
	_ = cacher.Set(ctx, tokenKeyPrefix+tokenValue, []byte(subject+suffix))
}

// LoadTokenSubject loads the subject via the given token,
// if the token is active, returns the subject,
// if the token is inactive, returns a none nil subject,
// if the token is not recorded, returns a nil subject.
func LoadTokenSubject(ctx context.Context, tokenValue string) (*string, bool) {
	tokenValue = strings.TrimSpace(tokenValue)
	if tokenValue != "" {
		var bs, _ = cacher.Get(ctx, tokenKeyPrefix+tokenValue)
		if len(bs) < 7 {
			return nil, false
		}
		if bytes.Equal(bs[len(bs)-7:], []byte(activeSuffix)) {
			if len(bs[:len(bs)-7]) == 0 {
				return nil, false
			}
			var s = string(bs[:len(bs)-7])
			return &s, true
		}
	}
	var s string
	return &s, false
}

// CleanTokenSubject cleans the subject of the given token.
func CleanTokenSubject(ctx context.Context, tokenValue string) {
	if tokenValue == "" {
		return
	}
	_ = cacher.Delete(ctx, tokenKeyPrefix+tokenValue)
}

// CleanTokenSubjects cleans all subjects about token.
func CleanTokenSubjects(ctx context.Context) {
	_ = cacher.Iterate(ctx, cache.HasPrefix(tokenKeyPrefix),
		func(ctx context.Context, e cache.Entry) (bool, error) {
			_ = cacher.Delete(ctx, e.Key())
			return true, nil
		})
}
