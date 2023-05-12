package dynacert

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/seal-io/seal/utils/log"
)

func GenCA() (crypto.Signer, *x509.Certificate, error) {
	// Generate private key.
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Generate certificate.
	var (
		now      = time.Now()
		x509Cert = &x509.Certificate{
			BasicConstraintsValid: true,
			IsCA:                  true,
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			NotBefore:             now.Add(-5 * time.Minute).UTC(),
			NotAfter:              now.Add(time.Hour * 24 * 365 * 10).UTC(), // 10 years.
			SerialNumber:          new(big.Int).SetInt64(now.Unix()),
			Subject: pkix.Name{
				CommonName:   "dynacert-ca@" + strconv.FormatInt(now.Unix(), 10),
				Organization: []string{"seal.io"},
			},
		}
	)

	certDER, err := x509.CreateCertificate(rand.Reader, x509Cert, x509Cert, key.Public(), key)
	if err != nil {
		return nil, nil, err
	}

	x509Cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, err
	}

	log.WithName("dynacert").Debugf("generated CA %q, not before %v, not after %v",
		x509Cert.Subject, x509Cert.NotBefore, x509Cert.NotAfter)

	return key, x509Cert, nil
}

func GenServerCert(
	caKey crypto.Signer,
	caX509Cert *x509.Certificate,
	key crypto.Signer,
	server ...string,
) (*x509.Certificate, error) {
	var (
		dnsNames    []string
		ipAddresses []net.IP
	)

	for i := range server {
		if ip := net.ParseIP(server[i]); ip != nil {
			ipAddresses = append(ipAddresses, ip)
		} else {
			dnsNames = append(dnsNames, server[i])
		}
	}

	now := time.Now()
	x509Cert := &x509.Certificate{
		DNSNames:     dnsNames,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  ipAddresses,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:     time.Now().Add(time.Hour * 24 * 92).UTC(), // 3 months.
		NotBefore:    caX509Cert.NotBefore,
		SerialNumber: new(big.Int).SetInt64(now.Unix()),
		Subject: pkix.Name{
			CommonName:   "dynacert@" + strconv.FormatInt(now.Unix(), 10),
			Organization: []string{"seal.io"},
		},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, x509Cert, caX509Cert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}

	x509Cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}

	log.WithName("dynacert").Debugf("generated server %v certificate %q signed by CA %q, not before %v, not after %v",
		server, x509Cert.Subject, caX509Cert.Subject, x509Cert.NotBefore, x509Cert.NotAfter)

	return x509Cert, nil
}

func LoadOrGenSelfSignedCA(ctx context.Context, cache Cache) (crypto.Signer, *x509.Certificate, error) {
	bs, err := cache.Get(ctx, "ca")
	if err != nil {
		if !errors.Is(err, autocert.ErrCacheMiss) {
			return nil, nil, err
		}

		key, x509Cert, err := GenCA()
		if err != nil {
			return nil, nil, err
		}

		bs, err = encodeTlsCertificate(&tls.Certificate{
			PrivateKey: key,
			Certificate: [][]byte{
				x509Cert.Raw,
			},
		})
		if err != nil {
			return nil, nil, err
		}

		err = cache.Put(ctx, "ca", bs)
		if err != nil {
			return nil, nil, err
		}

		return key, x509Cert, nil
	}

	tlsCert, err := decodeTlsCertificate(bs, "")
	if err != nil {
		return nil, nil, err
	}

	return tlsCert.PrivateKey.(crypto.Signer), tlsCert.Leaf, nil
}

func encodeTlsCertificate(tlsCert *tls.Certificate) ([]byte, error) {
	var buf bytes.Buffer

	// Encode private key.
	switch k := tlsCert.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}

		pb := &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
		if err = pem.Encode(&buf, pb); err != nil {
			return nil, err
		}
	case *rsa.PrivateKey:
		b := x509.MarshalPKCS1PrivateKey(k)

		pb := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}
		if err := pem.Encode(&buf, pb); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown private key type")
	}

	// Encode public key.
	for _, b := range tlsCert.Certificate {
		pb := &pem.Block{Type: "CERTIFICATE", Bytes: b}
		if err := pem.Encode(&buf, pb); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func decodeTlsCertificate(data []byte, hostname string) (*tls.Certificate, error) {
	// Decode private key.
	keyBlock, data := pem.Decode(data)
	if keyBlock == nil || !strings.Contains(keyBlock.Type, "PRIVATE") {
		return nil, errors.New("corrupt private key: invalid format")
	}

	var (
		key crypto.PrivateKey
		// Decode public key.
		certs [][]byte
		// Get leaf certificate.
		n int
	)

	key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		// Try RSA.
		k, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
		if err != nil {
			return nil, errors.New("corrupt private key: parse failed")
		}
		key = k.(*rsa.PrivateKey)
	}

	for len(data) > 0 {
		var certBlock *pem.Block

		certBlock, data = pem.Decode(data)
		if certBlock == nil {
			break
		}

		certs = append(certs, certBlock.Bytes)
	}

	if len(data) > 0 {
		return nil, errors.New("corrupt certificate: invalid format")
	}

	for i := range certs {
		n += len(certs[i])
	}
	der := make([]byte, n)
	n = 0

	for i := range certs {
		n += copy(der[n:], certs[i])
	}

	x509Certs, err := x509.ParseCertificates(der)
	if err != nil {
		return nil, err
	}

	if len(x509Certs) == 0 {
		return nil, errors.New("corrupt certificate: not found")
	}
	leaf := x509Certs[0]

	// Validate leaf certificate.
	err = verifyCert(leaf, hostname)
	if err != nil {
		return nil, err
	}

	switch prv := key.(type) {
	case *ecdsa.PrivateKey:
		pub, ok := leaf.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("corrupt certificate: not ECDSA type")
		}

		if pub.X.Cmp(prv.X) != 0 || pub.Y.Cmp(prv.Y) != 0 {
			return nil, errors.New("corrupt private key: not match ECDSA certificate")
		}
	case *rsa.PrivateKey:
		pub, ok := leaf.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("corrupt certificate: not RSA type")
		}

		if pub.N.Cmp(prv.N) != 0 {
			return nil, errors.New("corrupt private key: not match RSA certificate")
		}
	}

	return &tls.Certificate{
		Certificate: certs,
		PrivateKey:  key,
		Leaf:        leaf,
	}, nil
}

func verifyCert(x509Cert *x509.Certificate, hostname string) error {
	now := time.Now()
	if now.Before(x509Cert.NotBefore) {
		return errors.New("corrupt certificate: not valid yet")
	}

	if now.After(x509Cert.NotAfter) {
		return errors.New("corrupt certificate: expired")
	}

	if hostname != "" {
		err := x509Cert.VerifyHostname(hostname)
		if err != nil {
			return fmt.Errorf("corrupt certificate: %w", err)
		}
	}

	return nil
}
