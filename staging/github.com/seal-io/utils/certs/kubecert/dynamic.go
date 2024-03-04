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
	"k8s.io/klog/v2"

	"github.com/seal-io/utils/certs"
)

// DynamicManager manages the Kubernetes Certificate Signing Request certificates dynamically,
// which generates the server certificate according to the hello information.
type DynamicManager struct {
	// CertCli indicates the client for generating certificates.
	// If nil, the DynamicManager will not attempt to generate new certificates.
	CertCli certs.CertificateSigningRequestInterface

	// Cache optionally stores and retrieves previously-obtained certificates
	// and other state. If nil, certs will only be cached for the lifetime of
	// the DynamicManager.
	//
	// Multiple DynamicManager instances can share the same Cache.
	Cache certs.Cache

	// HostPolicy controls which domains the DynamicManager will attempt
	// to retrieve new certificates for. It does not affect cached certs.
	//
	// If non-nil, HostPolicy is called before requesting a new getCert.
	// If nil, all hosts are currently allowed.
	HostPolicy certs.HostPolicy

	state sync.Map
}

// GetCertificate implements the tls.Config.GetCertificate hook.
// It provides a TLS certificate for hello.ServerName host.
// All other fields of hello are ignored.
//
// If m.HostPolicy is non-nil, GetCertificate calls the policy before requesting
// a new getCert. A non-nil error returned from m.HostPolicy halts TLS negotiation.
// The error is propagated back to the caller of GetCertificate and is user-visible.
// This does not affect cached certs.
func (m *DynamicManager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
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

	// Timeout needs to process the worst-case scenario.
	ctx, cancel := context.WithTimeout(hello.Context(), 30*time.Second)
	defer cancel()

	// Get certificate by server.
	ck := certKey{
		host:     host,
		alterIPs: hostIP,
		rsa:      !certs.RequestECDSA(hello),
	}

	tlsCert, err := m.getCert(ctx, ck)
	if err == nil {
		return tlsCert, nil
	}

	if !errors.Is(err, certs.ErrCacheMiss) {
		return nil, fmt.Errorf("kubecert: get cert: %w", err)
	}

	// Create certificate to server.
	err = m.allowHost(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("kubecert: disallowed host: %s", host)
	}

	tlsCert, err = m.createCert(ctx, ck)
	if err != nil {
		return nil, fmt.Errorf("kubecert: create cert: %w", err)
	}

	return tlsCert, nil
}

// createCert starts the server ownership verification and returns a certificate
// for that server upon success.
//
// If the server is already being verified, it waits for the existing verification to complete.
// Either way, createCert blocks for the duration of the whole process.
func (m *DynamicManager) createCert(ctx context.Context, ck certKey) (*tls.Certificate, error) {
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
func (m *DynamicManager) getCert(ctx context.Context, ck certKey) (*tls.Certificate, error) {
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
func (m *DynamicManager) cacheGet(ctx context.Context, ck certKey) (*tls.Certificate, error) {
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
func (m *DynamicManager) cachePut(ctx context.Context, ck certKey, tlsCert *tls.Certificate) {
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

// allowHost returns true if the host is allowed.
func (m *DynamicManager) allowHost(ctx context.Context, hostname string) error {
	if m.HostPolicy != nil {
		return m.HostPolicy(ctx, hostname)
	}

	return nil
}
