package cache

import (
	"context"
	"time"

	"github.com/seal-io/seal/utils/cache"
)

// cacher keeps authn using process cache.
var cacher = cache.MustNewMemoryWithConfig(context.Background(),
	cache.MemoryConfig{
		EntryMaxAge: 5 * time.Minute,
	})
