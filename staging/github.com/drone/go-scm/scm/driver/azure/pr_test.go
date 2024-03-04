// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(201).
		Type("application/json").
		File("testdata/pr.json")

	input := scm.PullRequestInput{
		Title:  "test_pr",
		Body:   "test_pr_body",
		Source: "pr_branch",
		Target: "main",
	}

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.PullRequests.Create(context.Background(), "REPOID", &input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullFind(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/pullrequests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.PullRequests.Find(context.Background(), "REPOID", 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/pullRequests/1/commits").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.PullRequests.ListCommits(context.Background(), "REPOID", 1, scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
