package builtin

import (
	"context"

	"github.com/seal-io/walrus/pkg/caches"
	"github.com/seal-io/walrus/utils/json"
)

const cacheKeyPrefix = "auths/authn/builtin/"

type cacheValue struct {
	Domain string   `json:"domain"`
	Groups []string `json:"groups"`
	Name   string   `json:"name"`
}

func getCached(ctx context.Context, sv string) (domain string, groups []string, name string, exist bool) {
	bs, _ := caches.Get(ctx, cacheKeyPrefix+sv)
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

func cache(ctx context.Context, sv, domain string, groups []string, name string) {
	v := cacheValue{
		Domain: domain,
		Groups: groups,
		Name:   name,
	}

	bs, err := json.Marshal(v)
	if err != nil {
		return
	}

	_ = caches.Set(ctx, cacheKeyPrefix+sv, bs)
}

func delCached(ctx context.Context, sv string) {
	_ = caches.Delete(ctx, cacheKeyPrefix+sv)
}
