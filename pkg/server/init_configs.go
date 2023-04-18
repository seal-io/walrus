package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/cryptox"
	"github.com/seal-io/seal/utils/strs"
)

func (r *Server) initConfigs(ctx context.Context, opts initOptions) error {
	var err = configureDataEncryption(ctx, opts.ModelClient, r.DataSourceDataEncryptAlg, r.DataSourceDataEncryptKey)
	if err != nil {
		return err
	}

	// configures casdoor cookie max idle duration.
	casdoor.MaxIdleDurationConfig.Set(r.AuthnSessionMaxIdle)
	// configures casdoor securing cookie.
	casdoor.SecureConfig.Set(r.EnableTls)

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
