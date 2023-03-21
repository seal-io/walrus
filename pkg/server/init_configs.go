package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/seal-io/seal/pkg/caches"
	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/cds"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/cache"
	"github.com/seal-io/seal/utils/cryptox"
	"github.com/seal-io/seal/utils/strs"
)

func (r *Server) initConfigs(ctx context.Context, opts initOptions) (err error) {
	err = configureLevelCaches(ctx, opts.CacheDriver)
	if err != nil {
		return
	}

	err = configureDataEncryption(ctx, opts.ModelClient, r.DataSourceDataEncryptAlg, r.DataSourceDataEncryptKey)
	if err != nil {
		return
	}

	// configures casdoor max idle duration.
	casdoor.MaxIdleDurationConfig.Set(r.AuthnSessionMaxIdle)

	return
}

// configureLevelCaches configures a memory cache as level 1,
// a remote cache service as level 2.
func configureLevelCaches(ctx context.Context, drv cds.Driver) error {
	var stack = make([]cache.Cache, 0, 2)

	var processCfg = cache.MemoryConfig{
		EntryMaxAge: 15 * time.Minute,
	}
	var process, err = cache.NewMemoryWithConfig(ctx, processCfg)
	if err != nil {
		return err
	}
	stack = append(stack, process)

	dialect, cli, err := drv.Underlay(ctx)
	if err != nil {
		return err
	}
	switch dialect {
	case "redis", "rediss":
		var remoteCfg = cache.RemoteRedisConfig{
			Namespace:   "seal",
			EntryMaxAge: 30 * time.Minute,
			Client:      cli.(redis.UniversalClient),
		}
		remote, err := cache.NewRemoteRedisWithConfig(ctx, remoteCfg)
		if err != nil {
			return err
		}
		stack = append(stack, remote)
	}

	caches.Stack.Set(stack)
	return nil
}

// configureDataEncryption validates settings.DataEncryptionSentry with data encryption key,
// and enables encrypting data.
func configureDataEncryption(ctx context.Context, mc model.ClientSet, alg string, key []byte) error {
	var enc cryptox.Encryptor
	if key != nil {
		var err error
		switch alg {
		default:
			err = fmt.Errorf("unknown data encryptor algorithm: %s", alg)
		case "aesgcm":
			enc, err = cryptox.AesGcm(key)
		}
		if err != nil {
			return fmt.Errorf("error gaining data encryptor: %w", err)
		}
		crypto.EncryptorConfig.Set(enc)
	} else {
		enc = crypto.EncryptorConfig.Get()
	}

	var pt = "YOU CAN SEE ME"
	return settings.DataEncryptionSentry.Cas(ctx, mc, func(ctb64 string) (string, error) {
		if ctb64 == "" {
			// first time launching, just encrypt.
			var ctbs, err = enc.Encrypt(strs.ToBytes(&pt), nil)
			if err != nil {
				return "", err
			}
			return strs.EncodeBase64(strs.FromBytes(&ctbs)), nil
		}
		// decrypt and compare.
		var ct, _ = strs.DecodeBase64(ctb64)
		var ptbs, _ = enc.Decrypt(strs.ToBytes(&ct), nil)
		if !bytes.Equal(strs.ToBytes(&pt), ptbs) {
			return "", errors.New("inconsistent sentry")
		}
		return ctb64, nil
	})
}
