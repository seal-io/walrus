package system

import (
	"path/filepath"
	"time"

	"github.com/seal-io/utils/osx"
)

// DataDir is the path to expose the data.
var _DataDir = osx.Getenv("WALRUS_DATA_DIR", "/var/run/walrus")

// DataDir returns the path to the data directory.
func DataDir() string {
	return _DataDir
}

// SubDataDir returns the path to the subdirectory of DataDir.
func SubDataDir(sub string) string {
	if osx.Getenv("_RUNNING_INSIDE_CONTAINER_", "false") == "true" {
		return filepath.Join(_DataDir, sub)
	}
	// NB(thxCode): nice for development.
	return osx.SubTempDir(filepath.Join(time.Now().Format(time.DateOnly), _DataDir, sub))
}

// LibDir is the path to access the metadata.
var _LibDir = osx.Getenv("WALRUS_LIB_DIR", "/var/lib/walrus")

// LibDir returns the path to the lib directory.
func LibDir() string {
	return _LibDir
}

// SubLibDir returns the path to the subdirectory of LibDir.
func SubLibDir(sub string) string {
	if osx.Getenv("_RUNNING_INSIDE_CONTAINER_", "false") == "true" {
		return filepath.Join(_LibDir, sub)
	}
	// NB(thxCode): nice for development.
	return osx.SubTempDir(filepath.Join(time.Now().Format(time.DateOnly), _LibDir, sub))
}
