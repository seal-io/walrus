package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/azure"
	"github.com/drone/go-scm/scm/transport"
)

func TestCreatePR(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Acceptance test")
	}
	client = azure.NewDefault(organization, project)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			},
		},
	}
	input := &scm.PullRequestInput{
		Title:  "test_pr",
		Body:   "test_pr_body",
		Source: "pr_branch",
		Target: "main",
	}
	outputPR, response, listerr := client.PullRequests.Create(context.Background(), repoID, input)
	if listerr != nil {
		t.Errorf("PullRequests.Create got an error %v", listerr)
	}
	if response.Status != http.StatusCreated {
		t.Errorf("PullRequests.Create did not get a 201 back %v", response.Status)
	}
	if outputPR.Title != "test_pr" {
		t.Errorf("PullRequests.Create does not have the correct title %v", outputPR.Title)
	}
}

func TestPullRequestFind(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Acceptance test")
	}
	client = azure.NewDefault(organization, project)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			},
		},
	}
	outputPR, response, err := client.PullRequests.Find(context.Background(), repoID, 1)
	if err != nil {
		t.Errorf("PullRequests.Find got an error %v", err)
	}
	if response.Status != http.StatusOK {
		t.Errorf("PullRequests.Find did not get a 200 back %v", response.Status)
	}
	if outputPR.Title != "test_pr" {
		t.Errorf("PullRequests.Find does not have the correct title %v", outputPR.Title)
	}
}

func TestPullRequestCommits(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Acceptance test")
	}
	client = azure.NewDefault(organization, project)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			},
		},
	}
	commits, response, err := client.PullRequests.ListCommits(context.Background(), repoID, 1, scm.ListOptions{})
	if err != nil {
		t.Errorf("PullRequests.ListCommits got an error %v", err)
	}
	if response.Status != http.StatusOK {
		t.Errorf("PullRequests.ListCommits did not get a 200 back %v", response.Status)
	}
	if len(commits) < 1 {
		t.Errorf("PullRequests.ListCommits there should be at least 1 commit %d", len(commits))
	}
	if commits[0].Sha == "" {
		t.Errorf("PullRequests.ListCommits first entry did not get a sha back %v", commits[0].Sha)
	}
}
