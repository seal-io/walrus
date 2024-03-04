package fakecert

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"k8s.io/klog/v2"

	"github.com/seal-io/utils/certs"
)

// certKey is the key by which certificates are tracked in state and cache.
type certKey struct {
	// host indicates the host for generating.
	host string
	// alterIPs indicates the alternative IPs for generating,
	// which is separated by comma.
	alterIPs string
	// alterDNSNames indicates the alternative DNS names for generating,
	// which is separated by comma.
	alterDNSNames string
	// rsa indicates to use RSA algorithm, default is to use ECDSA.
	rsa bool
}

func (c certKey) String() string {
	if c.rsa {
		return strings.Join([]string{c.host, c.alterIPs, c.alterDNSNames, "rsa"}, ":")
	}

	return strings.Join([]string{c.host, c.alterIPs, c.alterDNSNames}, ":")
}

// certState is the state by which certificates are ready to read.
type certState struct {
	c       *sync.Cond
	tlsCert *tls.Certificate
	err     error
}

func (s *certState) Done() bool {
	return s.tlsCert != nil || s.err != nil
}

func (s *certState) Generate(
	ctx context.Context,
	caKey crypto.Signer, caCert *x509.Certificate,
	host, alterIPs, alterDNSNames string, isRSA bool,
) {
	if caKey == nil || caCert == nil {
		s.err = errors.New("fakecert: CA is nil")
		return
	}

	s.c.L.Lock()
	defer func() {
		s.c.L.Unlock()
		s.c.Broadcast()
	}()

	var key crypto.Signer
	if isRSA {
		key, s.err = rsa.GenerateKey(rand.Reader, 2048)
	} else {
		key, s.err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if s.err != nil {
		return
	}

	if host == "" {
		s.err = errors.New("fakecert: missing host")
		return
	}

	hosts := []string{host}
	if alterIPs != "" {
		hosts = append(hosts, strings.Split(alterIPs, ",")...)
	}
	if alterDNSNames != "" {
		hosts = append(hosts, strings.Split(alterDNSNames, ",")...)
	}

	var srvCert *x509.Certificate
	srvCert, s.err = genCert(caKey, caCert, key, hosts)
	if s.err != nil {
		return
	}

	s.tlsCert = &tls.Certificate{
		PrivateKey: key,
		Certificate: [][]byte{
			srvCert.Raw,
		},
		Leaf: srvCert,
	}
}

func (s *certState) Get() (*tls.Certificate, error) {
	if !s.Done() {
		s.c.L.Lock()
		for !s.Done() {
			s.c.Wait()
		}
		s.c.L.Unlock()
	}

	return s.tlsCert, s.err
}

// loadOrGenCA loads or generates a CA certificate and private key.
func loadOrGenCA(ctx context.Context, c certs.Cache) (crypto.Signer, *x509.Certificate, error) {
	bs, err := c.Get(ctx, "ca")
	if err != nil {
		if !errors.Is(err, autocert.ErrCacheMiss) {
			return nil, nil, err
		}

		key, x509Cert, err := genCA()
		if err != nil {
			return nil, nil, err
		}

		bs, err = certs.EncodeTLSCertificate(&tls.Certificate{
			PrivateKey: key,
			Certificate: [][]byte{
				x509Cert.Raw,
			},
		})
		if err != nil {
			return nil, nil, err
		}

		err = c.Put(ctx, "ca", bs)
		if err != nil {
			return nil, nil, err
		}

		return key, x509Cert, nil
	}

	tlsCert, err := certs.DecodeTLSCertificate(bs, "")
	if err != nil {
		return nil, nil, err
	}

	return tlsCert.PrivateKey.(crypto.Signer), tlsCert.Leaf, nil
}

// genCA generates a CA certificate and private key.
func genCA() (crypto.Signer, *x509.Certificate, error) {
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
				CommonName:   "fakecert-ca@" + strconv.FormatInt(now.Unix(), 10),
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

	klog.Background().WithName("fakecert").V(5).
		Info("generated",
			"CA", x509Cert.Subject,
			"not before", x509Cert.NotBefore,
			"not after", x509Cert.NotAfter)

	return key, x509Cert, nil
}

// genCert generates a certificate.
func genCert(caKey crypto.Signer, caX509Cert *x509.Certificate, key crypto.Signer, servers []string) (*x509.Certificate, error) {
	var (
		dnsNames    []string
		ipAddresses []net.IP
	)
	for i := range servers {
		if ip := net.ParseIP(servers[i]); ip != nil {
			ipAddresses = append(ipAddresses, ip)
		} else {
			dnsNames = append(dnsNames, servers[i])
		}
	}

	now := time.Now()
	x509Cert := &x509.Certificate{
		DNSNames:     dnsNames,
		IPAddresses:  ipAddresses,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:     time.Now().Add(time.Hour * 24 * 92).UTC(), // 3 months.
		NotBefore:    caX509Cert.NotBefore,
		SerialNumber: new(big.Int).SetInt64(now.Unix()),
		Subject: pkix.Name{
			CommonName:   "fakecert@" + strconv.FormatInt(now.Unix(), 10),
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

	klog.Background().WithName("fakecert").V(5).
		Info("cert generated",
			"servers", servers,
			"certificate", x509Cert.Subject,
			"not before", x509Cert.NotBefore,
			"not after", x509Cert.NotAfter)

	return x509Cert, nil
}
