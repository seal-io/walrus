// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/pulls/6").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr.json")

	client := NewDefault()
	got, res, err := client.PullRequests.Find(context.Background(), "kit101/drone-yml-test", 6)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestPullFindComment(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/pulls/comments/6922557").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_comment.json")

	client := NewDefault()
	got, res, err := client.PullRequests.FindComment(context.Background(), "kit101/drone-yml-test", 7, 6922557)
	if err != nil {
		t.Error(err)
		return
	}
	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestPullList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/pulls").
		MatchParam("page", "1").
		MatchParam("per_page", "3").
		MatchParam("state", "all").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pulls.json")

	client := NewDefault()
	got, res, err := client.PullRequests.List(context.Background(), "kit101/drone-yml-test", scm.PullRequestListOptions{Page: 1, Size: 3, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.PullRequest{}
	raw, _ := ioutil.ReadFile("testdata/pulls.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Page", testPage(res))
}

func TestPullListChanges(t *testing.T) {
	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/pulls/6/files").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pr_files.json")

	client := NewDefault()
	got, res, err := client.PullRequests.ListChanges(context.Background(), "kit101/drone-yml-test", 6, scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/pr_files.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestPullListComments(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/pulls/7/comments").
		MatchParam("page", "1").
		MatchParam("per_page", "3").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pr_comments.json")

	client := NewDefault()
	got, res, err := client.PullRequests.ListComments(context.Background(), "kit101/drone-yml-test", 7, scm.ListOptions{Page: 1, Size: 3})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Comment{}
	raw, _ := ioutil.ReadFile("testdata/pr_comments.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Page", testPage(res))
}

func TestPullListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/pulls/7/commits").
		Reply(200).
		Type("application/json").
		File("testdata/pr_commits.json")

	client := NewDefault()
	got, _, err := client.PullRequests.ListCommits(context.Background(), "kit101/drone-yml-test", 7, scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/pr_commits.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullMerge(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Put("/repos/kit101/drone-yml-test/pulls/6/merge").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.PullRequests.Merge(context.Background(), "kit101/drone-yml-test", 6)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
}

func TestPullClose(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Patch("/repos/kit101/drone-yml-test/pulls/6").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.PullRequests.Close(context.Background(), "kit101/drone-yml-test", 6)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/drone-yml-test/pulls").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr.json")

	input := scm.PullRequestInput{
		Title:  "new-feature",
		Body:   "Please pull these awesome changes",
		Source: "crud",
		Target: "master",
	}

	client := NewDefault()
	got, res, err := client.PullRequests.Create(context.Background(), "kit101/drone-yml-test", &input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestPullCommentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/drone-yml-test/pulls/7/comments").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_comment.json")

	client := NewDefault()
	input := scm.CommentInput{
		Body: "test comment 1",
	}
	got, res, err := client.PullRequests.CreateComment(context.Background(), "kit101/drone-yml-test", 7, &input)
	if err != nil {
		t.Error(err)
		return
	}
	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestPullCommentDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Delete("/repos/kit101/drone-yml-test/pulls/comments/6990588").
		Reply(204).
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.PullRequests.DeleteComment(context.Background(), "kit101/drone-yml-test", 7, 6990588)
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 204; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
}
