// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
)

type RepositoryService struct {
	client *wrapper
}

func (s *RepositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out), res, err
}

func (s *RepositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks/%s", repo, id)
	out := new(hook)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertHook(out), res, err
}

func (s *RepositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	repos, res, err := s.Find(ctx, repo)
	if err == nil {
		return repos.Perm, res, err
	}
	return nil, res, err
}

func (s *RepositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("user/repos?%s", encodeListOptions(opts))
	out := []*repository{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertRepositoryList(out), res, err
}
func (s *RepositoryService) ListV2(ctx context.Context, opts scm.RepoListOptions) ([]*scm.Repository, *scm.Response, error) {
	// gitee does not support search filters, hence calling List api without search filtering
	return s.List(ctx, opts.ListOptions)
}

func (s *RepositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks?%s", repo, encodeListOptions(opts))
	out := []*hook{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertHookList(out), res, err
}

func (s *RepositoryService) ListStatus(context.Context, string, string, scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *RepositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks", repo)
	in := new(hook)
	// 1: signature
	in.EncryptionType = 1
	in.Password = input.Secret
	in.URL = input.Target
	convertFromHookEvents(input.Events, in)

	out := new(hook)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertHook(out), res, err
}

func (s *RepositoryService) CreateStatus(context.Context, string, string, *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *RepositoryService) UpdateHook(ctx context.Context, repo, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks/%s", repo, id)
	in := new(hook)
	// 1: signature
	in.EncryptionType = 1
	in.Password = input.Secret
	in.URL = input.Target
	convertFromHookEvents(input.Events, in)

	out := new(hook)
	res, err := s.client.do(ctx, "PATCH", path, in, out)
	return convertHook(out), res, err
}

func (s *RepositoryService) DeleteHook(ctx context.Context, repo, id string) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks/%s", repo, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

type repository struct {
	ID    int `json:"id"`
	Owner struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
	Namespace     namespace `json:"namespace"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	HumanName     string    `json:"human_name"`
	Path          string    `json:"path"`
	Public        bool      `json:"public"`
	Private       bool      `json:"private"`
	Internal      bool      `json:"internal"`
	Fork          bool      `json:"fork"`
	URL           string    `json:"url"`
	HtmlURL       string    `json:"html_url"`
	SshURL        string    `json:"ssh_url"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Permission    struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permission"`
}

type hook struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	// password encryption type, 0: password, 1: signature
	EncryptionType      int    `json:"encryption_type"`
	Password            string `json:"password"`
	ProjectID           int    `json:"project_id"`
	Result              string `json:"result"`
	ResultCode          int    `json:"result_code"`
	PushEvents          bool   `json:"push_events"`
	TagPushEvents       bool   `json:"tag_push_events"`
	IssuesEvents        bool   `json:"issues_events"`
	NoteEvents          bool   `json:"note_events"`
	MergeRequestsEvents bool   `json:"merge_requests_events"`
}

type namespace struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	HtmlURL string `json:"html_url"`
}

func convertRepositoryList(from []*repository) []*scm.Repository {
	to := []*scm.Repository{}
	for _, v := range from {
		to = append(to, convertRepository(v))
	}
	return to
}
func convertRepository(from *repository) *scm.Repository {
	return &scm.Repository{
		ID:        strconv.Itoa(from.ID),
		Name:      from.Path,
		Namespace: from.Namespace.Path,
		Perm: &scm.Perm{
			Push:  from.Permission.Push,
			Pull:  from.Permission.Pull,
			Admin: from.Permission.Admin,
		},
		Link:     from.HtmlURL,
		Branch:   from.DefaultBranch,
		Private:  from.Private,
		Clone:    from.HtmlURL,
		CloneSSH: from.SshURL,
		Created:  from.CreatedAt,
		Updated:  from.UpdatedAt,
	}
}

func convertHookList(from []*hook) []*scm.Hook {
	to := []*scm.Hook{}
	for _, v := range from {
		to = append(to, convertHook(v))
	}
	return to
}
func convertHook(from *hook) *scm.Hook {
	return &scm.Hook{
		ID:         strconv.Itoa(from.ID),
		Active:     true,
		Target:     from.URL,
		Events:     convertHookEvent(from),
		SkipVerify: true,
	}
}

func convertHookEvent(from *hook) []string {
	var events []string
	if from.PushEvents {
		events = append(events, "push")
	}
	if from.TagPushEvents {
		events = append(events, "tag_push")
	}
	if from.IssuesEvents {
		events = append(events, "issues")
	}
	if from.NoteEvents {
		events = append(events, "note")
	}
	if from.MergeRequestsEvents {
		events = append(events, "merge_requests")
	}
	return events
}

// convertFromHookEvents not support: Branch, Deployment
func convertFromHookEvents(from scm.HookEvents, to *hook) {
	if from.Push {
		to.PushEvents = true
	}
	if from.PullRequest {
		to.MergeRequestsEvents = true
	}
	if from.Issue {
		to.IssuesEvents = true
	}
	if from.IssueComment || from.PullRequestComment || from.ReviewComment {
		to.NoteEvents = true
	}
	if from.Tag {
		to.TagPushEvents = true
	}
}
