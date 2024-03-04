// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/drone/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(ctx context.Context, repo string, number int) (*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d", encode(repo), number)
	out := new(pr)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) FindComment(ctx context.Context, repo string, index, id int) (*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/notes/%d", encode(repo), index, id)
	out := new(issueComment)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertIssueComment(out), res, err
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests?%s", encode(repo), encodePullRequestListOptions(opts))
	out := []*pr{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertPullRequestList(out), res, err
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/changes?%s", encode(repo), number, encodeListOptions(opts))
	out := new(changes)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertChangeList(out.Changes), res, err
}

func (s *pullService) ListComments(ctx context.Context, repo string, index int, opts scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/notes?%s", encode(repo), index, encodeListOptions(opts))
	out := []*issueComment{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertIssueCommentList(out), res, err
}

func (s *pullService) ListCommits(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/commits?%s", encode(repo), number, encodeListOptions(opts))
	out := []*commit{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommitList(out), res, err
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	in := url.Values{}
	in.Set("title", input.Title)
	in.Set("description", input.Body)
	in.Set("source_branch", input.Source)
	in.Set("target_branch", input.Target)
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests?%s", encode(repo), in.Encode())
	out := new(pr)
	res, err := s.client.do(ctx, "POST", path, nil, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) CreateComment(ctx context.Context, repo string, index int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	in := url.Values{}
	in.Set("body", input.Body)
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/notes?%s", encode(repo), index, in.Encode())
	out := new(issueComment)
	res, err := s.client.do(ctx, "POST", path, nil, out)
	return convertIssueComment(out), res, err
}

func (s *pullService) DeleteComment(ctx context.Context, repo string, index, id int) (*scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/notes/%d", encode(repo), index, id)
	res, err := s.client.do(ctx, "DELETE", path, nil, nil)
	return res, err
}

func (s *pullService) Merge(ctx context.Context, repo string, number int) (*scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/merge", encode(repo), number)
	res, err := s.client.do(ctx, "PUT", path, nil, nil)
	return res, err
}

func (s *pullService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d?state_event=closed", encode(repo), number)
	res, err := s.client.do(ctx, "PUT", path, nil, nil)
	return res, err
}

type pr struct {
	Number int    `json:"iid"`
	Sha    string `json:"sha"`
	Title  string `json:"title"`
	Desc   string `json:"description"`
	State  string `json:"state"`
	Link   string `json:"web_url"`
	Author struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Avatar   string `json:"avatar_url"`
	}
	SourceBranch string    `json:"source_branch"`
	TargetBranch string    `json:"target_branch"`
	Created      time.Time `json:"created_at"`
	Updated      time.Time `json:"updated_at"`
	Closed       time.Time
	Labels       []string `json:"labels"`
}

type changes struct {
	Changes []*change
}

type change struct {
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
	Added   bool   `json:"new_file"`
	Renamed bool   `json:"renamed_file"`
	Deleted bool   `json:"deleted_file"`
}

func convertPullRequestList(from []*pr) []*scm.PullRequest {
	to := []*scm.PullRequest{}
	for _, v := range from {
		to = append(to, convertPullRequest(v))
	}
	return to
}

func convertPullRequest(from *pr) *scm.PullRequest {
	var labels []scm.Label
	for _, label := range from.Labels {
		labels = append(labels, scm.Label{
			Name: label,
		})
	}
	return &scm.PullRequest{
		Number: from.Number,
		Title:  from.Title,
		Body:   from.Desc,
		Sha:    from.Sha,
		Ref:    fmt.Sprintf("refs/merge-requests/%d/head", from.Number),
		Source: from.SourceBranch,
		Target: from.TargetBranch,
		Link:   from.Link,
		Closed: from.State != "opened",
		Merged: from.State == "merged",
		Author: scm.User{
			Name:   from.Author.Name,
			Login:  from.Author.Username,
			Avatar: from.Author.Avatar,
		},
		Created: from.Created,
		Updated: from.Updated,
		Labels:  labels,
	}
}

func convertChangeList(from []*change) []*scm.Change {
	to := []*scm.Change{}
	for _, v := range from {
		to = append(to, convertChange(v))
	}
	return to
}

func convertChange(from *change) *scm.Change {
	to := &scm.Change{
		Path:    from.NewPath,
		Added:   from.Added,
		Deleted: from.Deleted,
		Renamed: from.Renamed,
	}
	if to.Path == "" {
		to.Path = from.OldPath
	}
	return to
}
