// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/h2non/gock"
)

func TestListCommits(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/commits").
			Reply(200).
			Type("application/json").
			File("testdata/commits.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Git.ListCommits(context.Background(), harnessRepo, scm.CommitListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	wantErr := json.Unmarshal(raw, &want)
	if wantErr != nil {
		t.Error(wantErr)
		return
	}
	if harnessPAT != "" && len(got) > 0 {
		// if testing against a real system and we get commits
		return
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestFindCommit(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/commits/1d640265d8bdd818175fa736f0fcbad2c9b716c9").
			Reply(200).
			Type("application/json").
			File("testdata/commit.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Git.FindCommit(context.Background(), harnessRepo, "1d640265d8bdd818175fa736f0fcbad2c9b716c9")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Commit)
	raw, _ := ioutil.ReadFile("testdata/commit.json.golden")
	wantErr := json.Unmarshal(raw, &want)
	if wantErr != nil {
		t.Error(wantErr)
		return
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestFindBranch(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/branches/main").
			Reply(200).
			Type("application/json").
			File("testdata/branch.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Git.FindBranch(context.Background(), harnessRepo, "main")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/branch.json.golden")
	wantErr := json.Unmarshal(raw, &want)
	if wantErr != nil {
		t.Error(wantErr)
		return
	}

	if diff := cmp.Diff(got, want, cmpopts.IgnoreFields(scm.Reference{}, "Sha")); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestListBranches(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/branches").
			Reply(200).
			Type("application/json").
			File("testdata/branches.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Git.ListBranches(context.Background(), harnessRepo, scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/branches.json.golden")
	wantErr := json.Unmarshal(raw, &want)
	if wantErr != nil {
		t.Error(wantErr)
		return
	}

	if diff := cmp.Diff(got, want, cmpopts.IgnoreFields(scm.Reference{}, "Sha")); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestCreateBranch(t *testing.T) {

	defer gock.Off()

	gock.New(gockOrigin).
		Post("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/branches").
		Reply(200).
		Type("application/json").
		File("testdata/branch.json")

	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	input := &scm.ReferenceInput{
		Name: "test",
		Sha:  "e8ef0374ca0cee8048e94b28eaf0d9e2e2515a14",
	}
	result, err := client.Git.CreateBranch(context.Background(), harnessRepo, input)
	if err != nil {
		t.Error(err)
		return
	}

	if result.Status != 200 {
		t.Errorf("Unexpected Results")
	}

}

func TestCompareChanges(t *testing.T) {
	source := "542ddabd47d7bfa79359b7b4e2af7f975354e35f"
	target := "c7d0d4b21d5cfdf47475ff1f6281ef1a91883d"
	defer gock.Off()

	gock.New(gockOrigin).
		Get(fmt.Sprintf("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/diff/%s...%s", source, target)).
		Reply(200).
		Type("application/json").
		File("testdata/gitdiff.json")

	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, result, err := client.Git.CompareChanges(context.Background(), harnessRepo, source, target, scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if result.Status != 200 {
		t.Errorf("Unexpected Results")
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/gitdiff.json.golden")
	wantErr := json.Unmarshal(raw, &want)
	if wantErr != nil {
		t.Error(wantErr)
		return
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
