// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(ctx context.Context, repo string, index int) (*scm.PullRequest, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d", harnessURI, index)
	out := new(pr)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err

}

func (s *pullService) FindComment(context.Context, string, int, int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/pullreq?%s", harnessURI, encodePullRequestListOptions(opts))
	out := []*pr{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertPullRequestList(out), res, err
}

func (s *pullService) ListComments(context.Context, string, int, scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListCommits(ctx context.Context, repo string, index int, opts scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d/commits?%s", harnessURI, index, encodeListOptions(opts))
	out := []*commit{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommits(out), res, err
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d/diff", harnessURI, number)
	out := []*fileDiff{}
	res, err := s.client.do(ctx, "POST", path, nil, &out)
	return convertFileDiffs(out), res, err
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/pullreq", harnessURI)
	in := &prInput{
		Title:        input.Title,
		Description:  input.Body,
		SourceBranch: input.Source,
		TargetBranch: input.Target,
	}
	out := new(pr)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) CreateComment(ctx context.Context, repo string, prNumber int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	harnessQueryParams := fmt.Sprintf("?accountIdentifier=%s&orgIdentifier=%s&projectIdentifier=%s", s.client.account, s.client.organization, s.client.project)
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d/comments%s", repo, prNumber, harnessQueryParams)
	in := &prComment{
		Text: input.Body,
	}
	out := new(prCommentResponse)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertComment(out), res, err
}

func (s *pullService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Merge(ctx context.Context, repo string, index int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Close(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

// native data structures
type (
	pr struct {
		Author struct {
			Created     int    `json:"created"`
			DisplayName string `json:"display_name"`
			Email       string `json:"email"`
			ID          int    `json:"id"`
			Type        string `json:"type"`
			UID         string `json:"uid"`
			Updated     int    `json:"updated"`
		} `json:"author"`
		Created       int    `json:"created"`
		Description   string `json:"description"`
		Edited        int    `json:"edited"`
		IsDraft       bool   `json:"is_draft"`
		MergeBaseSha  string `json:"merge_base_sha"`
		MergeHeadSha  string `json:"merge_head_sha"`
		MergeStrategy string `json:"merge_strategy"`
		Merged        int    `json:"merged"`
		Merger        struct {
			Created     int    `json:"created"`
			DisplayName string `json:"display_name"`
			Email       string `json:"email"`
			ID          int    `json:"id"`
			Type        string `json:"type"`
			UID         string `json:"uid"`
			Updated     int    `json:"updated"`
		} `json:"merger"`
		Number       int    `json:"number"`
		SourceBranch string `json:"source_branch"`
		SourceRepoID int    `json:"source_repo_id"`
		State        string `json:"state"`
		Stats        struct {
			Commits       int `json:"commits"`
			Conversations int `json:"conversations"`
			FilesChanged  int `json:"files_changed"`
		} `json:"stats"`
		TargetBranch string `json:"target_branch"`
		TargetRepoID int    `json:"target_repo_id"`
		Title        string `json:"title"`
	}

	reference struct {
		Repo repository `json:"repo"`
		Name string     `json:"ref"`
		Sha  string     `json:"sha"`
	}

	prInput struct {
		Description   string `json:"description"`
		IsDraft       bool   `json:"is_draft"`
		SourceBranch  string `json:"source_branch"`
		SourceRepoRef string `json:"source_repo_ref"`
		TargetBranch  string `json:"target_branch"`
		Title         string `json:"title"`
	}

	commit struct {
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
	prComment struct {
		LineEnd         int    `json:"line_end"`
		LineEndNew      bool   `json:"line_end_new"`
		LineStart       int    `json:"line_start"`
		LineStartNew    bool   `json:"line_start_new"`
		ParentID        int    `json:"parent_id"`
		Path            string `json:"path"`
		SourceCommitSha string `json:"source_commit_sha"`
		TargetCommitSha string `json:"target_commit_sha"`
		Text            string `json:"text"`
	}
	prCommentResponse struct {
		Id        int         `json:"id"`
		Created   int64       `json:"created"`
		Updated   int64       `json:"updated"`
		Edited    int64       `json:"edited"`
		ParentId  interface{} `json:"parent_id"`
		RepoId    int         `json:"repo_id"`
		PullreqId int         `json:"pullreq_id"`
		Order     int         `json:"order"`
		SubOrder  int         `json:"sub_order"`
		Type      string      `json:"type"`
		Kind      string      `json:"kind"`
		Text      string      `json:"text"`
		Payload   struct{}    `json:"payload"`
		Metadata  interface{} `json:"metadata"`
		Author    struct {
			Id          int    `json:"id"`
			Uid         string `json:"uid"`
			DisplayName string `json:"display_name"`
			Email       string `json:"email"`
			Type        string `json:"type"`
			Created     int64  `json:"created"`
			Updated     int64  `json:"updated"`
		} `json:"author"`
	}
)

// native data structure conversion
func convertPullRequests(src []*pr) []*scm.PullRequest {
	dst := []*scm.PullRequest{}
	for _, v := range src {
		dst = append(dst, convertPullRequest(v))
	}
	return dst
}

func convertPullRequest(src *pr) *scm.PullRequest {
	return &scm.PullRequest{
		Number: src.Number,
		Title:  src.Title,
		Body:   src.Description,
		Source: src.SourceBranch,
		Target: src.TargetBranch,
		Merged: src.Merged != 0,
		Author: scm.User{
			Login: src.Author.Email,
			Name:  src.Author.DisplayName,
			ID:    src.Author.UID,
			Email: src.Author.Email,
		},
		Fork:   "fork",
		Ref:    fmt.Sprintf("refs/pullreq/%d/head", src.Number),
		Closed: src.State == "closed",
	}
}

func convertCommits(src []*commit) []*scm.Commit {
	dst := []*scm.Commit{}
	for _, v := range src {
		dst = append(dst, convertCommit(v))
	}
	return dst
}

func convertCommit(src *commit) *scm.Commit {
	return &scm.Commit{
		Message: src.Message,
		Sha:     src.Sha,
		Author: scm.Signature{
			Name:  src.Author.Identity.Name,
			Email: src.Author.Identity.Email,
		},
		Committer: scm.Signature{
			Name:  src.Committer.Identity.Name,
			Email: src.Committer.Identity.Email,
		},
	}
}

func convertFileDiffs(diff []*fileDiff) []*scm.Change {
	var dst []*scm.Change
	for _, v := range diff {
		dst = append(dst, convertFileDiff(v))
	}
	return dst
}

func convertFileDiff(diff *fileDiff) *scm.Change {
	return &scm.Change{
		Path:         diff.Path,
		Added:        strings.EqualFold(diff.Status, "ADDED"),
		Renamed:      strings.EqualFold(diff.Status, "RENAMED"),
		Deleted:      strings.EqualFold(diff.Status, "DELETED"),
		Sha:          diff.SHA,
		BlobID:       "",
		PrevFilePath: diff.OldPath,
	}
}

func convertPullRequestList(from []*pr) []*scm.PullRequest {
	to := []*scm.PullRequest{}
	for _, v := range from {
		to = append(to, convertPullRequest(v))
	}
	return to
}

func convertComment(comment *prCommentResponse) *scm.Comment {
	return &scm.Comment{
		ID:   comment.Id,
		Body: comment.Text,
		Author: scm.User{
			Login:   comment.Author.Uid,
			Name:    comment.Author.DisplayName,
			ID:      strconv.Itoa(comment.Author.Id),
			Email:   comment.Author.Email,
			Created: time.UnixMilli(comment.Author.Created),
			Updated: time.UnixMilli(comment.Author.Updated),
		},
		Created: time.UnixMilli(comment.Created),
		Updated: time.UnixMilli(comment.Updated),
	}
}
