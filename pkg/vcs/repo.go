package vcs

import (
	"errors"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-version"

	"github.com/seal-io/seal/utils/log"
)

type Repository struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Reference   string `json:"reference"`
}

// ParseURLToGit parses a raw URL to a git repository.
// Since the return repository only contains the namespace and name,
// It only used for create template from catalog.
func ParseURLToRepo(rawURL string) (*Repository, error) {
	// Trim git:: prefix.
	rawURL = strings.TrimPrefix(rawURL, "git::")
	ref := ""
	name := ""
	namespace := ""

	endpoint, err := transport.NewEndpoint(rawURL)
	if err != nil {
		return nil, err
	}

	path := endpoint.Path

	// Get ref from path.
	if strings.Contains(path, "?ref=") {
		parts := strings.Split(endpoint.Path, "?ref=")
		ref = parts[1]
		path = strings.TrimSuffix(path, "?ref="+ref)
		rawURL = strings.TrimSuffix(rawURL, "?ref="+ref)
	}

	// Trim .git suffix.
	path = strings.TrimSuffix(path, ".git")

	switch endpoint.Protocol {
	case "https", "http":
		parts := strings.Split(path, "/")
		if len(parts) < 3 {
			return nil, errors.New("invalid repository path")
		}
		namespace = parts[1]
		name = parts[2]

	case "ssh":
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			return nil, errors.New("invalid repository path")
		}
		namespace = parts[0]
		name = parts[1]
	case "file":
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			return nil, errors.New("invalid repository path")
		}
		name = parts[len(parts)-1]
		namespace = strings.Join(parts[:len(parts)-1], "/")
	}

	return &Repository{
		Namespace: namespace,
		Name:      name,
		Link:      rawURL,
		Reference: ref,
	}, nil
}

// HardResetGitRepo hard resets a git repository to a specific hash.
func HardResetGitRepo(r *git.Repository, ref string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(ref),
		Mode:   git.HardReset,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetGitRepoVersions returns all versions of a git repository in descending order.
func GetGitRepoVersions(r *git.Repository) ([]*version.Version, error) {
	logger := log.WithName("vcs")

	tagRefs, err := r.Tags()
	if err != nil {
		return nil, err
	}

	var versions []*version.Version

	err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
		v, verr := version.NewVersion(ref.Name().Short())
		if verr != nil {
			logger.Warnf("failed to parse tag %s: %v", ref.Name().Short(), err)
		}

		if v != nil {
			versions = append(versions, v)
		}

		return nil
	})

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].LessThan(versions[j])
	})

	return versions, err
}

// CloneGitRepo clones a git repository to a specific directory.
func CloneGitRepo(link, dir string) (*git.Repository, error) {
	logger := log.WithName("template")

	endpoint, err := transport.NewEndpoint(link)
	if err != nil {
		return nil, err
	}

	var src string

	switch endpoint.Protocol {
	case "http", "https":
		src = endpoint.Host + endpoint.Path
	default:
		src = link
	}

	// Clone git repository.
	err = getter.Get(dir, src)
	if err != nil {
		logger.Errorf("failed to get %s: %v", link, err)

		return nil, err
	}

	return git.PlainOpen(dir)
}
