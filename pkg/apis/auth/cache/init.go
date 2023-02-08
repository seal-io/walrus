package cache

import (
	"context"

	"github.com/seal-io/seal/utils/cache"
)

var cacher = cache.MustNew(context.Background())
