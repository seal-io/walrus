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
	}
}

func testPullRequestList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.PullRequestListOptions{
			Open:   true,
			Closed: true,
		}
		result, _, err := client.PullRequests.List(context.Background(), "kit101/drone-yml-test", opts)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Errorf("Got empty pull request list")
		}
		for _, pr := range result {
			if pr.Number == 7 {
				t.Run("PullRequest", testPullRequest(pr))
			}
		}
	}
}

func testPullRequestFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.PullRequests.Find(context.Background(), "kit101/drone-yml-test", 7)
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
		result, _, err := client.PullRequests.FindComment(context.Background(), "kit101/drone-yml-test", 7, 6922557)
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
		result, _, err := client.PullRequests.ListComments(context.Background(), "kit101/drone-yml-test", 7, opts)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Errorf("Got empty pull request comment list")
		}
		for _, comment := range result {
			if comment.ID == 6922557 {
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
		result, _, err := client.PullRequests.ListChanges(context.Background(), "kit101/drone-yml-test", 7, opts)
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

//
// struct sub-tests
//

func testPullRequest(pr *scm.PullRequest) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := pr.Number, 7; got != want {
			t.Errorf("Want pr Number %d, got %d", want, got)
		}
		if got, want := pr.Title, "feat add 3"; got != want {
			t.Errorf("Want pr Title %q, got %q", want, got)
		}
		if got, want := pr.Body, ""; got != want {
			t.Errorf("Want pr Body %q, got %q", want, got)
		}
		if got, want := pr.Source, "feat-2"; got != want {
			t.Errorf("Want pr Source %q, got %q", want, got)
		}
		if got, want := pr.Target, "master"; got != want {
			t.Errorf("Want pr Target %q, got %q", want, got)
		}
		if got, want := pr.Ref, "refs/pull/7/head"; got != want {
			t.Errorf("Want pr Ref %q, got %q", want, got)
		}
		if got, want := pr.Sha, "6168d9dae737b47f00c59fafca10c913a6850c3a"; got != want {
			t.Errorf("Want pr Sha %q, got %q", want, got)
		}
		if got, want := pr.Link, "https://gitee.com/kit101/drone-yml-test/pulls/7"; got != want {
			t.Errorf("Want pr Link %q, got %q", want, got)
		}
		if got, want := pr.Diff, "https://gitee.com/kit101/drone-yml-test/pulls/7.diff"; got != want {
			t.Errorf("Want pr Diff %q, got %q", want, got)
		}
		if got, want := pr.Author.Login, "kit101"; got != want {
			t.Errorf("Want pr Author Login %q, got %q", want, got)
		}
		if got, want := pr.Author.Avatar, "https://portrait.gitee.com/uploads/avatars/user/511/1535738_qkssk1711_1578953939.png"; got != want {
			t.Errorf("Want pr Author Avatar %q, got %q", want, got)
		}
		if got, want := pr.Closed, false; got != want {
			t.Errorf("Want pr Closed %v, got %v", want, got)
		}
		if got, want := pr.Merged, false; got != want {
			t.Errorf("Want pr Merged %v, got %v", want, got)
		}
		if got, want := pr.Created.Unix(), int64(1632996627); got != want {
			t.Errorf("Want pr Created %d, got %d", want, got)
		}
	}
}

func testPullRequestComment(comment *scm.Comment) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := comment.ID, 6922557; got != want {
			t.Errorf("Want pr comment ID %d, got %d", want, got)
		}
		if got, want := comment.Body, "test comment 1"; got != want {
			t.Errorf("Want pr comment Body %q, got %q", want, got)
		}
		if got, want := comment.Author.Login, "kit101"; got != want {
			t.Errorf("Want pr comment Author Login %q, got %q", want, got)
		}
		if got, want := comment.Author.Name, "kit101"; got != want {
			t.Errorf("Want pr comment Author Name %q, got %q", want, got)
		}
		if got, want := comment.Created.Unix(), int64(1633570005); got != want {
			t.Errorf("Want pr comment Created %d, got %d", want, got)
		}
		if got, want := comment.Updated.Unix(), int64(1633570005); got != want {
			t.Errorf("Want pr comment Updated %d, got %d", want, got)
		}
	}
}

func testChange(change *scm.Change) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := change.Path, "change/add3.txt"; got != want {
			t.Errorf("Want file change Path %q, got %q", want, got)
		}
		if got, want := change.Added, true; got != want {
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
