package vcs

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-version"

	"github.com/seal-io/walrus/pkg/vcs/driver/github"
	"github.com/seal-io/walrus/utils/log"
)

type Repository struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Link is the link of the repository.
	Link string `json:"link"`
	// Reference is the reference of the repository. E.G: main, dev, v0.0.1.
	Reference string `json:"reference"`
	// Driver is the driver of the repository. E.G: github.
	Driver string `json:"driver"`
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

	var driver string

	switch endpoint.Host {
	case "github.com":
		driver = github.Driver
	case "gitlab.com":
		driver = github.Driver
	}

	return &Repository{
		Namespace: namespace,
		Name:      name,
		Link:      rawURL,
		Reference: ref,
		Driver:    driver,
	}, nil
}

// HardResetGitRepo hard resets a git repository to a specific hash.
func HardResetGitRepo(r *git.Repository, ref string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	resetRef, err := GetRepoRef(r, ref)
	if err != nil {
		return err
	}

	err = w.Reset(&git.ResetOptions{
		Commit: resetRef.Hash(),
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
			logger.Warnf("failed to parse tag %s: %v", ref.Name().Short(), verr)
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
func CloneGitRepo(ctx context.Context, link, dir string) (*git.Repository, error) {
	logger := log.WithName("template")

	src, err := GetGitSource(link)
	if err != nil {
		return nil, err
	}

	// Clone git repository.
	err = getter.Get(dir, src, getter.WithContext(ctx))
	if err != nil {
		logger.Errorf("failed to get %s: %v", link, err)

		return nil, err
	}

	return git.PlainOpen(dir)
}

// GetGitSource get git source for template.
// When source's protocol is http or https,
// the prefix git:: will be added for template to use.
func GetGitSource(link string) (string, error) {
	endpoint, err := transport.NewEndpoint(link)
	if err != nil {
		return "", err
	}

	var src string

	switch endpoint.Protocol {
	case "http", "https":
		src = "git::" + endpoint.String()
	default:
		src = link
	}

	return src, nil
}

// GetRepoRef returns a reference from a git repository.
func GetRepoRef(r *git.Repository, name string) (*plumbing.Reference, error) {
	if ref, err := r.Reference(plumbing.NewTagReferenceName(name), true); err == nil {
		return ref, nil
	}

	if ref, err := r.Reference(plumbing.NewBranchReferenceName(name), true); err == nil {
		return ref, nil
	}

	if ref, err := r.Reference(plumbing.NewRemoteReferenceName("origin", name), true); err == nil {
		return ref, nil
	}

	if ref, err := r.Reference(plumbing.NewNoteReferenceName(name), true); err == nil {
		return ref, nil
	}

	return nil, fmt.Errorf("failed to get reference: %s", name)
}

// GetOrgRepos returns full repositories list from the given org.
func GetOrgRepos(ctx context.Context, client *scm.Client, orgName string) ([]*scm.Repository, error) {
	opts := scm.ListOptions{Size: 100}

	var list []*scm.Repository

	for {
		repos, meta, err := client.Organizations.ListRepositories(ctx, orgName, opts)
		if err != nil {
			return nil, err
		}

		for _, src := range repos {
			if src != nil {
				list = append(list, src)
			}
		}

		opts.Page = meta.Page.Next
		opts.URL = meta.Page.NextURL

		if opts.Page == 0 && opts.URL == "" {
			break
		}
	}

	return list, nil
}

// GetOrgFromGitURL parses the organization from the given URL.
func GetOrgFromGitURL(source string) (string, error) {
	u, err := url.Parse(source)
	if err != nil {
		return "", err
	}

	parts := strings.Split(u.Path, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid git url")
	}

	return strings.TrimPrefix(u.Path, "/"), nil
}
