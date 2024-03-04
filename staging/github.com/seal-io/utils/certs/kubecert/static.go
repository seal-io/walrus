package kubecert

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/idna"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	"k8s.io/klog/v2"

	"github.com/seal-io/utils/certs"
	"github.com/seal-io/utils/stringx"
)

// StaticManager manages the Kubernetes Certificate Signing Request certificates,
// which generates the server certificate according to the given Host, IPs, and DNS names.
type StaticManager struct {
	// CertCli indicates the client for generating certificates.
	// If nil, the StaticManager will not attempt to generate new certificates.
	CertCli certs.CertificateSigningRequestInterface

	// Cache optionally stores and retrieves previously-obtained certificates
	// and other state. If nil, certs will only be cached for the lifetime of
	// the StaticManager.
	//
	// Multiple StaticManager instances can share the same Cache.
	Cache certs.Cache

	// Host indicates the host for generating certificates.
	// If blank, the StaticManager will not attempt to generate new certificates.
	Host string

	// AlternateIPs indicates the alternate IPs for generating certificates.
	AlternateIPs []net.IP

	// AlternateDNSNames indicates the alternate DNS names for generating certificates.
	AlternateDNSNames []string

	state sync.Map
}

// GetCertificate implements the tls.Config.GetCertificate hook.
// It only provides a TLS certificate for the StaticManager's Host, IPs, and DNS names.
func (m *StaticManager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	var (
		host   = hello.ServerName
		hostIP string
	)
	if host != "" {
		n, err := idna.Lookup.ToASCII(host)
		if err != nil {
			return nil, errors.New("kubecert: host contains invalid character")
		}
		host = strings.TrimSuffix(n, ".")
	} else {
		addr := hello.Conn.LocalAddr()
		if addr == nil {
			return nil, errors.New("kubecert: missing host")
		}
		tcpAddr, ok := addr.(*net.TCPAddr)
		if !ok {
			return nil, errors.New("kubecert: invalid host")
		}
		hostIP = tcpAddr.IP.String()
	}
	if host == "" {
		host = hostIP
		hostIP = ""
	}

	if m.Host != host && m.Host != hostIP {
		return nil, fmt.Errorf("kubecert: disallowed host: %s", host)
	}
	if !certs.RequestECDSA(hello) {
		return nil, errors.New("kubecert: disallowed ECDSA")
	}

	// Timeout needs to process the worst-case scenario.
	ctx, cancel := context.WithTimeout(hello.Context(), 30*time.Second)
	defer cancel()

	// Get certificate by server.
	ck := certKey{
		host:          m.Host,
		alterIPs:      strings.Join(stringx.Strings(m.AlternateIPs), ","),
		alterDNSNames: strings.Join(m.AlternateDNSNames, ","),
	}

	tlsCert, err := m.getCert(ctx, ck)
	if err == nil {
		return tlsCert, nil
	}

	if !errors.Is(err, certs.ErrCacheMiss) {
		return nil, fmt.Errorf("kubecert: get cert: %w", err)
	}

	// Create certificate to server.
	tlsCert, err = m.createCert(ctx, ck)
	if err != nil {
		return nil, fmt.Errorf("kubecert: create cert: %w", err)
	}

	return tlsCert, nil
}

func (m *StaticManager) Name() string {
	return "kubecert.static"
}

func (m *StaticManager) AddListener(_ dynamiccertificates.Listener) {}

func (m *StaticManager) CurrentCertKeyContent() ([]byte, []byte) {
	logger := klog.Background().WithName("kubecert").WithName("static")

	// Timeout needs to process the worst-case scenario.
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	ck := certKey{
		host:          m.Host,
		alterIPs:      strings.Join(stringx.Strings(m.AlternateIPs), ","),
		alterDNSNames: strings.Join(m.AlternateDNSNames, ","),
	}

	tlsCert, err := m.getCert(ctx, ck)
	if err != nil {
		if !errors.Is(err, certs.ErrCacheMiss) {
			logger.Error(err, "get cert")
			return nil, nil
		}

		// Create certificate to server.
		tlsCert, err = m.createCert(ctx, ck)
		if err != nil {
			logger.Error(err, "create cert")
			return nil, nil
		}
	}

	certPEM, keyPEM, err := certs.EncodeTLSCertificateAndKey(tlsCert)
	if err != nil {
		logger.Error(err, "encode tls certificate and key")
	}

	return certPEM, keyPEM
}

// createCert starts the server ownership verification and returns a certificate
// for that server upon success.
//
// If the server is already being verified, it waits for the existing verification to complete.
// Either way, createCert blocks for the duration of the whole process.
func (m *StaticManager) createCert(ctx context.Context, ck certKey) (*tls.Certificate, error) {
	s := &certState{c: sync.NewCond(&sync.Mutex{})}
	if v, ok := m.state.LoadOrStore(ck, s); ok {
		s = v.(*certState)
	} else {
		s.Generate(ctx, m.CertCli, ck.host, ck.alterIPs, ck.alterDNSNames, ck.rsa)
	}

	tlsCert, err := s.Get()
	if err != nil {
		return nil, err
	}

	m.cachePut(ctx, ck, tlsCert)

	return tlsCert, nil
}

// getCert returns an existing certificate either from m.state or cache,
// if a certificate is found in cache but not in m.state, the latter will be filled
// with the cached value.
func (m *StaticManager) getCert(ctx context.Context, ck certKey) (*tls.Certificate, error) {
	logger := klog.Background().WithName("kubecert").WithName("dynamic")

	if v, ok := m.state.Load(ck); ok {
		s := v.(*certState)

		tlsCert, err := s.Get()
		if err != nil {
			return nil, err
		}

		err = certs.VerifyX509Certificate(tlsCert.Leaf, ck.host)
		if err != nil {
			logger.V(5).
				Error(err, "verify certificate")
			// Treat as miss cache,
			// so the GetCertificate will regenerate.
			return nil, certs.ErrCacheMiss
		}

		return tlsCert, nil
	}

	tlsCert, err := m.cacheGet(ctx, ck)
	if err != nil {
		return nil, err
	}

	logger.V(5).
		Info("loaded certificate from cache", "host", ck.host)
	m.state.Store(ck, &certState{tlsCert: tlsCert})

	return tlsCert, nil
}

// cacheGet loads from the cache,
// and decodes private key and certificate.
func (m *StaticManager) cacheGet(ctx context.Context, ck certKey) (*tls.Certificate, error) {
	logger := klog.Background().WithName("kubecert").WithName("dynamic")

	if m.Cache == nil {
		return nil, certs.ErrCacheMiss
	}

	bs, err := m.Cache.Get(ctx, ck.String())
	if err != nil {
		return nil, err
	}

	tlsCert, err := certs.DecodeTLSCertificate(bs, ck.host)
	if err != nil {
		logger.V(5).
			Error(err, "decode tls certificate")
		// Treat as miss cache,
		// so the GetCertificate will regenerate.
		return nil, certs.ErrCacheMiss
	}

	return tlsCert, nil
}

// cachePut encodes private key and certificate together,
// and stores to the cache.
func (m *StaticManager) cachePut(ctx context.Context, ck certKey, tlsCert *tls.Certificate) {
	logger := klog.Background().WithName("kubecert").WithName("dynamic")

	if m.Cache == nil {
		return
	}

	bs, err := certs.EncodeTLSCertificate(tlsCert)
	if err != nil {
		logger.V(5).
			Error(err, "encode tls certificate")
		return
	}

	err = m.Cache.Put(ctx, ck.String(), bs)
	if err != nil {
		logger.V(5).
			Error(err, "cache tls certificate bytes")
	}
}
