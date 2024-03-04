// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) CreateBranch(ctx context.Context, repo string, params *scm.ReferenceInput) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/branches", repo)
	in := &branchCreate{
		Refs:       params.Sha,
		BranchName: params.Name,
	}
	res, err := s.client.do(ctx, "POST", path, in, nil)
	return res, err
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/branches/%s", repo, name)
	out := new(branch)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertBranch(out), res, err
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/commits/%s", repo, ref)
	out := new(commit)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommit(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	tags, res, err := s.ListTags(ctx, repo, scm.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	for _, tag := range tags {
		if tag.Name == name {
			return tag, res, err
		}
	}
	return nil, res, err
}

func (s *gitService) ListBranches(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/branches", repo)
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertBranchList(out), res, err
}

func (s *gitService) ListBranchesV2(ctx context.Context, repo string, opts scm.BranchListOptions) ([]*scm.Reference, *scm.Response, error) {
	// Gitee doesnt provide support listing based on searchTerm
	// Hence calling the ListBranches
	return s.ListBranches(ctx, repo, opts.PageListOptions)
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/commits?%s", repo, encodeCommitListOptions(opts))
	out := []*commit{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommitList(out), res, err
}

func (s *gitService) ListTags(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/tags", repo)
	out := []*releasesTags{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertTagList(out), res, err
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/commits/%s", repo, ref)
	out := new(commit)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertChangeList(out.Files), res, err
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/compare/%s...%s", repo, source, target)
	out := new(compare)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertChangeList(out.Files), res, err
}

type branchCreate struct {
	Refs       string `json:"refs"`
	BranchName string `json:"branch_name"`
}

type branch struct {
	//Links         string `json:"_links"`
	Name          string `json:"name"`
	Commit        tree   `json:"commit"`
	Protected     bool   `json:"protected"`
	ProtectionURL string `json:"protection_url"`
}

type commit struct {
	URL         string `json:"url"`
	Sha         string `json:"sha"`
	HtmlURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Commit      struct {
		Author    committer `json:"author"`
		Committer committer `json:"committer"`
		Message   string    `json:"message"`
		Tree      tree      `json:"tree"`
	} `json:"commit"`
	Author    author `json:"author"`
	Committer author `json:"committer"`
	Parents   []tree `json:"parents"`
	Stats     struct {
		ID        string `json:"id"`
		Additions int64  `json:"additions"`
		Deletions int64  `json:"deletions"`
		Total     int64  `json:"total"`
	} `json:"stats"`
	Files []*file `json:"files"`
}

type committer struct {
	Name  string    `json:"name"`
	Date  time.Time `json:"date"`
	Email string    `json:"email"`
}

type tree struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

type file struct {
	SHA        string `json:"sha"`
	Filename   string `json:"filename"`
	Status     string `json:"status"`
	Additions  int64  `json:"additions"`
	Deletions  int64  `json:"deletions"`
	Changes    int64  `json:"changes"`
	BlobURL    string `json:"blob_url"`
	RawURL     string `json:"raw_url"`
	ContentURL string `json:"content_url"`
	Patch      string `json:"patch"`
}

type compare struct {
	BaseCommit struct {
		URL         string        `json:"url"`
		Sha         string        `json:"sha"`
		HTMLURL     string        `json:"html_url"`
		CommentsURL string        `json:"comments_url"`
		Commit      compareCommit `json:"commit"`
		Author      compareAuthor `json:"author"`
		Committer   committer     `json:"committer"`
		Parents     []tree        `json:"parents"`
	} `json:"base_commit"`
	MergeBaseCommit struct {
		URL         string        `json:"url"`
		Sha         string        `json:"sha"`
		HTMLURL     string        `json:"html_url"`
		CommentsURL string        `json:"comments_url"`
		Commit      compareCommit `json:"commit"`
		Author      compareAuthor `json:"author"`
		Committer   committer     `json:"committer"`
		Parents     []tree        `json:"parents"`
	} `json:"merge_base_commit"`
	Commits []struct {
		URL         string        `json:"url"`
		Sha         string        `json:"sha"`
		HTMLURL     string        `json:"html_url"`
		CommentsURL string        `json:"comments_url"`
		Commit      compareCommit `json:"commit"`
		Author      compareAuthor `json:"author"`
		Committer   committer     `json:"committer"`
		Parents     []tree        `json:"parents"`
	} `json:"commits"`
	Files []*file `json:"files"`
}
type compareAuthor struct {
	Name  string    `json:"name"`
	Date  time.Time `json:"date"`
	Email string    `json:"email"`
}
type compareCommit struct {
	Author    compareAuthor `json:"author"`
	Committer committer     `json:"committer"`
	Message   string        `json:"message"`
	Tree      tree          `json:"tree"`
}

type releasesTags struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Commit  struct {
		Sha  string    `json:"sha"`
		Date time.Time `json:"date"`
	} `json:"commit"`
}

type author struct {
	ID                int    `json:"id"`
	Login             string `json:"login"`
	Name              string `json:"name"`
	AvatarURL         string `json:"avatar_url"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
}

func convertCommitList(from []*commit) []*scm.Commit {
	to := []*scm.Commit{}
	for _, v := range from {
		to = append(to, convertCommit(v))
	}
	return to
}

func convertCommit(from *commit) *scm.Commit {
	return &scm.Commit{
		Message: from.Commit.Message,
		Sha:     from.Sha,
		Link:    from.HtmlURL,
		Author: scm.Signature{
			Name:   from.Commit.Author.Name,
			Email:  from.Commit.Author.Email,
			Date:   from.Commit.Author.Date,
			Login:  from.Author.Login,
			Avatar: from.Author.AvatarURL,
		},
		Committer: scm.Signature{
			Name:   from.Commit.Committer.Name,
			Email:  from.Commit.Committer.Email,
			Date:   from.Commit.Committer.Date,
			Login:  from.Committer.Login,
			Avatar: from.Committer.AvatarURL,
		},
	}
}

func convertBranchList(from []*branch) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from {
		to = append(to, convertBranch(v))
	}
	return to
}

func convertBranch(from *branch) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.Name),
		Path: scm.ExpandRef(from.Name, "refs/heads/"),
		Sha:  from.Commit.Sha,
	}
}

func convertTagList(from []*releasesTags) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from {
		to = append(to, convertTag(v))
	}
	return to
}

func convertTag(from *releasesTags) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.Name),
		Path: scm.ExpandRef(from.Name, "refs/tags/"),
		Sha:  from.Commit.Sha,
	}
}

func convertChangeList(from []*file) []*scm.Change {
	to := []*scm.Change{}
	for _, v := range from {
		to = append(to, convertChange(v))
	}
	return to
}

func convertChange(from *file) *scm.Change {
	return &scm.Change{
		Path:    from.Filename,
		Added:   from.Status == "added",
		Deleted: from.Status == "removed",
		Renamed: from.Status == "modified" && from.Additions == 0 && from.Deletions == 0 && from.Changes == 0,
		BlobID:  from.SHA,
	}
}
