package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/cryptox"
	"github.com/seal-io/seal/utils/strs"
)

// initDataEncryption validates settings.DataEncryptionSentry with data encryption key,
// and enables encrypting data.
func (r *Server) initDataEncryption(ctx context.Context, opts initOptions) error {
	var enc cryptox.Encryptor
	if r.DataSourceDataEncryptKey != nil {
		var err error
		switch r.DataSourceDataEncryptAlg {
		default:
			err = fmt.Errorf("unknown data encryptor algorithm: %s", r.DataSourceDataEncryptAlg)
		case "aesgcm":
			enc, err = cryptox.AesGcm(r.DataSourceDataEncryptKey)
		}
		if err != nil {
			return fmt.Errorf("error gaining data encryptor: %w", err)
		}
		crypto.EncryptorConfig.Set(enc)
	} else {
		enc = crypto.EncryptorConfig.Get()
	}

	var pt = "YOU CAN SEE ME"
	return settings.DataEncryptionSentry.Cas(ctx, opts.ModelClient, func(ctb64 string) (string, error) {
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
