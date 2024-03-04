// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/commits/131cb13f4aed12e725177bc4b7c28db67839bf9f").
		Reply(200).
		Type("application/json").
		File("testdata/commit.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.FindCommit(context.Background(), "PRJ/my-repo", "131cb13f4aed12e725177bc4b7c28db67839bf9f")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Commit)
	raw, _ := ioutil.ReadFile("testdata/commit.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitFindBranch(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/branches").
		MatchParam("filterText", "master").
		Reply(200).
		Type("application/json").
		File("testdata/branch.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.FindBranch(context.Background(), "PRJ/my-repo", "master")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/branch.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitFindTag(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/tags").
		MatchParam("filterText", "v1.0.0").
		Reply(200).
		Type("application/json").
		File("testdata/tag.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.FindTag(context.Background(), "PRJ/my-repo", "v1.0.0")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/tag.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/commits").
		MatchParam("until", "").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.ListCommits(context.Background(), "PRJ/my-repo", scm.CommitListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/branches").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/branches.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.ListBranches(context.Background(), "PRJ/my-repo", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/branches.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
	//
	// t.Run("Page", testPage(res))
}

func TestGitListBranchesWithBranchFilter(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/branches").
		MatchParam("filterText", "mast").
		Reply(200).
		Type("application/json").
		File("testdata/branches_filter.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.ListBranchesV2(context.Background(), "PRJ/my-repo", scm.BranchListOptions{SearchTerm: "mast"})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/branches_filter.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
	//
	// t.Run("Page", testPage(res))
}

func TestGitListTags(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/tags").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/tags.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.ListTags(context.Background(), "PRJ/my-repo", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/tags.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	// t.Run("Page", testPage(res))
}

func TestGitListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/commits/131cb13f4aed12e725177bc4b7c28db67839bf9f/changes").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/changes.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.ListChanges(context.Background(), "PRJ/my-repo", "131cb13f4aed12e725177bc4b7c28db67839bf9f", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/changes.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitCompareChanges(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/compare/changes").
		MatchParam("from", "4f4b0ef1714a5b6cafdaf2f53c7f5f5b38fb9348").
		MatchParam("to", "131cb13f4aed12e725177bc4b7c28db67839bf9f").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/compare.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Git.CompareChanges(context.Background(), "PRJ/my-repo", "4f4b0ef1714a5b6cafdaf2f53c7f5f5b38fb9348", "131cb13f4aed12e725177bc4b7c28db67839bf9f", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/compare.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestCreateBranch(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("/rest/api/1.0/projects/PRJ/repos/my-repo/branches").
		Reply(201).
		Type("application/json").
		File("testdata/branch_create.json")

	client, _ := New("http://example.com:7990")
	params := &scm.ReferenceInput{
		Name: "Hello",
		Sha:  "312797ba52425353dec56871a255e2a36fc96344",
	}
	res, err := client.Git.CreateBranch(context.Background(), "PRJ/my-repo", params)

	if err != nil {
		t.Errorf("Encountered err while creating branch " + err.Error())
	}

	if res.Status != 201 {
		t.Errorf("The error response of branch creation is not 201")
	}
}
