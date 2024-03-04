// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) CreateBranch(ctx context.Context, repo string, params *scm.ReferenceInput) (*scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/branches", harnessURI)
	in := &branchInput{
		Name:   params.Name,
		Target: params.Sha,
	}
	return s.client.do(ctx, "POST", path, in, nil)
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/branches/%s", harnessURI, name)
	out := new(branch)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertBranch(out), res, err
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/commits/%s", harnessURI, ref)
	out := new(commitInfo)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommitInfo(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/branches?%s", harnessURI, encodeListOptions(opts))
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertBranchList(out), res, err
}

func (s *gitService) ListBranchesV2(ctx context.Context, repo string, opts scm.BranchListOptions) ([]*scm.Reference, *scm.Response, error) {
	// Harness doesnt provide support listing based on searchTerm
	// Hence calling the ListBranches
	return s.ListBranches(ctx, repo, opts.PageListOptions)
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/commits?%s", harnessURI, encodeCommitListOptions(opts))
	out := new(commits)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommitList(out), res, err
}

func (s *gitService) ListTags(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/commits/%s/diff?%s", harnessURI, ref, encodeListOptions(opts))
	out := []*fileDiff{}
	res, err := s.client.do(ctx, "POST", path, nil, &out)
	return convertFileDiffs(out), res, err
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/diff/%s...%s", harnessURI, source, target)
	out := []*fileDiff{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertChangeList(out), res, err
}

// native data structures
type (
	commits struct {
		Commits []commitInfo `json:"commits"`
	}

	commitInfo struct {
		Author struct {
			Identity struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"committer"`
		Message string `json:"message"`
		Sha     string `json:"sha"`
		Title   string `json:"title"`
	}
	branchInput struct {
		Name   string `json:"name"`
		Target string `json:"target"`
	}
	branch struct {
		Commit struct {
			Author struct {
				Identity struct {
					Email string `json:"email"`
					Name  string `json:"name"`
				} `json:"identity"`
				When time.Time `json:"when"`
			} `json:"author"`
			Committer struct {
				Identity struct {
					Email string `json:"email"`
					Name  string `json:"name"`
				} `json:"identity"`
				When time.Time `json:"when"`
			} `json:"committer"`
			Message string `json:"message"`
			Sha     string `json:"sha"`
			Title   string `json:"title"`
		} `json:"commit"`
		Name string `json:"name"`
		Sha  string `json:"sha"`
	}
	fileDiff struct {
		SHA         string `json:"sha"`
		OldSHA      string `json:"old_sha,omitempty"`
		Path        string `json:"path"`
		OldPath     string `json:"old_path,omitempty"`
		Status      string `json:"status"`
		Additions   int64  `json:"additions"`
		Deletions   int64  `json:"deletions"`
		Changes     int64  `json:"changes"`
		ContentURL  string `json:"content_url"`
		Patch       []byte `json:"patch,omitempty"`
		IsBinary    bool   `json:"is_binary"`
		IsSubmodule bool   `json:"is_submodule"`
	}
)

//
// native data structure conversion
//

func convertBranchList(src []*branch) []*scm.Reference {
	dst := []*scm.Reference{}
	for _, v := range src {
		dst = append(dst, convertBranch(v))
	}
	return dst
}

func convertBranch(src *branch) *scm.Reference {
	return &scm.Reference{
		Name: src.Name,
		Path: scm.ExpandRef(src.Name, "refs/heads/"),
		Sha:  src.Sha,
	}
}

func convertCommitList(src *commits) []*scm.Commit {
	var dst []*scm.Commit
	for _, v := range src.Commits {
		dst = append(dst, convertCommitInfo(&v))
	}
	return dst
}

func convertChangeList(src []*fileDiff) []*scm.Change {
	dst := []*scm.Change{}
	for _, v := range src {
		dst = append(dst, convertChange(v))
	}
	return dst
}

func convertCommitInfo(src *commitInfo) *scm.Commit {
	return &scm.Commit{
		Sha:     src.Sha,
		Message: src.Message,
		Author: scm.Signature{
			Name:  src.Author.Identity.Name,
			Email: src.Author.Identity.Email,
			Date:  src.Author.When,
		},
		Committer: scm.Signature{
			Name:  src.Committer.Identity.Name,
			Email: src.Committer.Identity.Email,
			Date:  src.Committer.When,
		},
	}
}

func convertChange(src *fileDiff) *scm.Change {
	return &scm.Change{
		Path:         src.Path,
		PrevFilePath: src.OldPath,
		Added:        strings.EqualFold(src.Status, "ADDED"),
		Renamed:      strings.EqualFold(src.Status, "RENAMED"),
		Deleted:      strings.EqualFold(src.Status, "DELETED"),
	}
}
