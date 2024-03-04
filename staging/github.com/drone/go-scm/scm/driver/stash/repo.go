// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/drone/go-scm/scm"
)

type repository struct {
	Slug          string `json:"slug"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ScmID         string `json:"scmId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Forkable      bool   `json:"forkable"`
	Project       struct {
		Key    string `json:"key"`
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Public bool   `json:"public"`
		Type   string `json:"type"`
		Links  struct {
			Self []link `json:"self"`
		} `json:"links"`
	} `json:"project"`
	Public bool `json:"public"`
	Links  struct {
		Clone []link `json:"clone"`
		Self  []link `json:"self"`
	} `json:"links"`
}

type repositories struct {
	pagination
	Values []*repository `json:"values"`
}

type link struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type perms struct {
	Values []*perm `json:"values"`
}

type perm struct {
	Permissions string `json:"permission"`
}

type hooks struct {
	pagination
	Values []*hook `json:"values"`
}

type hook struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	CreatedDate int64    `json:"createdDate"`
	UpdatedDate int64    `json:"updatedDate"`
	Events      []string `json:"events"`
	URL         string   `json:"url"`
	Active      bool     `json:"active"`
	Config      struct {
		Secret string `json:"secret"`
	} `json:"configuration"`
}

type hookInput struct {
	Name   string   `json:"name"`
	Events []string `json:"events"`
	URL    string   `json:"url"`
	Active bool     `json:"active"`
	Config struct {
		Secret string `json:"secret,omitempty"`
	} `json:"configuration"`
}

