package integration

import (
	"context"
	"os"

	"github.com/drone/go-scm/scm"
)

var (
	client *scm.Client
	token  = os.Getenv("AZURE_TOKEN")

	organization = "tphoney"
	project      = "test_project"
	repoID       = "fde2d21f-13b9-4864-a995-83329045289a"
)

func GetCurrentCommitOfBranch(client *scm.Client, branch string) (string, error) {
	commit, _, err := client.Contents.List(context.Background(), repoID, "", "main", scm.ListOptions{})
	if err != nil {
		return "", err
	}
	return commit[0].Sha, nil
}
