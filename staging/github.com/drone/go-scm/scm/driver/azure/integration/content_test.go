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

func TestContentsFind(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Integration test")
	}
	client = azure.NewDefault(organization, project)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			},
		},
	}
	content, response, err := client.Contents.Find(context.Background(), repoID, "README.md", "")
	if err != nil {
		t.Errorf("We got an error %v", err)
	}
	if content.Sha == "" {
		t.Errorf("we did not get a sha back %v", content.Sha)
	}
	if response.Status != http.StatusOK {
		t.Errorf("we did not get a 200 back %v", response.Status)
	}
}

func TestCreateUpdateDeleteFileAzure(t *testing.T) {
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
	// get latest commit first
	currentCommit, commitErr := GetCurrentCommitOfBranch(client, "main")
	if commitErr != nil {
		t.Errorf("we got an error %v", commitErr)
	}
	// create a new file
	createParams := scm.ContentParams{
		Message: "go-scm create crud file",
		Data:    []byte("hello"),
		Branch:  "refs/heads/main",
		Ref:     currentCommit,
	}
	createResponse, createErr := client.Contents.Create(context.Background(), repoID, "CRUD", &createParams)
	if createErr != nil {
		t.Errorf("Contents.Create we got an error %v", createErr)
	}
	if createResponse.Status != http.StatusCreated {
		t.Errorf("Contents.Create we did not get a 201 back %v", createResponse.Status)
	}
	// get latest commit first
	currentCommit, commitErr = GetCurrentCommitOfBranch(client, "main")
	if commitErr != nil {
		t.Errorf("we got an error %v", commitErr)
	}
	// update the file
	updateParams := scm.ContentParams{
		Message: "go-scm update crud file",
		Data:    []byte("updated test data"),
		Branch:  "refs/heads/main",
		Sha:     currentCommit,
	}
	updateResponse, updateErr := client.Contents.Update(context.Background(), repoID, "CRUD", &updateParams)
	if updateErr != nil {
		t.Errorf("Contents.Update we got an error %v", updateErr)
	}
	if updateResponse.Status != http.StatusCreated {
		t.Errorf("Contents.Update we did not get a 201 back %v", updateResponse.Status)
	}
	// get latest commit first
	currentCommit, commitErr = GetCurrentCommitOfBranch(client, "main")
	if commitErr != nil {
		t.Errorf("we got an error %v", commitErr)
	}
	// delete the file
	deleteParams := scm.ContentParams{
		Message: "go-scm delete crud file",
		Branch:  "refs/heads/main",
		Sha:     currentCommit,
	}
	deleteResponse, deleteErr := client.Contents.Delete(context.Background(), repoID, "CRUD", &deleteParams)
	if deleteErr != nil {
		t.Errorf("Contents.Delete we got an error %v", deleteErr)
	}
	if deleteResponse.Status != http.StatusCreated {
		t.Errorf("Contents.Delete we did not get a 201 back %v", deleteResponse.Status)
	}
}

func TestListFiles(t *testing.T) {
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
	contentInfo, listResponse, listerr := client.Contents.List(context.Background(),
		repoID, "", "", scm.ListOptions{})
	if listerr != nil {
		t.Errorf("Contents.List we got an error %v", listerr)
	}
	if listResponse.Status != http.StatusOK {
		t.Errorf("Contents.Delete we did not get a 200 back %v", listResponse.Status)
	}
	if 2 >= len(contentInfo) {
		t.Errorf("Contents.List there should be at least 2 files %d", len(contentInfo))
	}

}
