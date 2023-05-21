package server

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/seal-io/seal/pkg/auths"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/cryptox"
	"github.com/seal-io/seal/utils/strs"
)

func (r *Server) initConfigs(ctx context.Context, opts initOptions) (err error) {
	err = configureDataEncryption(ctx, opts.ModelClient, r.DataSourceDataEncryptAlg, r.DataSourceDataEncryptKey)
	if err != nil {
		return
	}

	err = configureAuths(ctx, opts.ModelClient, r.AuthnSessionMaxIdle, r.EnableTls)
	if err != nil {
		return
	}

	return
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

	pt := "YOU CAN SEE ME"

	return settings.DataEncryptionSentry.Cas(ctx, mc, func(ctb64 string) (string, error) {
		if ctb64 == "" {
			// First time launching, just encrypt.
			ctbs, err := enc.Encrypt(strs.ToBytes(&pt), nil)
			if err != nil {
				return "", err
			}

			return strs.EncodeBase64(strs.FromBytes(&ctbs)), nil
		}
		// Decrypt and compare.
		ct, _ := strs.DecodeBase64(ctb64)

		ptbs, _ := enc.Decrypt(strs.ToBytes(&ct), nil)
		if !bytes.Equal(strs.ToBytes(&pt), ptbs) {
			return "", errors.New("inconsistent sentry")
		}

		return ctb64, nil
	})
}

// configureAuths configures the auths.
func configureAuths(ctx context.Context, mc model.ClientSet, maxIdle time.Duration, secure bool) error {
	// Configures cookie max idle duration.
	auths.MaxIdleDurationConfig.Set(maxIdle)

	// Configures securing cookie.
	auths.SecureConfig.Set(secure)

	// Configures token encryptor.
	var enc cryptox.Encryptor
	{
		hexKey, err := settings.AuthsEncryptionAesGcmKey.Value(ctx, mc)
		if err != nil {
			return err
		}
		key, err := hex.DecodeString(hexKey)
		if err != nil {
			return fmt.Errorf("failed to decode hex string: %w", err)
		}
		enc, err = cryptox.AesGcm(key)
		if err != nil {
			return fmt.Errorf("failed to create aes-gcm encryptor: %w", err)
		}
	}
	auths.EncryptorConfig.Set(enc)

	return nil
}
