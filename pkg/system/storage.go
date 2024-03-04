package system

import (
	"path/filepath"
	"time"

	"github.com/seal-io/utils/osx"
)

// DataDir is the path to expose the data.
const DataDir = "/var/run/walrus"

// SubDataDir returns the path to the subdirectory of DataDir.
func SubDataDir(sub string) string {
	if osx.Getenv("_RUNNING_INSIDE_CONTAINER_", "false") == "true" {
		return filepath.Join(DataDir, sub)
	}
	// NB(thxCode): nice for development.
	return osx.SubTempDir(time.Now().Format(time.DateOnly))
}
