// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) CreateBranch(ctx context.Context, repo string, params *scm.ReferenceInput) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/branches/%s", repo, name)
	out := new(branch)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertBranch(out), res, err
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	ref = scm.TrimRef(ref)
	path := fmt.Sprintf("api/v1/repos/%s/git/commits/%s", repo, url.PathEscape(ref))
	out := new(commitInfo)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommitInfo(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	name = scm.TrimRef(name)
	path := fmt.Sprintf("api/v1/repos/%s/git/refs/tags/%s", repo, url.PathEscape(name))
	out := []*tag{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if err != nil {
		return nil, res, err
	}
	for _, tag := range convertTagList(out) {
		if tag.Name == name {
			return tag, res, nil
		}
	}
	return nil, res, scm.ErrNotFound
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/branches?%s", repo, encodeListOptions(opts))
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertBranchList(out), res, err
}

func (s *gitService) ListBranchesV2(ctx context.Context, repo string, opts scm.BranchListOptions) ([]*scm.Reference, *scm.Response, error) {
	// Gitea doesnt provide support listing based on searchTerm
	// Hence calling the ListBranches
	return s.ListBranches(ctx, repo, opts.PageListOptions)
}

func (s *gitService) ListCommits(ctx context.Context, repo string, _ scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/commits", repo)
	out := []*commitInfo{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommitList(out), res, err
}

func (s *gitService) ListTags(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/git/refs/tags", repo)
	out := []*tag{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertTagList(out), res, err
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

//
// native data structures
//

type (
	// gitea branch object.
	branch struct {
		Name   string `json:"name"`
		Commit commit `json:"commit"`
	}

	// gitea commit object.
	commit struct {
		ID        string    `json:"id"`
		Sha       string    `json:"sha"`
		Message   string    `json:"message"`
		URL       string    `json:"url"`
		Author    signature `json:"author"`
		Committer signature `json:"committer"`
		Timestamp time.Time `json:"timestamp"`
	}

	// gitea commit info object.
	commitInfo struct {
		Sha       string `json:"sha"`
		Commit    commit `json:"commit"`
		Author    user   `json:"author"`
		Committer user   `json:"committer"`
	}

	// gitea signature object.
	signature struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	// gitea tag object
	tag struct {
		Ref    string `json:"ref"`
		URL    string `json:"url"`
		Object struct {
			Type string `json:"type"`
			Sha  string `json:"sha"`
			URL  string `json:"url"`
		} `json:"object"`
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
		Name: scm.TrimRef(src.Name),
		Path: scm.ExpandRef(src.Name, "refs/heads/"),
		Sha:  src.Commit.ID,
	}
}

func convertCommitList(src []*commitInfo) []*scm.Commit {
	dst := []*scm.Commit{}
	for _, v := range src {
		dst = append(dst, convertCommitInfo(v))
	}
	return dst
}

func convertCommitInfo(src *commitInfo) *scm.Commit {
	return &scm.Commit{
		Sha:       src.Sha,
		Link:      src.Commit.URL,
		Message:   src.Commit.Message,
		Author:    convertUserSignature(src.Author),
		Committer: convertUserSignature(src.Committer),
	}
}

func convertSignature(src signature) scm.Signature {
	return scm.Signature{
		Login: src.Username,
		Email: src.Email,
		Name:  src.Name,
	}
}

func convertUserSignature(src user) scm.Signature {
	return scm.Signature{
		Login:  userLogin(&src),
		Email:  src.Email,
		Name:   src.Fullname,
		Avatar: src.Avatar,
	}
}

func convertTagList(src []*tag) []*scm.Reference {
	var dst []*scm.Reference
	for _, v := range src {
		dst = append(dst, convertTag(v))
	}
	return dst
}

func convertTag(src *tag) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(src.Ref),
		Path: src.Ref,
		Sha:  src.Object.Sha,
	}
}