type status struct {
	State string `json:"state"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Desc  string `json:"description"`
}

type repositoryService struct {
	client *wrapper
}

// Find returns the repository by name.
func (s *repositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s", namespace, name)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	outputRepo := convertRepository(out)

	branch := new(branch)
	pathBranch := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/branches/default", namespace, name)
	_, errBranch := s.client.do(ctx, "GET", pathBranch, nil, branch)
	if errBranch == nil {
		outputRepo.Branch = branch.DisplayID
	}
	if err == nil {
		err = errBranch
	}

	return outputRepo, res, err
}

// FindHook returns a repository hook.
func (s *repositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks/%s", namespace, name, id)
	out := new(hook)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertHook(out), res, err
}

// FindPerms returns the repository permissions.
func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	// HACK: test if the user has read access to the repository.
	_, _, err := s.Find(ctx, repo)
	if err != nil {
		return &scm.Perm{
			Pull:  false,
			Push:  false,
			Admin: false,
		}, nil, nil
	}

	// HACK: test if the user has admin access to the repository.
	_, _, err = s.ListHooks(ctx, repo, scm.ListOptions{})
	if err == nil {
		return &scm.Perm{
			Pull:  true,
			Push:  true,
			Admin: true,
		}, nil, nil
	}
	// HACK: test if the user has write access to the repository.
	namespace, _ := scm.Split(repo)
	repos, _, _ := s.listWrite(ctx, repo)
	for _, repo := range repos {
		if repo.Namespace == namespace {
			return &scm.Perm{
				Pull:  true,
				Push:  true,
				Admin: false,
			}, nil, nil
		}
	}

	return &scm.Perm{
		Pull:  true,
		Push:  false,
		Admin: false,
	}, nil, nil
}

// List returns the user repository list.
func (s *repositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("rest/api/1.0/repos?%s", encodeListRoleOptions(opts))
	out := new(repositories)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if res != nil && !out.pagination.LastPage.Bool {
		res.Page.First = 1
		res.Page.Next = opts.Page + 1
	}
	return convertRepositoryList(out), res, err
}

// ListV2 returns the user repository list based on the searchTerm passed.
func (s *repositoryService) ListV2(ctx context.Context, opts scm.RepoListOptions) ([]*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("rest/api/1.0/repos?%s", encodeRepoListOptions(opts))
	out := new(repositories)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if res != nil && !out.pagination.LastPage.Bool {
		res.Page.First = 1
		res.Page.Next = opts.ListOptions.Page + 1
	}
	return convertRepositoryList(out), res, err
}

// listWrite returns the user repository list.
func (s *repositoryService) listWrite(ctx context.Context, repo string) ([]*scm.Repository, *scm.Response, error) {
	_, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/repos?size=1000&permission=REPO_WRITE&name=%s", name)
	out := new(repositories)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepositoryList(out), res, err
}

// ListHooks returns a list or repository hooks.
func (s *repositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks?%s", namespace, name, encodeListOptions(opts))
	out := new(hooks)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	if res != nil && !out.pagination.LastPage.Bool {
		res.Page.First = 1
		res.Page.Next = opts.Page + 1
	}
	return convertHookList(out), res, err
}

// ListStatus returns a list of commit statuses.
func (s *repositoryService) ListStatus(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// CreateHook creates a new repository webhook.
func (s *repositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks", namespace, name)
	in := new(hookInput)
	in.URL = input.Target
	in.Active = true
	in.Name = input.Name
	in.Config.Secret = input.Secret
	in.Events = append(
		input.NativeEvents,
		convertFromHookEvents(input.Events)...,
	)
	out := new(hook)
	res, err := s.client.do(ctx, "POST", path, in, out)
	if err != nil && isUnknownHookEvent(err) {
		downgradeHookInput(in)
		res, err = s.client.do(ctx, "POST", path, in, out)
	}
	return convertHook(out), res, err
}

// CreateStatus creates a new commit status.
func (s *repositoryService) CreateStatus(ctx context.Context, repo, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	path := fmt.Sprintf("rest/build-status/1.0/commits/%s", ref)
	in := status{
		State: convertFromState(input.State),
		Key:   input.Label,
		Name:  input.Label,
		URL:   input.Target,
		Desc:  input.Desc,
	}
	res, err := s.client.do(ctx, "POST", path, in, nil)
	return &scm.Status{
		State:  input.State,
		Label:  input.Label,
		Desc:   input.Desc,
		Target: input.Target,
	}, res, err
}

// UpdateHook updates new repository webhook.
func (s *repositoryService) UpdateHook(ctx context.Context, repo, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks/%s", namespace, name, id)
	in := new(hookInput)
	in.URL = input.Target
	in.Active = true
	in.Name = input.Name
	in.Config.Secret = input.Secret
	in.Events = append(
		input.NativeEvents,
		convertFromHookEvents(input.Events)...,
	)
	out := new(hook)
	res, err := s.client.do(ctx, "PUT", path, in, out)
	if err != nil && isUnknownHookEvent(err) {
		downgradeHookInput(in)
		res, err = s.client.do(ctx, "PUT", path, in, out)
	}
	return convertHook(out), res, err
}

// DeleteHook deletes a repository webhook.
func (s *repositoryService) DeleteHook(ctx context.Context, repo string, id string) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks/%s", namespace, name, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// helper function to convert from the gogs repository list to
// the common repository structure.
func convertRepositoryList(from *repositories) []*scm.Repository {
	to := []*scm.Repository{}
	for _, v := range from.Values {
		to = append(to, convertRepository(v))
	}
	return to
}

// helper function to convert from the gogs repository structure
// to the common repository structure.
func convertRepository(from *repository) *scm.Repository {
	return &scm.Repository{
		ID:        strconv.Itoa(from.ID),
		Name:      from.Slug,
		Namespace: from.Project.Key,
		Link:      extractSelfLink(from.Links.Self),
		Branch:    "master",
		Private:   !from.Public,
		CloneSSH:  extractLink(from.Links.Clone, "ssh"),
		Clone:     anonymizeLink(extractLink(from.Links.Clone, "http")),
	}
}

func extractLink(links []link, name string) (href string) {
	for _, link := range links {
		if link.Name == name {
			return link.Href
		}
	}
	return
}

func extractSelfLink(links []link) (href string) {
	for _, link := range links {
		return link.Href
	}
	return
}

func anonymizeLink(link string) (href string) {
	parsed, err := url.Parse(link)
	if err != nil {
		return link
	}
	parsed.User = nil
	return parsed.String()
}

func convertHookList(from *hooks) []*scm.Hook {
	to := []*scm.Hook{}
	for _, v := range from.Values {
		to = append(to, convertHook(v))
	}
	return to
}

func convertHook(from *hook) *scm.Hook {
	return &scm.Hook{
		ID:     strconv.Itoa(from.ID),
		Name:   from.Name,
		Active: from.Active,
		Target: from.URL,
		Events: from.Events,
	}
}

func convertFromHookEvents(from scm.HookEvents) []string {
	var events []string
	if from.Push || from.Branch || from.Tag {
		events = append(events, "repo:refs_changed")
	}
	if from.PullRequest {
		events = append(events, "pr:declined")
		events = append(events, "pr:modified")
		events = append(events, "pr:deleted")
		events = append(events, "pr:opened")
		events = append(events, "pr:merged")
		events = append(events, "pr:from_ref_updated")
	}
	if from.PullRequestComment {
		events = append(events, "pr:comment:added")
		events = append(events, "pr:comment:deleted")
		events = append(events, "pr:comment:edited")
	}
	return events
}

func isUnknownHookEvent(err error) bool {
	return strings.Contains(err.Error(), "pr:from_ref_updated is unknown")
}

func downgradeHookInput(in *hookInput) {
	var events []string
	for _, event := range in.Events {
		if event != "pr:from_ref_updated" {
			events = append(events, event)
		}
	}
	in.Events = events
}

func convertFromState(from scm.State) string {
	switch from {
	case scm.StatePending, scm.StateRunning:
		return "INPROGRESS"
	case scm.StateSuccess:
		return "SUCCESSFUL"
	default:
		return "FAILED"
	}
}

func convertState(from string) scm.State {
	switch from {
	case "FAILED":
		return scm.StateFailure
	case "INPROGRESS":
		return scm.StatePending
	case "SUCCESSFUL":
		return scm.StateSuccess
	default:
		return scm.StateUnknown
	}
}
