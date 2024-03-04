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

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(ctx context.Context, repo string, number int) (*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d", repo, number)
	out := new(pr)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) FindComment(ctx context.Context, repo string, _ int, id int) (*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/comments/%d", repo, id)
	out := new(prComment)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequestComment(out), res, err
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls?%s", repo, encodePullRequestListOptions(opts))
	out := []*pr{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertPullRequestList(out), res, err
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/files", repo, number)
	out := []*prFile{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertPrChangeList(out), res, err
}

func (s *pullService) ListComments(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/comments/?%s", repo, number, encodeListOptions(opts))
	out := []*prComment{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertPullRequestComments(out), res, err
}

func (s *pullService) ListCommits(ctx context.Context, repo string, number int, _ scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("/repos/%s/pulls/%d/commits", repo, number)
	out := []*prCommit{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertPrCommitList(out), res, err
}

func (s *pullService) Merge(ctx context.Context, repo string, number int) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/merge", repo, number)
	res, err := s.client.do(ctx, "PUT", path, nil, nil)
	return res, err
}

func (s *pullService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d", repo, number)
	data := map[string]string{"state": "closed"}
	res, err := s.client.do(ctx, "PATCH", path, &data, nil)
	return res, err
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls", repo)
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

func (s *pullService) CreateComment(ctx context.Context, repo string, number int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/comments", repo, number)
	in := &prCommentInput{
		Body: input.Body,
	}
	out := new(prComment)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertPullRequestComment(out), res, err
}

func (s *pullService) DeleteComment(ctx context.Context, repo string, _ int, id int) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/comments/%d", repo, id)
	res, err := s.client.do(ctx, "DELETE", path, nil, nil)
	return res, err
}

type (
	pr struct {
		ID                int        `json:"id"`
		URL               string     `json:"url"`
		HtmlURL           string     `json:"html_url"`
		DiffURL           string     `json:"diff_url"`
		PatchURL          string     `json:"patch_url"`
		IssueURL          string     `json:"issue_url"`
		CommitsURL        string     `json:"commits_url"`
		ReviewCommentsURL string     `json:"review_comments_url"`
		ReviewCommentURL  string     `json:"review_comment_url"`
		CommentsURL       string     `json:"comments_url"`
		Number            int        `json:"number"`
		State             string     `json:"state"`
		Title             string     `json:"title"`
		Body              string     `json:"body"`
		AssigneesNumber   int        `json:"assignees_number"`
		TestersNumber     int        `json:"testers_number"`
		Assignees         []assignee `json:"assignees"`
		Testers           []tester   `json:"testers"`
		Milestone         milestone  `json:"milestone"`
		Labels            []label    `json:"labels"`
		Locked            bool       `json:"locked"`
		CreatedAt         time.Time  `json:"created_at"`
		UpdatedAt         time.Time  `json:"updated_at"`
		ClosedAt          time.Time  `json:"closed_at"`
		MergedAt          time.Time  `json:"merged_at"`
		Mergeable         bool       `json:"mergeable"`
		CanMergeCheck     bool       `json:"can_merge_check"`
		Head              headOrBase `json:"head"`
		Base              headOrBase `json:"base"`
		User              user       `json:"user"`
	}
	assignee struct {
		ID                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarURL         string `json:"avatar_url"`
		URL               string `json:"url"`
		HtmlURL           string `json:"html_url"`
		Remark            string `json:"remark"`
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
	tester struct {
		ID                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarURL         string `json:"avatar_url"`
		URL               string `json:"url"`
		HtmlURL           string `json:"html_url"`
		Remark            string `json:"remark"`
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
	milestone struct {
		URL          string      `json:"url"`
		HtmlURL      string      `json:"html_url"`
		ID           int         `json:"id"`
		Number       int         `json:"number"`
		RepositoryID interface{} `json:"repository_id"`
		State        string      `json:"state"`
		Title        string      `json:"title"`
		Description  string      `json:"description"`
		UpdatedAt    time.Time   `json:"updated_at"`
		CreatedAt    time.Time   `json:"created_at"`
		OpenIssues   int         `json:"open_issues"`
		ClosedIssues int         `json:"closed_issues"`
		DueOn        string      `json:"due_on"`
	}
	label struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		Color        string    `json:"color"`
		RepositoryID int       `json:"repository_id"`
		URL          string    `json:"url"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
	prRepo struct {
		ID        int    `json:"id"`
		FullName  string `json:"full_name"`
		HumanName string `json:"human_name"`
		URL       string `json:"url"`
		//Namespace   namespace `json:"namespace"`
		Path        string   `json:"path"`
		Name        string   `json:"name"`
		Owner       user     `json:"owner"`
		Assigner    assignee `json:"assigner"`
		Description string   `json:"description"`
		Private     bool     `json:"private"`
		Public      bool     `json:"public"`
		Internal    bool     `json:"internal"`
		Fork        bool     `json:"fork"`
		HtmlURL     string   `json:"html_url"`
		SshURL      string   `json:"ssh_url"`
	}
	headOrBase struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Sha   string `json:"sha"`
		User  user   `json:"user"`
		Repo  prRepo `json:"repo"`
	}

	prComment struct {
		URL            string    `json:"url"`
		ID             int       `json:"id"`
		User           user      `json:"user"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Body           string    `json:"body"`
		HtmlURL        string    `json:"html_url"`
		PullRequestURL string    `json:"pull_request_url"`
	}

	prFile struct {
		Sha       string `json:"sha"`
		Filename  string `json:"filename"`
		Status    string `json:"status"`
		Additions string `json:"additions"`
		Deletions string `json:"deletions"`
		BlobURL   string `json:"blob_url"`
		RawURL    string `json:"raw_url"`
		Patch     struct {
			Diff        string `json:"diff"`
			NewPath     string `json:"new_path"`
			OldPath     string `json:"old_path"`
			AMode       string `json:"a_mode"`
			BMode       string `json:"b_mode"`
			NewFile     bool   `json:"new_file"`
			RenamedFile bool   `json:"renamed_file"`
			DeletedFile bool   `json:"deleted_file"`
			TooLarge    bool   `json:"too_large"`
		} `json:"patch"`
	}

	prCommit struct {
		URL         string `json:"url"`
		Sha         string `json:"sha"`
		HtmlURL     string `json:"html_url"`
		CommentsURL string `json:"comments_url"`
		Commit      struct {
			URL    string `json:"url"`
			Author struct {
				Name  string    `json:"name"`
				Date  time.Time `json:"date"`
				Email string    `json:"email"`
			} `json:"author"`
			Committer struct {
				Name  string    `json:"name"`
				Date  time.Time `json:"date"`
				Email string    `json:"email"`
			} `json:"committer"`
			Message      string `json:"message"`
			CommentCount int    `json:"comment_count"`
		} `json:"commit"`
		Author    author `json:"author"`
		Committer author `json:"committer"`
		Parents   struct {
			URL string `json:"url"`
			Sha string `json:"sha"`
		} `json:"parents"`
	}
)

type prInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}
type prCommentInput struct {
	Body string `json:"body"`
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
			Name:  label.Name,
			Color: label.Color,
		})
	}
	merged := from.State == "merged"
	closed := from.State == "closed"
	if merged {
		closed = true
	}
	return &scm.PullRequest{
		Number: from.Number,
		Title:  from.Title,
		Body:   from.Body,
		Sha:    from.Head.Sha,
		Ref:    fmt.Sprintf("refs/pull/%d/head", from.Number),
		Source: from.Head.Ref,
		Target: from.Base.Ref,
		Fork:   from.Head.Repo.FullName,
		Link:   from.HtmlURL,
		Diff:   from.DiffURL,
		Closed: closed,
		Merged: merged,
		Head: scm.Reference{
			Name: from.Head.Ref,
			Path: scm.ExpandRef(from.Head.Ref, "refs/heads"),
			Sha:  from.Head.Sha,
		},
		Base: scm.Reference{
			Name: from.Base.Ref,
			Path: scm.ExpandRef(from.Base.Ref, "refs/heads"),
			Sha:  from.Base.Sha,
		},
		Author: scm.User{
			Login:  from.User.Login,
			Avatar: from.User.AvatarURL,
		},
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
		Labels:  labels,
	}
}

func convertPullRequestComments(from []*prComment) []*scm.Comment {
	to := []*scm.Comment{}
	for _, v := range from {
		to = append(to, convertPullRequestComment(v))
	}
	return to
}
func convertPullRequestComment(from *prComment) *scm.Comment {
	return &scm.Comment{
		ID:   from.ID,
		Body: from.Body,
		Author: scm.User{
			Login:  from.User.Login,
			Name:   from.User.Name,
			Avatar: from.User.AvatarURL,
		},
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
	}
}

func convertPrChangeList(from []*prFile) []*scm.Change {
	to := []*scm.Change{}
	for _, v := range from {
		to = append(to, convertPrChange(v))
	}
	return to
}
func convertPrChange(from *prFile) *scm.Change {
	return &scm.Change{
		Path:    from.Filename,
		Added:   from.Status == "added",
		Deleted: from.Status == "deleted",
		Renamed: from.Status == "renamed",
		BlobID:  from.Sha,
	}
}

func convertPrCommitList(from []*prCommit) []*scm.Commit {
	to := []*scm.Commit{}
	for _, v := range from {
		to = append(to, convertPrCommit(v))
	}
	return to
}
func convertPrCommit(from *prCommit) *scm.Commit {
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
