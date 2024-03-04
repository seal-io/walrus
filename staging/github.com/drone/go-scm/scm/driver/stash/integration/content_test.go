package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/stash"
	"github.com/drone/go-scm/scm/transport"
)

func TestCreateUpdateDeleteFileStash(t *testing.T) {
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
	// get latest commit first
	currentCommit, commitErr := GetCurrentCommitOfBranch(client, "master")
	if commitErr != nil {
		t.Errorf("we got an error %v", commitErr)
	}
	// create a new file
	createParams := scm.ContentParams{
		Message: "go-scm create crud file",
		Data:    []byte("hello"),
		Branch:  "master",
		Sha:     currentCommit,
	}
	createResponse, createErr := client.Contents.Create(context.Background(), repoID, "README5", &createParams)
	if createErr != nil {
		t.Errorf("Contents.Create we got an error %v", createErr)
	}
	if createResponse.Status != http.StatusOK {
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
		Branch:  "master",
		Sha:     currentCommit,
	}
	updateResponse, updateErr := client.Contents.Update(context.Background(), repoID, "README5", &updateParams)
	if updateErr != nil {
		t.Errorf("Contents.Update we got an error %v", updateErr)
	}
	if updateResponse.Status != http.StatusOK {
		t.Errorf("Contents.Update we did not get a 201 back %v", updateResponse.Status)
	}
}
