// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(ctx context.Context, repo string, index int) (*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/pulls/%d", repo, index)
	out := new(pr)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) FindComment(context.Context, string, int, int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/pulls?%s", repo, encodePullRequestListOptions(opts))
	out := []*pr{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertPullRequests(out), res, err
}

func (s *pullService) ListComments(context.Context, string, int, scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListCommits(context.Context, string, int, scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListChanges(context.Context, string, int, scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/pulls", repo)
	in := &prInput{
		Title: input.Title,
		Body:  input.Body,
		Head:  input.Source,
		Base:  input.Target,
	}
	out := new(pr)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) CreateComment(context.Context, string, int, *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Merge(ctx context.Context, repo string, index int) (*scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/pulls/%d/merge", repo, index)
	res, err := s.client.do(ctx, "POST", path, nil, nil)
	return res, err
}

func (s *pullService) Close(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

//
// native data structures
//

type pr struct {
	ID         int        `json:"id"`
	Number     int        `json:"number"`
	User       user       `json:"user"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	State      string     `json:"state"`
	HeadBranch string     `json:"head_branch"`
	HeadRepo   repository `json:"head_repo"`
	Head       reference  `json:"head"`
	BaseBranch string     `json:"base_branch"`
	BaseRepo   repository `json:"base_repo"`
	Base       reference  `json:"base"`
	HTMLURL    string     `json:"html_url"`
	DiffURL    string     `json:"diff_url"`
	Mergeable  bool       `json:"mergeable"`
	Merged     bool       `json:"merged"`
	Created    time.Time  `json:"created_at"`
	Updated    time.Time  `json:"updated_at"`
	Labels     []struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"labels"`
}

type reference struct {
	Repo repository `json:"repo"`
	Name string     `json:"ref"`
	Sha  string     `json:"sha"`
}

type prInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

//
// native data structure conversion
//

func convertPullRequests(src []*pr) []*scm.PullRequest {
	dst := []*scm.PullRequest{}
	for _, v := range src {
		dst = append(dst, convertPullRequest(v))
	}
	return dst
}

func convertPullRequest(src *pr) *scm.PullRequest {
	var labels []scm.Label
	for _, label := range src.Labels {
		labels = append(labels, scm.Label{
			Name:  label.Name,
			Color: label.Color,
		})
	}
	return &scm.PullRequest{
		Number:  src.Number,
		Title:   src.Title,
		Body:    src.Body,
		Sha:     src.Head.Sha,
		Source:  src.Head.Name,
		Target:  src.Base.Name,
		Link:    src.HTMLURL,
		Diff:    src.DiffURL,
		Fork:    src.Base.Repo.FullName,
		Ref:     fmt.Sprintf("refs/pull/%d/head", src.Number),
		Closed:  src.State == "closed",
		Author:  *convertUser(&src.User),
		Merged:  src.Merged,
		Created: src.Created,
		Updated: src.Updated,
		Labels:  labels,
	}
}

func convertPullRequestFromIssue(src *issue) *scm.PullRequest {
	return &scm.PullRequest{
		Number:  src.Number,
		Title:   src.Title,
		Body:    src.Body,
		Closed:  src.State == "closed",
		Author:  *convertUser(&src.User),
		Merged:  src.PullRequest.Merged,
		Created: src.Created,
		Updated: src.Updated,
	}
}
