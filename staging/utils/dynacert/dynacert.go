package dynacert

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/idna"

	"github.com/seal-io/seal/utils/log"
)

// HostPolicy specifies which host names the Manager is allowed to respond to.
// It returns a non-nil error if the host should be rejected.
// The returned error is accessible via tls.Conn.Handshake and its callers.
type HostPolicy = autocert.HostPolicy

// HostWhitelist returns a policy where only the specified host names are allowed.
// Only exact matches are currently supported. Subdomains, regexp or wildcard
// will not match.
//
// Note that all hosts will be converted to Punycode via idna.Lookup.ToASCII so that
// Manager.GetCertificate can handle the Unicode IDN and mixedcase hosts correctly.
// Invalid hosts will be silently ignored.
func HostWhitelist(hosts ...string) HostPolicy {
	return autocert.HostWhitelist(hosts...)
}

// Manager manages the self-signed server certificates.
type Manager struct {
	// Cache optionally stores and retrieves previously-obtained certificates
	// and other state. If nil, certs will only be cached for the lifetime of
	// the Manager. Multiple Managers can share the same Cache.
	Cache Cache

	// HostPolicy controls which domains the Manager will attempt
	// to retrieve new certificates for. It does not affect cached certs.
	//
	// If non-nil, HostPolicy is called before requesting a new getCert.
	// If nil, all hosts are currently allowed.
	HostPolicy HostPolicy

	state sync.Map
}

// certKey is the key by which certificates are tracked in state and cache.
type certKey struct {
	// Server indicates the server for generating.
	server string
	// Rsa indicates to use RSA algorithm, default is to use ECDSA.
	rsa bool
}

