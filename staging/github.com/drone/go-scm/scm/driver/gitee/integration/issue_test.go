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
// issue sub-tests
//

func testIssues(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("List", testIssueList(client))
		t.Run("Find", testIssueFind(client))
		t.Run("Comments", testIssueComments(client))
	}
}

func testIssueList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.IssueListOptions{
			Open:   true,
			Closed: true,
		}
		result, _, err := client.Issues.List(context.Background(), "kit101/drone-yml-test", opts)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Errorf("Got empty issue list")
		}
	}
}

func testIssueFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Issues.Find(context.Background(), "kit101/drone-yml-test", 735267685380)
		if err != nil {
			t.Error(err)
		}
		t.Run("Issue", testIssue(result))
	}
}

//
// issue comment sub-tests
//

func testIssueComments(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("List", testIssueCommentList(client))
		t.Run("Find", testIssueCommentFind(client))
	}
}

func testIssueCommentList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		opts := scm.ListOptions{}
		result, _, err := client.Issues.ListComments(context.Background(), "kit101/drone-yml-test", 735267685380, opts)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty issue comment list")
		}
		for _, comment := range result {
			if comment.ID == 6877445 {
				t.Run("Comment", testIssueComment(comment))
			}
		}
	}
}

func testIssueCommentFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Issues.FindComment(context.Background(), "kit101/drone-yml-test", 735267685380, 6877445)
		if err != nil {
			t.Error(err)
		}
		t.Run("Comment", testIssueComment(result))
	}
}

//
// struct sub-tests
//

func testIssue(issue *scm.Issue) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := issue.Number, 735267685380; got != want {
			t.Errorf("Want issue Number %d, got %d", want, got)
		}
		if got, want := issue.Title, "test issue 1"; got != want {
			t.Errorf("Want issue Title %q, got %q", want, got)
		}
		if got, want := issue.Body, "test issue 1"; got != want {
			t.Errorf("Want issue Body %q, got %q", want, got)
		}
		if got, want := issue.Link, "https://gitee.com/kit101/drone-yml-test/issues/I4CD5P"; got != want {
			t.Errorf("Want issue Link %q, got %q", want, got)
		}
		if got, want := issue.Author.Login, "kit101"; got != want {
			t.Errorf("Want issue Author Login %q, got %q", want, got)
		}
		if got, want := issue.Author.Avatar, "https://portrait.gitee.com/uploads/avatars/user/511/1535738_qkssk1711_1578953939.png"; got != want {
			t.Errorf("Want issue Author Avatar %q, got %q", want, got)
		}
		if got, want := issue.Closed, true; got != want {
			t.Errorf("Want issue Closed %v, got %v", want, got)
		}
		if got, want := issue.Created.Unix(), int64(1632878488); got != want {
			t.Errorf("Want issue Created %d, got %d", want, got)
		}
	}
}

func testIssueComment(comment *scm.Comment) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := comment.ID, 6877445; got != want {
			t.Errorf("Want issue comment ID %d, got %d", want, got)
		}
		if got, want := comment.Body, "it's ok."; got != want {
			t.Errorf("Want issue comment Body %q, got %q", want, got)
		}
		if got, want := comment.Author.Login, "kit101"; got != want {
			t.Errorf("Want issue comment Author Login %q, got %q", want, got)
		}
		if got, want := comment.Created.Unix(), int64(1632884450); got != want {
			t.Errorf("Want issue comment Created %d, got %d", want, got)
		}
		if got, want := comment.Updated.Unix(), int64(1632884450); got != want {
			t.Errorf("Want issue comment Updated %d, got %d", want, got)
		}
	}
}
