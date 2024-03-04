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

func TestListBranches(t *testing.T) {
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
	references, response, listerr := client.Git.ListBranches(context.Background(), repoID, scm.ListOptions{})
	if listerr != nil {
		t.Errorf("ListBranches got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("ListBranches did not get a 200 back %v", response.Status)
	}
	if len(references) < 1 {
		t.Errorf("ListBranches  should have at least 1 branch %d", len(references))
	}
	if references[0].Sha == "" {
		t.Errorf("ListBranches first entry did not get a sha back %v", references[0].Sha)
	}
}

func TestCreateBranch(t *testing.T) {
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
	currentCommit, commitErr := GetCurrentCommitOfBranch(client, "main")
	if commitErr != nil {
		t.Errorf("we got an error %v", commitErr)
	}
	input := &scm.ReferenceInput{
		Name: "test_branch",
		Sha:  currentCommit,
	}
	response, listerr := client.Git.CreateBranch(context.Background(), repoID, input)
	if listerr != nil {
		t.Errorf("CreateBranch got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("CreateBranch did not get a 200 back %v", response.Status)
	}

}

func TestFindCommit(t *testing.T) {
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
	currentCommit, commitErr := GetCurrentCommitOfBranch(client, "main")
	if commitErr != nil {
		t.Errorf("we got an error %v", commitErr)
	}
	commit, response, listerr := client.Git.FindCommit(context.Background(), repoID, currentCommit)
	if listerr != nil {
		t.Errorf("FindCommit got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("FindCommit did not get a 200 back %v", response.Status)
	}
	if commit.Author.Name == "" {
		t.Errorf("There is no author %v", commit.Author)
	}
}

func TestListCommits(t *testing.T) {
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
	commits, response, listerr := client.Git.ListCommits(context.Background(), repoID, scm.CommitListOptions{})
	if listerr != nil {
		t.Errorf("ListCommits  got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("ListCommitsdid not get a 200 back %v", response.Status)
	}
	if len(commits) < 1 {
		t.Errorf("Contents.List there should be at least 1 commit %d", len(commits))
	}
	if commits[0].Sha == "" {
		t.Errorf("Contents.List first entry did not get a sha back %v", commits[0].Sha)
	}
}

func TestCompareChanges(t *testing.T) {
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
	// get all the commits
	commits, _, err := client.Git.ListCommits(context.Background(), repoID, scm.CommitListOptions{})
	if err != nil {
		t.Errorf("we got an error %v", err)
	}
	// compare the last two commits
	changes, response, listerr := client.Git.CompareChanges(context.Background(), repoID, commits[10].Sha, commits[0].Sha, scm.ListOptions{})
	if listerr != nil {
		t.Errorf("CompareChanges got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("CompareChanges did not get a 200 back %v", response.Status)
	}
	if len(changes) == 0 {
		t.Errorf("There is at least one change %d", len(changes))
	}
}
