package cache

import (
	"golang.org/x/crypto/acme/autocert"

	"github.com/seal-io/utils/certs"
)

// NewDirCache returns a new DirCache instance with the given directory.
func NewDirCache(dir string) certs.Cache {
	return autocert.DirCache(dir)
}
