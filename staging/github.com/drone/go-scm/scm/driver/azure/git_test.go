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

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/commit.json")

	client := NewDefault("ORG", "PROJ")

	got, _, err := client.Git.FindCommit(context.Background(), "REPOID", "14897f4465d2d63508242b5cbf68aa2865f693e7")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Commit)
	raw, _ := ioutil.ReadFile("testdata/commit.json.golden")

	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitCreateBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(201).
		Type("application/json").
		File("testdata/branch_create.json")

	params := &scm.ReferenceInput{
		Name: "test_branch",
		Sha:  "312797ba52425353dec56871a255e2a36fc96344",
	}

	client := NewDefault("ORG", "PROJ")
	res, err := client.Git.CreateBranch(context.Background(), "REPOID", params)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 201 {
		t.Errorf("Unexpected Results")
	}
}

func TestGitListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Git.ListCommits(context.Background(), "REPOID", scm.CommitListOptions{})
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

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/branches.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Git.ListBranches(context.Background(), "REPOID", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/branches.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitListBranchesV2(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/branches_filter.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Git.ListBranchesV2(context.Background(), "REPOID", scm.BranchListOptions{SearchTerm: "main"})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/branches_filter.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitCompareChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/compare.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Git.CompareChanges(context.Background(), "REPOID", "9788e5ddf8b387cb79228628f34d8dc18582d606", "66df312dad61e84dd896d1e8d14ee3dce53b62f0", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/compare.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}
