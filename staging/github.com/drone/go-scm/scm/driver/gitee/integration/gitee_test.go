// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"net/http"
	"os"
	"testing"

	"github.com/drone/go-scm/scm/driver/gitee"
	"github.com/drone/go-scm/scm/transport"
)

func TestGitee(t *testing.T) {
	accessToken := os.Getenv("GITEE_TOKEN")
	if accessToken == "" {
		t.Skipf("missing GITEE_TOKEN environment variable")
		return
	}

	client := gitee.NewDefault()
	client.Client = &http.Client{
		Transport: &transport.BearerToken{
			Token: accessToken,
		},
	}

	t.Run("Contents", testContents(client))
	t.Run("Git", testGit(client))
	t.Run("Issues", testIssues(client))
	t.Run("Organizations", testOrgs(client))
	t.Run("PullRequests", testPullRequests(client))
	t.Run("Repositories", testRepos(client))
	t.Run("Users", testUsers(client))
}
