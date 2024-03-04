// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/null"
)

type pullService struct {
	*issueService
}

func (s *pullService) Find(ctx context.Context, repo string, number int) (*scm.PullRequest, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-request?view=azure-devops-rest-6.0
	endpoint := fmt.Sprintf("%s/%s/_apis/git/repositories/%s/pullrequests/%d?api-version=6.0",
		s.client.owner, s.client.project, repo, number)
	out := new(pr)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListCommits(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-commits/get-pull-request-commits?view=azure-devops-rest-6.0
	endpoint := fmt.Sprintf("%s/%s/_apis/git/repositories/%s/pullRequests/%d/commits?api-version=6.0",
		s.client.owner, s.client.project, repo, number)
	out := new(commitList)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return convertCommitList(out.Value), res, err
}

func (s *pullService) Merge(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/create?view=azure-devops-rest-6.0
	endpoint := fmt.Sprintf("%s/%s/_apis/git/repositories/%s/pullrequests?api-version=6.0", s.client.owner, s.client.project, repo)
	in := &prInput{
		Title:         input.Title,
		Description:   input.Body,
		SourceRefName: scm.ExpandRef(input.Source, "refs/heads"),
		TargetRefName: scm.ExpandRef(input.Target, "refs/heads"),
	}
	out := new(pr)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertPullRequest(out), res, err
}

type prInput struct {
	SourceRefName string `json:"sourceRefName"`
	TargetRefName string `json:"targetRefName"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Reviewers     []struct {
		ID string `json:"id"`
	} `json:"reviewers"`
}

type pr struct {
	Repository struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		URL     string `json:"url"`
		Project struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			URL         string `json:"url"`
			State       string `json:"state"`
			Revision    int    `json:"revision"`
		} `json:"project"`
		RemoteURL string `json:"remoteUrl"`
	} `json:"repository"`
	PullRequestID int    `json:"pullRequestId"`
	CodeReviewID  int    `json:"codeReviewId"`
	Status        string `json:"status"`
	CreatedBy     struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		UniqueName  string `json:"uniqueName"`
		URL         string `json:"url"`
		ImageURL    string `json:"imageUrl"`
	} `json:"createdBy"`
	CreationDate          time.Time   `json:"creationDate"`
	ClosedDate            null.String `json:"closedDate"`
	Title                 string      `json:"title"`
	Description           string      `json:"description"`
	SourceRefName         string      `json:"sourceRefName"`
	TargetRefName         string      `json:"targetRefName"`
	MergeStatus           string      `json:"mergeStatus"`
	MergeID               string      `json:"mergeId"`
	LastMergeSourceCommit struct {
		CommitID string `json:"commitId"`
		URL      string `json:"url"`
	} `json:"lastMergeSourceCommit"`
	LastMergeTargetCommit struct {
		CommitID string `json:"commitId"`
		URL      string `json:"url"`
	} `json:"lastMergeTargetCommit"`
	Reviewers []struct {
		ReviewerURL string `json:"reviewerUrl"`
		Vote        int    `json:"vote"`
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		UniqueName  string `json:"uniqueName"`
		URL         string `json:"url"`
		ImageURL    string `json:"imageUrl"`
	} `json:"reviewers"`
	URL   string `json:"url"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Repository struct {
			Href string `json:"href"`
		} `json:"repository"`
		WorkItems struct {
			Href string `json:"href"`
		} `json:"workItems"`
		SourceBranch struct {
			Href string `json:"href"`
		} `json:"sourceBranch"`
		TargetBranch struct {
			Href string `json:"href"`
		} `json:"targetBranch"`
		SourceCommit struct {
			Href string `json:"href"`
		} `json:"sourceCommit"`
		TargetCommit struct {
			Href string `json:"href"`
		} `json:"targetCommit"`
		CreatedBy struct {
			Href string `json:"href"`
		} `json:"createdBy"`
		Iterations struct {
			Href string `json:"href"`
		} `json:"iterations"`
	} `json:"_links"`
	SupportsIterations bool   `json:"supportsIterations"`
	ArtifactID         string `json:"artifactId"`
}

func convertPullRequest(from *pr) *scm.PullRequest {
	return &scm.PullRequest{
		Number: from.PullRequestID,
		Title:  from.Title,
		Body:   from.Description,
		Sha:    from.LastMergeSourceCommit.CommitID,
		Source: scm.TrimRef(from.SourceRefName),
		Target: scm.TrimRef(from.TargetRefName),
		Link:   from.URL,
		Closed: from.ClosedDate.Valid,
		Merged: from.Status == "completed",
		Ref:    fmt.Sprintf("refs/pull/%d/merge", from.PullRequestID),
		Head: scm.Reference{
			Sha: from.LastMergeSourceCommit.CommitID,
		},
		Base: scm.Reference{

			Sha: from.LastMergeTargetCommit.CommitID,
		},
		Author: scm.User{
			Login:  from.CreatedBy.UniqueName,
			Avatar: from.CreatedBy.ImageURL,
		},
		Created: from.CreationDate,
	}
}
