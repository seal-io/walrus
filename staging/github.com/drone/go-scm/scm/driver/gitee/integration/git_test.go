// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

//
// git sub-tests
//

func testGit(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Branches", testBranches(client))
		t.Run("Commits", testCommits(client))
		t.Run("Tags", testTags(client))
	}
}

//
// branch sub-tests
//

func testBranches(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testBranchFind(client))
		t.Run("List", testBranchList(client))
	}
}

func testBranchFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Git.FindBranch(context.Background(), "kit101/drone-yml-test", "feat-4")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("Branch", testBranch(result))
	}
}

func testBranchList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.ListOptions{}
		result, _, err := client.Git.ListBranches(context.Background(), "kit101/drone-yml-test", opts)
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty branch list")
		}
		for _, branch := range result {
			if branch.Name == "feat-4" {
				t.Run("Branch", testBranch(branch))
			}
		}
	}
}

//
// branch sub-tests
//

func testTags(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testTagFind(client))
		t.Run("List", testTagList(client))
	}
}

func testTagFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		tag, _, err := client.Git.FindTag(context.Background(), "kit101/drone-yml-test", "1.1")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("Tag", testTag(tag))
	}
}

func testTagList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.ListOptions{}
		result, _, err := client.Git.ListTags(context.Background(), "kit101/drone-yml-test", opts)
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty tag list")
		}
		for _, tag := range result {
			if tag.Name == "1.1" {
				t.Run("Tag", testTag(tag))
			}
		}
	}
}

//
// commit sub-tests
//

func testCommits(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testCommitFind(client))
		t.Run("List", testCommitList(client))
	}
}

func testCommitFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Git.FindCommit(context.Background(), "kit101/drone-yml-test", "e3c0ff4d5cef439ea11b30866fb1ed79b420801d")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("Commit", testCommit(result))
	}
}

func testCommitList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.CommitListOptions{
			Ref: "feat-4",
		}
		result, _, err := client.Git.ListCommits(context.Background(), "kit101/drone-yml-test", opts)
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty commit list")
		}
		for _, commit := range result {
			if commit.Sha == "e3c0ff4d5cef439ea11b30866fb1ed79b420801d" {
				t.Run("Commit", testCommit(commit))
			}
		}
	}
}

//
// struct sub-tests
//

func testBranch(branch *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := branch.Name, "feat-4"; got != want {
			t.Errorf("Want branch Name %q, got %q", want, got)
		}
		if got, want := branch.Sha, "2eac1cac02c325058cf959725c45b0612d3e8177"; got != want {
			t.Errorf("Want branch Avatar %q, got %q", want, got)
		}
	}
}

func testTag(tag *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := tag.Name, "1.1"; got != want {
			t.Errorf("Want tag Name %q, got %q", want, got)
		}
		if got, want := tag.Sha, "5e7876efb3468ff679410b82a72f7c002382d41e"; got != want {
			t.Errorf("Want tag Avatar %q, got %q", want, got)
		}
	}
}

func testCommit(commit *scm.Commit) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := commit.Message, "add/update/delete/rename\n"; got != want {
			t.Errorf("Want commit Message %q, got %q", want, got)
		}
		if got, want := commit.Sha, "e3c0ff4d5cef439ea11b30866fb1ed79b420801d"; got != want {
			t.Errorf("Want commit Sha %q, got %q", want, got)
		}
		if got, want := commit.Author.Name, "kit101"; got != want {
			t.Errorf("Want commit author Name %q, got %q", want, got)
		}
		if got, want := commit.Author.Email, "qkssk1711@163.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Author.Date.Unix(), int64(1629733553); got != want {
			t.Errorf("Want commit author Date %d, got %d", want, got)
		}
		if got, want := commit.Committer.Name, "kit101"; got != want {
			t.Errorf("Want commit author Name %q, got %q", want, got)
		}
		if got, want := commit.Committer.Email, "qkssk1711@163.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Committer.Date.Unix(), int64(1629733553); got != want {
			t.Errorf("Want commit author Date %d, got %d", want, got)
		}
	}
}
