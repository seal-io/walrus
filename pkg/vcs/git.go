package vcs

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/seal-io/utils/osx"
)

type GitCloneOptions struct {
	URL             string
	Auth            transport.AuthMethod
	Proxy           transport.ProxyOptions
	InsecureSkipTLS bool
}

var regGitHashRef = regexp.MustCompile(`[0-9a-fA-F]{7,40}`)

// GitClone executes Git clone with the options and returns the cloned filesystem.
//
// If the given directory is empty, a temporary directory will be created.
func GitClone(ctx context.Context, dir string, o GitCloneOptions) (*ClonedFilesystem, error) {
	opts := git.CloneOptions{
		URL:               o.URL,
		Auth:              o.Auth,
		ProxyOptions:      o.Proxy,
		InsecureSkipTLS:   o.InsecureSkipTLS,
		Mirror:            false,
		ReferenceName:     plumbing.HEAD,
		RemoteName:        git.DefaultRemoteName,
		Tags:              git.AllTags,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	var subpath, ref string
	{
		if opts.URL == "" {
			return nil, git.ErrEmptyUrls
		}

		u, err := url.Parse(opts.URL)
		if err != nil {
			return nil, fmt.Errorf("parse given URL: %w", err)
		}
		u.Path, subpath, _ = strings.Cut(u.Path, "//")
		if q := u.Query(); q != nil {
			ref = q.Get("ref")
			q.Del("ref")
			u.RawQuery = q.Encode()
		}

		opts.URL = u.String()
	}

	optsSlice := []git.CloneOptions{opts}
	// If the reference is given, try to clone with the reference name.
	if ref != "" && !regGitHashRef.MatchString(ref) {
		opts.SingleBranch = true
		opts.Depth = 1
		optsSlice = make([]git.CloneOptions, 0, 2)
		o1 := opts
		o1.ReferenceName = plumbing.NewBranchReferenceName(ref)
		optsSlice = append(optsSlice, o1)
		o2 := opts
		o2.ReferenceName = plumbing.NewTagReferenceName(ref)
		optsSlice = append(optsSlice, o2)
	}

	if dir == "" {
		dir = osx.TempDir("git-clone-*")
	}

	var (
		repo *git.Repository
		err  error
	)
	for i := range optsSlice {
		repo, err = git.PlainCloneContext(ctx, dir, false, &optsSlice[i])
		if err != nil && !errors.As(err, &git.NoMatchingRefSpecError{}) {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("git clone : %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("prepare git worktree: %w", err)
	}

	if ref != "" && regGitHashRef.MatchString(ref) {
		var h plumbing.Hash
		if len(ref) == 40 {
			h = plumbing.NewHash(ref)
		} else {
			// Try slow path.
			hp, err := hex.DecodeString(ref[:len(ref)&^1])
			if err != nil {
				return nil, fmt.Errorf("decode given ref: %w", err)
			}
			hs := retrieveHashesWithPrefix(repo.Storer, hp)
			switch {
			case len(hs) == 0:
				return nil, fmt.Errorf("no matching hash for the given ref: %s", ref)
			case len(hs) > 1:
				return nil, fmt.Errorf("ambiguous hash for the given ref: %s", ref)
			}
			h = hs[0]
		}
		err = wt.Checkout(&git.CheckoutOptions{
			Hash: h,
		})
		if err != nil {
			return nil, fmt.Errorf("git checkout: %w", err)
		}
	}

	fs := wt.Filesystem
	if subpath != "" {
		// If the subpath is given, chroot the filesystem.
		fs, err = fs.Chroot(subpath)
		if err != nil {
			return nil, fmt.Errorf("chroot git worktree's subpath: %w", err)
		}
	}

	return &ClonedFilesystem{Filesystem: fs}, nil
}

// retrieveHashesWithPrefix is borrowed from the method expandPartialHash of github.com/go-git/go-git/v5/repository.go,
// which is used to expand the partial hash to the full hash.
func retrieveHashesWithPrefix(s storer.EncodedObjectStorer, hp []byte) []plumbing.Hash {
	// Fast path.
	type fastIter interface {
		HashesWithPrefix(prefix []byte) ([]plumbing.Hash, error)
	}
	if fi, ok := s.(fastIter); ok {
		h, err := fi.HashesWithPrefix(hp)
		if err != nil {
			return nil
		}
		return h
	}

	// Slow path.
	var hashes []plumbing.Hash
	si, err := s.IterEncodedObjects(plumbing.AnyObject)
	if err != nil {
		return nil
	}
	_ = si.ForEach(func(obj plumbing.EncodedObject) error {
		h := obj.Hash()
		if bytes.HasPrefix(h[:], hp) {
			hashes = append(hashes, h)
		}
		return nil
	})
	return hashes
}
