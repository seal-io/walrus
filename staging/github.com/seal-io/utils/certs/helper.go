package certs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/seal-io/utils/pools/bytespool"
)

// EncodeTLSCertificate encodes a TLS certificate to PEM format,
// include private key and public key.
func EncodeTLSCertificate(cert *tls.Certificate) ([]byte, error) {
	buf := bytespool.GetBuffer()

	// Encode private key.
	switch k := cert.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}

		pb := &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
		if err = pem.Encode(buf, pb); err != nil {
			return nil, err
		}
	case *rsa.PrivateKey:
		b := x509.MarshalPKCS1PrivateKey(k)

		pb := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}
		if err := pem.Encode(buf, pb); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown private key type")
	}

	// Encode public key.
	for _, b := range cert.Certificate {
		pb := &pem.Block{Type: "CERTIFICATE", Bytes: b}
		if err := pem.Encode(buf, pb); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// DecodeTLSCertificate decodes a TLS certificate from PEM format,
// include private key and public key.
func DecodeTLSCertificate(data []byte, hostname string) (*tls.Certificate, error) {
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
	err = VerifyX509Certificate(leaf, hostname)
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

// VerifyX509Certificate verifies the x509 certificate.
func VerifyX509Certificate(cert *x509.Certificate, hostname string) error {
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return errors.New("corrupt certificate: not valid yet")
	}

	if now.After(cert.NotAfter) {
		return errors.New("corrupt certificate: expired")
	}

	if hostname != "" {
		err := cert.VerifyHostname(hostname)
		if err != nil {
			return fmt.Errorf("corrupt certificate: %w", err)
		}
	}

	return nil
}

// EncodeX509CertificateRequest encodes a x509 certificate request to PEM format.
func EncodeX509CertificateRequest(certReq *x509.CertificateRequest, key crypto.Signer) ([]byte, error) {
	b, err := x509.CreateCertificateRequest(rand.Reader, certReq, key)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: b}), nil
}

// DecodeX509Certificate decodes a x509 certificate from PEM format.
func DecodeX509Certificate(data []byte, hostname string) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("corrupt certificate: invalid format")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	err = VerifyX509Certificate(cert, hostname)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// RequestECDSA returns true if requesting ECDSA algorithms.
func RequestECDSA(hello *tls.ClientHelloInfo) bool {
	if hello.SignatureSchemes != nil {
		var ecdsaOK bool
	schemeLoop:
		for _, scheme := range hello.SignatureSchemes {
			switch scheme {
			case tls.ECDSAWithSHA1, tls.ECDSAWithP256AndSHA256,
				tls.ECDSAWithP384AndSHA384, tls.ECDSAWithP521AndSHA512:
				ecdsaOK = true
				break schemeLoop
			}
		}

		if !ecdsaOK {
			return false
		}
	}

	if hello.SupportedCurves != nil {
		var ecdsaOK bool

		for _, curve := range hello.SupportedCurves {
			if curve == tls.CurveP256 {
				ecdsaOK = true
				break
			}
		}

		if !ecdsaOK {
			return false
		}
	}

	for _, suite := range hello.CipherSuites {
		switch suite {
		case tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:
			return true
		}
	}

	return false
}

// EncodeTLSCertificateAndKey is similar to EncodeTLSCertificate,
// but returns private key and certificate separately.
func EncodeTLSCertificateAndKey(cert *tls.Certificate) ([]byte, []byte, error) {
	// Encode private key.
	priBuf := bytespool.GetBuffer()
	switch k := cert.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, nil, err
		}

		pb := &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
		if err = pem.Encode(priBuf, pb); err != nil {
			return nil, nil, err
		}
	case *rsa.PrivateKey:
		b := x509.MarshalPKCS1PrivateKey(k)

		pb := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}
		if err := pem.Encode(priBuf, pb); err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, errors.New("unknown private key type")
	}

	// Encode public key.
	pubBuf := bytespool.GetBuffer()
	for _, b := range cert.Certificate {
		pb := &pem.Block{Type: "CERTIFICATE", Bytes: b}
		if err := pem.Encode(pubBuf, pb); err != nil {
			return nil, nil, err
		}
	}

	return pubBuf.Bytes(), priBuf.Bytes(), nil
}
