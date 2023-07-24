package token

import (
	"context"
	"encoding/json"

	tokenbus "github.com/seal-io/seal/pkg/bus/token"
	"github.com/seal-io/seal/pkg/caches"
)

const cacheKeyPrefix = "auths/authn/token/"

type cacheValue struct {
	Domain string   `json:"domain"`
	Groups []string `json:"groups"`
	Name   string   `json:"name"`
}

func getCached(ctx context.Context, tv string) (domain string, groups []string, name string, exist bool) {
	bs, _ := caches.Get(ctx, cacheKeyPrefix+tv)
	if bs == nil {
		return
	}

	var v cacheValue

	err := json.Unmarshal(bs, &v)
	if err != nil {
		return
	}

	return v.Domain, v.Groups, v.Name, true
}

func cache(ctx context.Context, tv, domain string, groups []string, name string) {
	v := cacheValue{
		Domain: domain,
		Groups: groups,
		Name:   name,
	}

	bs, err := json.Marshal(v)
	if err != nil {
		return
	}

	_ = caches.Set(ctx, cacheKeyPrefix+tv, bs)
}

// DelCached observes the model.Token deletions.
func DelCached(ctx context.Context, m tokenbus.BusMessage) error {
	for i := range m.Refers {
		tv := string(m.Refers[i].Value)
		if tv == "" {
			continue
		}

		_ = caches.Delete(ctx, cacheKeyPrefix+tv)
	}

	return nil
}
