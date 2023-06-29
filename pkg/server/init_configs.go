package server

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/seal-io/seal/pkg/auths"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/cryptox"
	"github.com/seal-io/seal/utils/strs"
)

func (r *Server) initConfigs(ctx context.Context, opts initOptions) (err error) {
	err = validateDataEncryption(ctx, opts.ModelClient)
	if err != nil {
		return
	}

	err = configureAuths(ctx, opts.ModelClient, r.AuthnSessionMaxIdle, r.EnableTls)
	if err != nil {
		return
	}

	return
}

// validateDataEncryption validates settings.DataEncryptionSentry with data encryption key.
func validateDataEncryption(ctx context.Context, mc model.ClientSet) error {
	pt := "YOU CAN SEE ME"

	err := settings.DataEncryptionSentry.Cas(ctx, mc, func(ctb64 string) (string, error) {
		enc := crypto.EncryptorConfig.Get()

		// First time launching, just encrypt.
		if ctb64 == "" {
			ctbs, err := enc.Encrypt(strs.ToBytes(&pt), nil)
			if err != nil {
				return "", err
			}

			return strs.EncodeBase64(strs.FromBytes(&ctbs)), nil
		}

		// Otherwise, decrypt and compare the sentry.
		ct, _ := strs.DecodeBase64(ctb64)

		ptbs, _ := enc.Decrypt(strs.ToBytes(&ct), nil)
		if !bytes.Equal(strs.ToBytes(&pt), ptbs) {
			return "", errors.New("inconsistent data encryption sentry")
		}

		return ctb64, nil
	})
	if err != nil {
		// Return clearer message if encryption processing failed.
		if strings.Contains(err.Error(), "cipher:") {
			return errors.New("inconsistent data encryption key")
		}

		return err
	}

	return nil
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
