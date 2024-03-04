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
// pull request sub-tests
//

func testPullRequests(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("List", testPullRequestList(client))
		t.Run("Find", testPullRequestFind(client))
		t.Run("Changes", testPullRequestChanges(client))
		t.Run("Comments", testPullRequestComments(client))
		t.Run("Commits", testPullRequestCommitList(client))
	}
}

func testPullRequestList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.PullRequestListOptions{
			Open:   true,
			Closed: true,
		}
		result, _, err := client.PullRequests.List(context.Background(), "gitlab-org/testme", opts)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Errorf("Got empty pull request list")
		}
		for _, pr := range result {
			if pr.Number == 1 {
				t.Run("PullRequest", testPullRequest(pr))
			}
		}
	}
}

func testPullRequestFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.PullRequests.Find(context.Background(), "gitlab-org/testme", 1)
		if err != nil {
			t.Error(err)
		}
		t.Run("PullRequest", testPullRequest(result))
	}
}

//
// pull request comment sub-tests
//

func testPullRequestComments(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("List", testPullRequestCommentFind(client))
		t.Run("Find", testPullRequestCommentList(client))
	}
}

func testPullRequestCommentFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.PullRequests.FindComment(context.Background(), "gitlab-org/testme", 1, 2990882)
		if err != nil {
			t.Error(err)
		}
		t.Run("Comment", testPullRequestComment(result))
	}
}

func testPullRequestCommentList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.ListOptions{}
		result, _, err := client.PullRequests.ListComments(context.Background(), "gitlab-org/testme", 1, opts)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Errorf("Got empty pull request comment list")
		}
		for _, comment := range result {
			if comment.ID == 2990882 {
				t.Run("Comment", testPullRequestComment(comment))
			}
		}
	}
}

//
// pull request changes sub-tests
//

func testPullRequestChanges(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.ListOptions{}
		result, _, err := client.PullRequests.ListChanges(context.Background(), "gitlab-org/testme", 1, opts)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Errorf("Got empty pull request change list")
			return
		}
		t.Run("File", testChange(result[0]))
	}
}

func testPullRequestCommitList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.ListOptions{}
		result, _, err := client.PullRequests.ListCommits(context.Background(), "gitlab-org/testme", 1, opts)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Errorf("Got empty pull request commit list")
			return
		}
		t.Run("Commit", testPullRequestCommit(result[0]))
	}
}

//
// struct sub-tests
//

func testPullRequest(pr *scm.PullRequest) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := pr.Number, 1; got != want {
			t.Errorf("Want pr Number %d, got %d", want, got)
		}
		if got, want := pr.Title, "JS fix"; got != want {
			t.Errorf("Want pr Title %q, got %q", want, got)
		}
		if got, want := pr.Body, "Signed-off-by: Dmitriy Zaporozhets <dmitriy.zaporozhets@gmail.com>"; got != want {
			t.Errorf("Want pr Body %q, got %q", want, got)
		}
		if got, want := pr.Source, "fix"; got != want {
			t.Errorf("Want pr Source %q, got %q", want, got)
		}
		if got, want := pr.Target, "master"; got != want {
			t.Errorf("Want pr Target %q, got %q", want, got)
		}
		if got, want := pr.Ref, "refs/merge-requests/1/head"; got != want {
			t.Errorf("Want pr Ref %q, got %q", want, got)
		}
		if got, want := pr.Sha, "12d65c8dd2b2676fa3ac47d955accc085a37a9c1"; got != want {
			t.Errorf("Want pr Sha %q, got %q", want, got)
		}
		if got, want := pr.Link, "https://gitlab.com/gitlab-org/testme/-/merge_requests/1"; got != want {
			t.Errorf("Want pr Link %q, got %q", want, got)
		}
		if got, want := pr.Author.Login, "dblessing"; got != want {
			t.Errorf("Want pr Author Login %q, got %q", want, got)
		}
		if got, want := pr.Author.Name, "Drew Blessing"; got != want {
			t.Errorf("Want pr Author Name %q, got %q", want, got)
		}
		if got, want := pr.Author.Avatar, "https://assets.gitlab-static.net/uploads/-/system/user/avatar/13356/avatar.png"; got != want {
			t.Errorf("Want pr Author Avatar %q, got %q", want, got)
		}
		if got, want := pr.Closed, true; got != want {
			t.Errorf("Want pr Closed %v, got %v", want, got)
		}
		if got, want := pr.Merged, false; got != want {
			t.Errorf("Want pr Merged %v, got %v", want, got)
		}
		if got, want := pr.Created.Unix(), int64(1450463393); got != want {
			t.Errorf("Want pr Created %d, got %d", want, got)
		}
		if got, want := pr.Updated.Unix(), int64(1450463422); got != want {
			t.Errorf("Want pr Updated %d, got %d", want, got)
		}
	}
}

