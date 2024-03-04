// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

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

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/pullrequests/4982").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.PullRequests.Find(context.Background(), "atlassian/atlaskit", 4982)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/pullrequests").
		MatchParam("pagelen", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		File("testdata/prs.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.PullRequests.List(context.Background(), "atlassian/atlaskit", scm.PullRequestListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.PullRequest{}
	raw, _ := ioutil.ReadFile("testdata/prs.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/pullrequests/1/diffstat").
		MatchParam("pagelen", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_diffstat.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.PullRequests.ListChanges(context.Background(), "atlassian/atlaskit", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/pr_diffstat.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullMerge(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("2.0/repositories/atlassian/atlaskit/pullrequests/1/merge").
		Reply(200).
		Type("application/json")

	client, _ := New("https://api.bitbucket.org")
	_, err := client.PullRequests.Merge(context.Background(), "atlassian/atlaskit", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestPullClose(t *testing.T) {
	client, _ := New("https://api.bitbucket.org")
	_, err := client.PullRequests.Close(context.Background(), "atlassian/atlaskit", 1)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/atlassian/atlaskit/pullrequests").
		Reply(201).
		Type("application/json").
		File("testdata/pr.json")

	input := &scm.PullRequestInput{
		Title:  "IOS date picker component duplicate March issue",
		Body:   "IOS date picker component duplicate March issue",
		Source: "Lachlan-Vass/ios-date-picker-component-duplicate-marc-1579222909688",
		Target: "master",
	}

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.PullRequests.Create(context.Background(), "atlassian/atlaskit", input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
func TestPullListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/pullrequests/1/commits").
		MatchParam("pagelen", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.PullRequests.ListCommits(context.Background(), "atlassian/atlaskit", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestCommentFind(t *testing.T) {
	_, _, err := NewDefault().PullRequests.FindComment(context.Background(), "", 0, 0)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestPullRequestListComments(t *testing.T) {
	_, _, err := NewDefault().PullRequests.ListComments(context.Background(), "", 0, scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestPullRequestCreateComment(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/atlassian/atlaskit/pullrequests/12").
		Reply(201).
		Type("application/json").
		File("testdata/prcomment.json")

	input := &scm.CommentInput{
		Body: "Lovely comment",
	}

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.PullRequests.CreateComment(context.Background(), "atlassian/atlaskit", 12, input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/prcomment.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestCommentDelete(t *testing.T) {
	_, err := NewDefault().PullRequests.DeleteComment(context.Background(), "", 0, 0)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
