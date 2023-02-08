package dynamicert

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rancher/dynamiclistener"
	"golang.org/x/crypto/acme/autocert"
	core "k8s.io/api/core/v1"

	"github.com/seal-io/seal/utils/json"
)

type Cache = autocert.Cache

type DirCache = autocert.DirCache

func cacheToStorage(cache Cache) dynamiclistener.TLSStorage {
	return storage{cache: cache}
}

type storage struct {
	cache Cache
}

func (s storage) Get() (*core.Secret, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var bs, err = s.cache.Get(ctx, "secret")
	if err != nil {
		if errors.Is(err, autocert.ErrCacheMiss) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting secret from cache: %w", err)
	}
	var secret core.Secret
	err = json.Unmarshal(bs, &secret)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling secret: %w", err)
	}
	return &secret, nil
}

func (s storage) Update(secret *core.Secret) error {
	if secret == nil {
		return nil
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var bs, err = json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("error marshalling secret: %w", err)
	}
	err = s.cache.Put(ctx, "secret", bs)
	if err != nil {
		return fmt.Errorf("error caching secret: %w", err)
	}

	err = s.cache.Put(ctx, core.TLSCertKey, secret.Data[core.TLSCertKey])
	if err != nil {
		return fmt.Errorf("error caching %s: %w", core.TLSCertKey, err)
	}
	err = s.cache.Put(ctx, core.TLSPrivateKeyKey, secret.Data[core.TLSPrivateKeyKey])
	if err != nil {
		return fmt.Errorf("error caching %s: %w", core.TLSPrivateKeyKey, err)
	}
	return nil
}
