package templates

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/scaleway/scaleway-sdk-go/logger"

	"github.com/seal-io/walrus/pkg/vcs"
)

// gitRepoIconURL retrieves template icon from a git repository and return icon URL.
func gitRepoIconURL(r *git.Repository, vr vcs.Repository) (string, error) {
	// Get icon path.
	p, err := gitRepoIconFilePath(r, vr.SubPath)
	if err != nil {
		logger.Errorf("failed to get icon url: %v", err)
		return "", err
	}

	u, err := vr.FileRawURL(p)
	if err != nil {
		return "", err
	}
	return u, nil
}

// gitRepoIconFileName retrieves template icon from a git repository and return icon path.
func gitRepoIconFilePath(repoLocal *git.Repository, subPath string) (string, error) {
	var (
		err error
		// Valid icon files.
		icons = []string{
			"icon.png",
			"icon.jpg",
			"icon.jpeg",
			"icon.svg",
		}
	)

	w, err := repoLocal.Worktree()
	if err != nil {
		return "", err
	}

	// Get icon URL.
	for _, icon := range icons {
		if subPath != "" {
			icon = filepath.Join(subPath, icon)
		}
		// If icon exists, get icon rawURL.
		if _, err := w.Filesystem.Stat(icon); err == nil {
			return icon, nil
		}
	}

	return "", nil
}
