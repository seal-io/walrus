// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"fmt"
	"net/url"

	"github.com/drone/go-scm/scm"
)

// RepositoryService implements the repository service for
// the GitHub driver.
type RepositoryService struct {
	client *wrapper
}

// Find returns the repository by name.
func (s *RepositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/repositories/get?view=azure-devops-rest-4.1
	if s.client.project == "" {
		return nil, nil, ProjectRequiredError()
	}
	endpoint := fmt.Sprintf("%s/%s/_apis/git/repositories/%s?api-version=6.0", s.client.owner, s.client.project, repo)

	out := new(repository)
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertRepository(out), res, err
}

// FindHook returns a repository hook.
func (s *RepositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// FindPerms returns the repository permissions.
func (s *RepositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// List returns the user repository list.
func (s *RepositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/repositories/list?view=azure-devops-rest-6.0
	var endpoint string
	if s.client.project == "" {
		endpoint = fmt.Sprintf("%s/_apis/git/repositories?api-version=6.0", s.client.owner)
	} else {
		endpoint = fmt.Sprintf("%s/%s/_apis/git/repositories?api-version=6.0", s.client.owner, s.client.project)
	}

	out := new(repositories)
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertRepositoryList(out), res, err
}

// ListV2 returns the user repository list.
func (s *RepositoryService) ListV2(ctx context.Context, opts scm.RepoListOptions) ([]*scm.Repository, *scm.Response, error) {
	// Azure does not support search filters, hence calling List api without search filtering
	return s.List(ctx, opts.ListOptions)
}

// ListHooks returns a list or repository hooks.
func (s *RepositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/hooks/subscriptions/list?view=azure-devops-rest-6.0
	if s.client.project == "" {
		return nil, nil, ProjectRequiredError()
	}
	projectID, projErr := s.getProjectIDFromProjectName(ctx, s.client.project)
	if projErr != nil {
		return nil, nil, fmt.Errorf("ListHooks was unable to look up the project's projectID, %s", projErr)
	}
	endpoint := fmt.Sprintf("%s/_apis/hooks/subscriptions?api-version=6.0", s.client.owner)
	out := new(subscriptions)
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertHookList(out.Value, projectID, repo), res, err
}

// ListStatus returns a list of commit statuses.
func (s *RepositoryService) ListStatus(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// CreateHook creates a new repository webhook.
func (s *RepositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/hooks/subscriptions/create?view=azure-devops-rest-6.0
	if s.client.project == "" {
		return nil, nil, ProjectRequiredError()
	}
	endpoint := fmt.Sprintf("%s/_apis/hooks/subscriptions?api-version=6.0", s.client.owner)
	in := new(subscription)
	in.Status = "enabled"
	in.PublisherID = "tfs"
	in.ResourceVersion = "1.0"
	in.ConsumerID = "webHooks"
	in.ConsumerActionID = "httpRequest"
	// we do not support scm hookevents, only native events
	if input.NativeEvents == nil {
		return nil, nil, fmt.Errorf("CreateHook, You must pass at least one native event")
	}
	if len(input.NativeEvents) > 1 {
		return nil, nil, fmt.Errorf("CreateHook, Azure only allows the creation of a single hook at a time %v", input.NativeEvents)
	}
	in.EventType = input.NativeEvents[0]
	// publisher
	projectID, projErr := s.getProjectIDFromProjectName(ctx, s.client.project)
	if projErr != nil {
		return nil, nil, fmt.Errorf("CreateHook was unable to look up the project's projectID, %s", projErr)
	}
	in.PublisherInputs.ProjectID = projectID
	in.PublisherInputs.Repository = repo
	// consumer
	in.ConsumerInputs.URL = input.Target
	if input.SkipVerify {
		in.ConsumerInputs.AcceptUntrustedCerts = "enabled"
	}
	// with version 1.0, azure provides incomplete data for issue-comment
	if in.EventType == "ms.vss-code.git-pullrequest-comment-event" {
		in.ResourceVersion = "2.0"
	}
	out := new(subscription)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertHook(out), res, err
}

// CreateStatus creates a new commit status.
func (s *RepositoryService) CreateStatus(ctx context.Context, repo, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// CreateDeployStatus creates a new deployment status.
func (s *RepositoryService) CreateDeployStatus(ctx context.Context, repo string, input *scm.DeployStatus) (*scm.DeployStatus, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// UpdateHook updates a repository webhook.
func (s *RepositoryService) UpdateHook(ctx context.Context, repo, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// DeleteHook deletes a repository webhook.
func (s *RepositoryService) DeleteHook(ctx context.Context, repo, id string) (*scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/hooks/subscriptions/delete?view=azure-devops-rest-6.0
	if s.client.project == "" {
		return nil, ProjectRequiredError()
	}
	endpoint := fmt.Sprintf("%s/_apis/hooks/subscriptions/%s?api-version=6.0", s.client.owner, id)
	return s.client.do(ctx, "DELETE", endpoint, nil, nil)
}

// helper function to return the projectID from the project name
func (s *RepositoryService) getProjectIDFromProjectName(ctx context.Context, projectName string) (string, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/list?view=azure-devops-rest-6.0
	projectName, err := url.PathUnescape(projectName)
	if err != nil {
		return "", fmt.Errorf("unable to unscape project: %s", projectName)
	}

	endpoint := fmt.Sprintf("%s/_apis/projects?api-version=6.0", s.client.owner)
	type projects struct {
		Count int64 `json:"count"`
		Value []struct {
			Description string `json:"description"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			State       string `json:"state"`
			URL         string `json:"url"`
		} `json:"value"`
	}

	out := new(projects)
	response, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	if err != nil {
		fmt.Println(response)
		return "", fmt.Errorf("failed to list projects: %s", err)
	}
	for _, v := range out.Value {
		if v.Name == projectName {
			return v.ID, nil
		}
	}
	return "", fmt.Errorf("failed to find project id for %s", projectName)
}

type repositories struct {
	Count int64         `json:"count"`
	Value []*repository `json:"value"`
}

type repository struct {
	DefaultBranch string `json:"defaultBranch"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Project       struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		State string `json:"state"`
		URL   string `json:"url"`
	} `json:"project"`
	RemoteURL string `json:"remoteUrl"`
	URL       string `json:"url"`
}

type subscriptions struct {
	Count int64           `json:"count"`
	Value []*subscription `json:"value"`
}

type subscription struct {
	ActionDescription string `json:"actionDescription"`
	ConsumerActionID  string `json:"consumerActionId"`
	ConsumerID        string `json:"consumerId"`
	ConsumerInputs    struct {
		AccountName          string `json:"accountName,omitempty"`
		AcceptUntrustedCerts string `json:"acceptUntrustedCerts,omitempty"`
		AddToTop             string `json:"addToTop,omitempty"`
		APIToken             string `json:"apiToken,omitempty"`
		BoardID              string `json:"boardId,omitempty"`
		BuildName            string `json:"buildName,omitempty"`
		BuildParameterized   string `json:"buildParameterized,omitempty"`
		FeedID               string `json:"feedId,omitempty"`
		ListID               string `json:"listId,omitempty"`
		PackageSourceID      string `json:"packageSourceId,omitempty"`
		Password             string `json:"password,omitempty"`
		ServerBaseURL        string `json:"serverBaseUrl,omitempty"`
		URL                  string `json:"url,omitempty"`
		UserToken            string `json:"userToken,omitempty"`
		Username             string `json:"username,omitempty"`
	} `json:"consumerInputs"`
	CreatedBy struct {
		ID string `json:"id"`
	} `json:"createdBy"`
	CreatedDate      string `json:"createdDate"`
	EventDescription string `json:"eventDescription"`
	EventType        string `json:"eventType"`
	ID               string `json:"id"`
	ModifiedBy       struct {
		ID string `json:"id"`
	} `json:"modifiedBy"`
	ModifiedDate     string `json:"modifiedDate"`
	ProbationRetries int64  `json:"probationRetries"`
	PublisherID      string `json:"publisherId"`
	PublisherInputs  struct {
		AreaPath          string `json:"areaPath,omitempty"`
		Branch            string `json:"branch,omitempty"`
		BuildStatus       string `json:"buildStatus,omitempty"`
		ChangedFields     string `json:"changedFields,omitempty"`
		CommentPattern    string `json:"commentPattern,omitempty"`
		DefinitionName    string `json:"definitionName,omitempty"`
		HostID            string `json:"hostId,omitempty"`
		Path              string `json:"path,omitempty"`
		ProjectID         string `json:"projectId,omitempty"`
		Repository        string `json:"repository,omitempty"`
		TfsSubscriptionID string `json:"tfsSubscriptionId,omitempty"`
		WorkItemType      string `json:"workItemType,omitempty"`
	} `json:"publisherInputs"`
	ResourceVersion string `json:"resourceVersion"`
	Status          string `json:"status"`
	URL             string `json:"url"`
}

// helper function to convert from the gogs repository list to
// the common repository structure.
func convertRepositoryList(from *repositories) []*scm.Repository {
	to := []*scm.Repository{}
	for _, v := range from.Value {
		to = append(to, convertRepository(v))
	}
	return to
}

// helper function to convert from the gogs repository structure
// to the common repository structure.
func convertRepository(from *repository) *scm.Repository {
	return &scm.Repository{
		ID:     from.ID,
		Name:   from.Name,
		Link:   from.URL,
		Branch: scm.TrimRef(from.DefaultBranch),
	}
}

func convertHookList(from []*subscription, projectFilter string, repositoryFilter string) []*scm.Hook {
	to := []*scm.Hook{}
	for _, v := range from {
		if repositoryFilter != "" && projectFilter == v.PublisherInputs.ProjectID && repositoryFilter == v.PublisherInputs.Repository {
			to = append(to, convertHook(v))
		}
	}
	return to
}

func convertHook(from *subscription) *scm.Hook {
	returnVal := &scm.Hook{
		ID: from.ID,

		Active:     from.Status == "enabled",
		Target:     from.ConsumerInputs.URL,
		Events:     []string{from.EventType},
		SkipVerify: from.ConsumerInputs.AcceptUntrustedCerts == "true",
	}

	return returnVal
}
