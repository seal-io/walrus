// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(context.Context, string, int) (*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) FindComment(context.Context, string, int, int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) List(context.Context, string, scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListComments(context.Context, string, int, scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListChanges(context.Context, string, int, scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListCommits(context.Context, string, int, scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) Create(context.Context, string, *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) CreateComment(context.Context, string, int, *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Merge(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Close(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

//
// native data structures
//

type pullRequest struct {
	ID         int        `json:"id"`
	Number     int        `json:"number"`
	User       user       `json:"user"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	State      string     `json:"state"`
	HeadBranch string     `json:"head_branch"`
	HeadRepo   repository `json:"head_repo"`
	BaseBranch string     `json:"base_branch"`
	BaseRepo   repository `json:"base_repo"`
	HTMLURL    string     `json:"html_url"`
	Mergeable  bool       `json:"mergeable"`
	Merged     bool       `json:"merged"`
}

//
// native data structure conversion
//

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