func testPullRequestComment(comment *scm.Comment) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := comment.ID, 2990882; got != want {
			t.Errorf("Want pr comment ID %d, got %d", want, got)
		}
		if got, want := comment.Body, "Status changed to closed"; got != want {
			t.Errorf("Want pr comment Body %q, got %q", want, got)
		}
		if got, want := comment.Author.Login, "dblessing"; got != want {
			t.Errorf("Want pr comment Author Login %q, got %q", want, got)
		}
		if got, want := comment.Author.Name, "Drew Blessing"; got != want {
			t.Errorf("Want pr comment Author Name %q, got %q", want, got)
		}
		if got, want := comment.Author.Avatar, "https://assets.gitlab-static.net/uploads/-/system/user/avatar/13356/avatar.png"; got != want {
			t.Errorf("Want pr comment Author Avatar %q, got %q", want, got)
		}
		if got, want := comment.Created.Unix(), int64(1450463422); got != want {
			t.Errorf("Want pr comment Created %d, got %d", want, got)
		}
		if got, want := comment.Updated.Unix(), int64(1450463422); got != want {
			t.Errorf("Want pr comment Updated %d, got %d", want, got)
		}
	}
}

func testChange(change *scm.Change) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := change.Path, "files/js/application.js"; got != want {
			t.Errorf("Want file change Path %q, got %q", want, got)
		}
		if got, want := change.Added, false; got != want {
			t.Errorf("Want file Added %v, got %v", want, got)
		}
		if got, want := change.Deleted, false; got != want {
			t.Errorf("Want file Deleted %v, got %v", want, got)
		}
		if got, want := change.Renamed, false; got != want {
			t.Errorf("Want file Renamed %v, got %v", want, got)
		}
	}
}

func testPullRequestCommit(commit *scm.Commit) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := commit.Message, "JS fix\n\nSigned-off-by: Dmitriy Zaporozhets <dmitriy.zaporozhets@gmail.com>\n"; got != want {
			t.Errorf("Want commit Message %q, got %q", want, got)
		}
		if got, want := commit.Sha, "12d65c8dd2b2676fa3ac47d955accc085a37a9c1"; got != want {
			t.Errorf("Want commit Sha %q, got %q", want, got)
		}
		if got, want := commit.Author.Name, "Dmitriy Zaporozhets"; got != want {
			t.Errorf("Want commit author Name %q, got %q", want, got)
		}
		if got, want := commit.Author.Email, "dmitriy.zaporozhets@gmail.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Author.Date.Unix(), int64(1393489620); got != want {
			t.Errorf("Want commit author Date %d, got %d", want, got)
		}
		if got, want := commit.Committer.Name, "Dmitriy Zaporozhets"; got != want {
			t.Errorf("Want commit author Name %q, got %q", want, got)
		}
		if got, want := commit.Committer.Email, "dmitriy.zaporozhets@gmail.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Committer.Date.Unix(), int64(1393489620); got != want {
			t.Errorf("Want commit author Date %d, got %d", want, got)
		}
	}
}
