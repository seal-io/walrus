package certs

import (
	"context"

	"golang.org/x/crypto/acme/autocert"
	cert "k8s.io/api/certificates/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// ErrCacheMiss is returned when a certificate is not found in cache.
var ErrCacheMiss = autocert.ErrCacheMiss

// Cache is used by Manager to store and retrieve previously obtained certificates
// and other account data as opaque blobs.
//
// Cache implementations should not rely on the key naming pattern. Keys can
// include any printable ASCII characters, except the following: \/:*?"<>|.
type Cache = autocert.Cache

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

type SecretInterface interface {
	Create(context.Context, *core.Secret, metav1.CreateOptions) (*core.Secret, error)
	Update(context.Context, *core.Secret, metav1.UpdateOptions) (*core.Secret, error)
	Delete(context.Context, string, metav1.DeleteOptions) error
	Get(context.Context, string, metav1.GetOptions) (*core.Secret, error)
	List(context.Context, metav1.ListOptions) (*core.SecretList, error)
	Watch(context.Context, metav1.ListOptions) (watch.Interface, error)
}

// CertificateSigningRequestInterface holds the operations on certificate signing requests.
type CertificateSigningRequestInterface interface {
	Create(context.Context, *cert.CertificateSigningRequest, metav1.CreateOptions) (*cert.CertificateSigningRequest, error)
	Update(context.Context, *cert.CertificateSigningRequest, metav1.UpdateOptions) (*cert.CertificateSigningRequest, error)
	Delete(context.Context, string, metav1.DeleteOptions) error
	Get(context.Context, string, metav1.GetOptions) (*cert.CertificateSigningRequest, error)
	List(context.Context, metav1.ListOptions) (*cert.CertificateSigningRequestList, error)
	Watch(context.Context, metav1.ListOptions) (watch.Interface, error)
	UpdateApproval(context.Context, string, *cert.CertificateSigningRequest, metav1.UpdateOptions) (
		*cert.CertificateSigningRequest, error)
}
