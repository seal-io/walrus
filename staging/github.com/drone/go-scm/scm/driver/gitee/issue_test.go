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

func TestEncodeAndDecodeIssueNumber(t *testing.T) {
	giteeIssueNumber := "I4CD5P"
	scmIssueNumber := 735267685380
	encodedScmIssueNumber := encodeNumber(giteeIssueNumber)
	if diff := cmp.Diff(encodedScmIssueNumber, scmIssueNumber); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
	encodedGiteeIssueNumber := decodeNumber(scmIssueNumber)
	if diff := cmp.Diff(encodedGiteeIssueNumber, giteeIssueNumber); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestIssueFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/issues/I4CD5P").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue.json")

	client := NewDefault()
	got, res, err := client.Issues.Find(context.Background(), "kit101/drone-yml-test", 735267685380)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Issue)
	raw, _ := ioutil.ReadFile("testdata/issue.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestIssueFindComment(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/issues/comments/6877445").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue_comment.json")

	client := NewDefault()
	got, res, err := client.Issues.FindComment(context.Background(), "kit101/drone-yml-test", 735267685380, 6877445)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/issue_comment.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestIssueList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/issues").
		MatchParam("page", "1").
		MatchParam("per_page", "3").
		MatchParam("state", "all").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/issues.json")

	client := NewDefault()
	got, res, err := client.Issues.List(context.Background(), "kit101/drone-yml-test", scm.IssueListOptions{Page: 1, Size: 3, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Issue{}
	raw, _ := ioutil.ReadFile("testdata/issues.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Page", testPage(res))
}

func TestIssueListComments(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml/issues/I4CD5P/comments").
		MatchParam("page", "1").
		MatchParam("per_page", "3").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/issue_comments.json")

	client := NewDefault()
	got, res, err := client.Issues.ListComments(context.Background(), "kit101/drone-yml", 735267685380, scm.ListOptions{Size: 3, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Comment{}
	raw, _ := ioutil.ReadFile("testdata/issue_comments.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Page", testPage(res))
}

func TestIssueCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/issues").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue.json")

	input := scm.IssueInput{
		Title: "test issue 1",
		Body:  "test issue 1",
	}

	client := NewDefault()
	got, res, err := client.Issues.Create(context.Background(), "kit101/drone-yml-test", &input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Issue)
	raw, _ := ioutil.ReadFile("testdata/issue.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestIssueCreateComment(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Post("/repos/kit101/drone-yml-test/issues/I4CD5P/comments").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue_comment.json")

	input := &scm.CommentInput{
		Body: "it's ok.",
	}

	client := NewDefault()
	got, res, err := client.Issues.CreateComment(context.Background(), "kit101/drone-yml-test", 735267685380, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/issue_comment.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestIssueDeleteComment(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Delete("/repos/kit101/drone-yml-test/issues/comments/6879139").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Issues.DeleteComment(context.Background(), "kit101/drone-yml-test", 735267685380, 6879139)
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 204; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
}

func TestIssueClose(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Patch("/repos/kit101/issues/I4CD5P").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue.json")

	client := NewDefault()
	res, err := client.Issues.Close(context.Background(), "kit101/drone-yml-test", 735267685380)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
}

func TestIssueLock(t *testing.T) {
	_, err := NewDefault().Issues.Lock(context.Background(), "kit101/drone-yml-test", 735267685380)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestIssueUnlock(t *testing.T) {
	_, err := NewDefault().Issues.Unlock(context.Background(), "kit101/drone-yml-test", 735267685380)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