func (c certKey) String() string {
	if c.rsa {
		return c.server + "+rsa"
	}
	return c.server
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

func (s *certState) Generate(ctx context.Context, cache Cache, server string, isRSA bool) {
	s.c.L.Lock()
	defer func() {
		s.c.L.Unlock()
		s.c.Broadcast()
	}()

	var (
		caKey      crypto.Signer
		caX509Cert *x509.Certificate
	)
	caKey, caX509Cert, s.err = LoadOrGenSelfSignedCA(ctx, cache)
	if s.err != nil {
		return
	}

	var key crypto.Signer
	if isRSA {
		key, s.err = rsa.GenerateKey(rand.Reader, 2048)
	} else {
		key, s.err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if s.err != nil {
		return
	}

	var x509Cert *x509.Certificate
	x509Cert, s.err = GenServerCert(caKey, caX509Cert, key, server)
	if s.err != nil {
		return
	}

	s.tlsCert = &tls.Certificate{
		PrivateKey: key,
		Certificate: [][]byte{
			x509Cert.Raw,
		},
		Leaf: x509Cert,
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

// GetCertificate implements the tls.Config.GetCertificate hook.
// It provides a TLS certificate for hello.ServerName host.
// All other fields of hello are ignored.
//
// If m.HostPolicy is non-nil, GetCertificate calls the policy before requesting
// a new getCert. A non-nil error returned from m.HostPolicy halts TLS negotiation.
// The error is propagated back to the caller of GetCertificate and is user-visible.
// This does not affect cached certs.
func (m *Manager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	name := hello.ServerName
	if name != "" {
		// Allow localhost hostname.
		if name != "localhost" && !strings.Contains(strings.Trim(name, "."), ".") {
			return nil, errors.New("dynacert: server name component count invalid")
		}
		// Validate invalid character.
		var err error
		name, err = idna.Lookup.ToASCII(name)
		if err != nil {
			return nil, errors.New("dynacert: server name contains invalid character")
		}
	} else {
		addr := hello.Conn.LocalAddr()
		if addr == nil {
			return nil, errors.New("dynacert: missing local address")
		}
		tcpAddr, ok := addr.(*net.TCPAddr)
		if !ok {
			return nil, errors.New("dynacert: invalid local address")
		}
		name = tcpAddr.IP.String()
	}

	// Timeout needs to process the worst-case scenario.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Get certificate by server.
	ck := certKey{
		server: strings.TrimSuffix(name, "."),
		rsa:    !requestECDSA(hello),
	}
	tlsCert, err := m.getCert(ctx, ck)
	if err == nil {
		return tlsCert, nil
	}
	if !errors.Is(err, autocert.ErrCacheMiss) {
		return nil, fmt.Errorf("dynacert: error getting cert: %w", err)
	}

	// Create certificate to server.
	err = m.allowHost(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("dynacert: disallowed host: %s", name)
	}
	tlsCert, err = m.createCert(ctx, ck)
	if err != nil {
		return nil, fmt.Errorf("dynacert: error creating cert: %w", err)
	}
	return tlsCert, nil
}

// createCert starts the server ownership verification and returns a certificate
// for that server upon success.
//
// If the server is already being verified, it waits for the existing verification to complete.
// Either way, createCert blocks for the duration of the whole process.
func (m *Manager) createCert(ctx context.Context, ck certKey) (*tls.Certificate, error) {
	s := &certState{c: sync.NewCond(&sync.Mutex{})}
	if v, ok := m.state.LoadOrStore(ck, s); ok {
		s = v.(*certState)
	} else {
		s.Generate(ctx, m.Cache, ck.server, ck.rsa)
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
func (m *Manager) getCert(ctx context.Context, ck certKey) (*tls.Certificate, error) {
	if v, ok := m.state.Load(ck); ok {
		s := v.(*certState)
		tlsCert, err := s.Get()
		if err != nil {
			return nil, err
		}
		err = verifyCert(tlsCert.Leaf, ck.server)
		if err != nil {
			log.WithName("dynacert").Warn(err)
			// Treat as miss cache,
			// so the GetCertificate will regenerate.
			return nil, autocert.ErrCacheMiss
		}
		return tlsCert, nil
	}
	tlsCert, err := m.cacheGet(ctx, ck)
	if err != nil {
		return nil, err
	}
	log.WithName("dynacert").Debugf("loaded %q certificate from cache", ck.server)
	m.state.Store(ck, &certState{tlsCert: tlsCert})
	return tlsCert, nil
}

// cacheGet loads from the cache,
// and decodes private key and certificate.
func (m *Manager) cacheGet(ctx context.Context, ck certKey) (*tls.Certificate, error) {
	if m.Cache == nil {
		return nil, autocert.ErrCacheMiss
	}

	bs, err := m.Cache.Get(ctx, ck.String())
	if err != nil {
		return nil, err
	}
	tlsCert, err := decodeTlsCertificate(bs, ck.server)
	if err != nil {
		log.WithName("dynacert").Warnf("error decoding tls certificate: %v", err)
		// Treat as miss cache,
		// so the GetCertificate will regenerate.
		return nil, autocert.ErrCacheMiss
	}
	return tlsCert, nil
}

// cachePut encodes private key and certificate together,
// and stores to the cache.
func (m *Manager) cachePut(ctx context.Context, ck certKey, tlsCert *tls.Certificate) {
	if m.Cache == nil {
		return
	}

	bs, err := encodeTlsCertificate(tlsCert)
	if err != nil {
		log.WithName("dynacert").Warnf("error encoding tls certificate: %v", err)
		return
	}
	err = m.Cache.Put(ctx, ck.String(), bs)
	if err != nil {
		log.WithName("dynacert").Warnf("error caching tls certificate bytes: %v", err)
	}
}

// allowHost returns true if the host is allowed.
func (m *Manager) allowHost(ctx context.Context, hostname string) error {
	if m.HostPolicy != nil {
		return m.HostPolicy(ctx, hostname)
	}
	return nil
}

// requestECDSA returns true if requesting ECDSA algorithms.
func requestECDSA(hello *tls.ClientHelloInfo) bool {
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
