package options

import (
	"context"
	"crypto/x509"
	"os"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/settings"
)

// GetCertPool returns the certificate pool from system default pool and the given CA file in settings.
func GetCertPool(ctx context.Context, client model.ClientSet) (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	data, err := GetTrustedCAData(ctx, client)
	if err != nil {
		return pool, err
	}

	pool.AppendCertsFromPEM(data)

	return pool, err
}

// GetTrustedCAData returns trusted CA certs data from the given CA file in settings.
func GetTrustedCAData(ctx context.Context, client model.ClientSet) ([]byte, error) {
	file, err := settings.SSLTrustedCAFile.Value(ctx, client)
	if err != nil {
		return nil, err
	}

	if file == "" {
		return nil, nil
	}

	certs, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return certs, nil
}
