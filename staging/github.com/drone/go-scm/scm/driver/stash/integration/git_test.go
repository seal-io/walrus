package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/stash"
	"github.com/drone/go-scm/scm/transport"
)

func TestCreateBranch(t *testing.T) {
	if token == "" || username == "" {
		t.Skip("Skipping, Acceptance test. Missing BITBUCKET_SERVER_TOKEN or BITBUCKET_USERNAME")
	}
	client, _ = stash.New(endpoint)
	client.Client = &http.Client{
		Transport: &transport.BasicAuth{
			Username: username,
			Password: token,
		},
	}

	commitId, _ := GetCurrentCommitOfBranch(client, "master")
	input := &scm.ReferenceInput{
		Name: "test_branch",
		Sha:  commitId,
	}
	response, listerr := client.Git.CreateBranch(context.Background(), repoID, input)
	if listerr != nil {
		t.Errorf("CreateBranch got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("CreateBranch did not get a 200 back %v", response.Status)
	}
}

func TestGetLatestCommitOfBranch(t *testing.T) {
	if token == "" || username == "" {
		t.Skip("Skipping, Acceptance test. Missing BITBUCKET_SERVER_TOKEN or BITBUCKET_USERNAME")
	}
	client, _ = stash.New(endpoint)
	client.Client = &http.Client{
		Transport: &transport.BasicAuth{
			Username: username,
			Password: token,
		},
	}

	commits, response, err := client.Git.ListCommits(context.Background(), repoID, scm.CommitListOptions{Ref: "master", Path: "README"})

	if err != nil {
		t.Errorf("GetLatestCommitOfFile got an error %v", err)
	} else {
		if response.Status != http.StatusOK {
			t.Errorf("GetLatestCommitOfFile did not get a 200 back %v", response.Status)
		}

		if commits[0].Sha != "2cc4dbe084f0d66761318b305c408cb0ea300c9a" {
			t.Errorf("Got the commitId %s instead of the top commit of the file", commits[0].Sha)
		}
	}
}

func TestGetLatestCommitOfNonDefaultBranch(t *testing.T) {
	if token == "" || username == "" {
		t.Skip("Skipping, Acceptance test. Missing BITBUCKET_SERVER_TOKEN or BITBUCKET_USERNAME")
	}
	client, _ = stash.New(endpoint)
	client.Client = &http.Client{
		Transport: &transport.BasicAuth{
			Username: username,
			Password: token,
		},
	}

	commits, response, err := client.Git.ListCommits(context.Background(), repoID, scm.CommitListOptions{Ref: "main", Path: "do-not-touch.txt"})

	if err != nil {
		t.Errorf("GetLatestCommitOfFile got an error %v", err)
	} else {
		if response.Status != http.StatusOK {
			t.Errorf("GetLatestCommitOfFile did not get a 200 back %v", response.Status)
		}

		if commits[0].Sha != "76fb1762048a277596d3fa330b3da140cd12d361" {
			t.Errorf("Got the commitId %s instead of the top commit of the file", commits[0].Sha)
		}
	}
}

func TestGetLatestCommitOfBranchWhenNoRefPassed(t *testing.T) {
	if token == "" || username == "" {
		t.Skip("Skipping, Acceptance test. Missing BITBUCKET_SERVER_TOKEN or BITBUCKET_USERNAME")
	}
	client, _ = stash.New(endpoint)
	client.Client = &http.Client{
		Transport: &transport.BasicAuth{
			Username: username,
			Password: token,
		},
	}

	commits, response, err := client.Git.ListCommits(context.Background(), repoID, scm.CommitListOptions{Path: "README"})

	if err != nil {
		t.Errorf("GetLatestCommitOfFile got an error %v", err)
	} else {
		if response.Status != http.StatusOK {
			t.Errorf("GetLatestCommitOfFile did not get a 200 back %v", response.Status)
		}

		if commits[0].Sha != "2cc4dbe084f0d66761318b305c408cb0ea300c9a" {
			t.Errorf("Got the commitId %s instead of the top commit of the file", commits[0].Sha)
		}
	}
}
