package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/stash"
	"github.com/drone/go-scm/scm/transport"
)

func TestListRepos(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Acceptance test")
	}
	client, _ = stash.New(endpoint)
	client.Client = &http.Client{
		Transport: &transport.BasicAuth{
			Username: username,
			Password: token,
		},
	}

	repos, response, listerr := client.Repositories.List(context.Background(), scm.ListOptions{})
	if listerr != nil {
		t.Errorf("List Repos got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("List Repos did not get a 200 back %v", response.Status)
	}

	if len(repos) == 0 {
		t.Errorf("Got Empty repo list")
	}

}
