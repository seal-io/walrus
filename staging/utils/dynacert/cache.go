package dynacert

import "golang.org/x/crypto/acme/autocert"

// ErrCacheMiss is returned when a certificate is not found in cache.
var ErrCacheMiss = autocert.ErrCacheMiss

// Cache is used by Manager to store and retrieve previously obtained certificates
// and other account data as opaque blobs.
//
// Cache implementations should not rely on the key naming pattern. Keys can
// include any printable ASCII characters, except the following: \/:*?"<>|.
type Cache = autocert.Cache

// DirCache implements Cache using a directory on the local filesystem.
// If the directory does not exist, it will be created with 0700 permissions.
type DirCache = autocert.DirCache
