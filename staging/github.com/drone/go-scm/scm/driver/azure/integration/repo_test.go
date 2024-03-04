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

func TestListRepos(t *testing.T) {
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
	references, response, listerr := client.Repositories.List(context.Background(), scm.ListOptions{})
	if listerr != nil {
		t.Errorf("List got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("List did not get a 200 back %v", response.Status)
	}
	if len(references) < 1 {
		t.Errorf("List should have at least 1 repo %d", len(references))
	}
}

func TestListHooks(t *testing.T) {
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
	hooks, response, listerr := client.Repositories.ListHooks(context.Background(), repoID, scm.ListOptions{})
	if listerr != nil {
		t.Errorf("List got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("List did not get a 200 back %v", response.Status)
	}
	if len(hooks) < 1 {
		t.Errorf("List should have at least 1 hook %d", len(hooks))
	}
}

func TestCreateDeleteHooks(t *testing.T) {
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
	originalHooks, _, _ := client.Repositories.ListHooks(context.Background(), repoID, scm.ListOptions{})
	// create a new hook
	inputHook := &scm.HookInput{
		Name:         "web",
		NativeEvents: []string{"git.push"},
		Target:       "http://www.example.com/webhook",
	}
	outHook, createResponse, createErr := client.Repositories.CreateHook(context.Background(), repoID, inputHook)
	if createErr != nil {
		t.Errorf("Create got an error %v", createErr)
	}
	if createResponse.Status != http.StatusOK {
		t.Errorf("Create did not get a 200 back %v", createResponse.Status)
	}
	if len(outHook.Events) != 1 {
		t.Errorf("New hook has one event %d", len(outHook.Events))
	}
	// get the hooks again, and make sure the new hook is there
	afterCreate, _, _ := client.Repositories.ListHooks(context.Background(), repoID, scm.ListOptions{})
	if len(afterCreate) != len(originalHooks)+1 {
		t.Errorf("After create, the number of hooks is not correct %d. It should be %d", len(afterCreate), len(originalHooks)+1)
	}
	// delete the hook we created
	deleteResponse, deleteErr := client.Repositories.DeleteHook(context.Background(), repoID, outHook.ID)
	if deleteErr != nil {
		t.Errorf("Delete got an error %v", deleteErr)
	}
	if deleteResponse.Status != http.StatusNoContent {
		t.Errorf("Delete did not get a 204 back, got %v", deleteResponse.Status)
	}
	// get the hooks again, and make sure the new hook is gone
	afterDelete, _, _ := client.Repositories.ListHooks(context.Background(), repoID, scm.ListOptions{})
	if len(afterDelete) != len(originalHooks) {
		t.Errorf("After Delete, the number of hooks is not correct %d. It should be %d", len(afterDelete), len(originalHooks))
	}
}
