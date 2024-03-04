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

func TestGitCreateBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/drone-yml-test/branches").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	input := scm.ReferenceInput{
		Name: "create-by-api",
		Sha:  "b72a4c4a2d838d96a545a42d41d7776ae5566f4a",
	}
	res, err := client.Git.CreateBranch(context.Background(), "kit101/drone-yml-test", &input)
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 201; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
}

func TestGitFindBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/branches/master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/branch.json")

	client := NewDefault()
	got, res, err := client.Git.FindBranch(context.Background(), "kit101/drone-yml-test", "master")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/branch.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/commits/e3c0ff4d5cef439ea11b30866fb1ed79b420801d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/commit.json")

	client := NewDefault()
	got, res, err := client.Git.FindCommit(context.Background(), "kit101/drone-yml-test", "e3c0ff4d5cef439ea11b30866fb1ed79b420801d")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Commit)
	raw, _ := ioutil.ReadFile("testdata/commit.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestGitFindTag(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/tags").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/tags.json")

	client := NewDefault()
	got, res, err := client.Git.FindTag(context.Background(), "kit101/drone-yml-test", "1.0")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/tag.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/branches").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/branches.json")

	client := NewDefault()
	got, res, err := client.Git.ListBranches(context.Background(), "kit101/drone-yml-test", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/branches.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestGitListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/commits").
		MatchParam("sha", "master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/commits.json")

	client := NewDefault()
	got, res, err := client.Git.ListCommits(context.Background(), "kit101/drone-yml-test",
		scm.CommitListOptions{Ref: "master", Page: 1, Size: 3})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Page", testPage(res))
}

func TestGitListTags(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/tags").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/tags.json")

	client := NewDefault()
	got, res, err := client.Git.ListTags(context.Background(), "kit101/drone-yml-test", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/tags.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestGitListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/repos/kit101/drone-yml-test/commits/7e84b6f94b8d4bfaa051910cc4ce16b73bcffd51").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/changes.json")

	client := NewDefault()
	got, res, err := client.Git.ListChanges(context.Background(), "kit101/drone-yml-test", "7e84b6f94b8d4bfaa051910cc4ce16b73bcffd51", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/changes.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestGitCompareChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/compare/e3c0ff4d5cef439ea11b30866fb1ed79b420801d...2700445cd84c08546f4d003f8aa54d2099a006b7").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/compare.json")

	client := NewDefault()
	got, res, err := client.Git.CompareChanges(context.Background(), "kit101/drone-yml-test", "e3c0ff4d5cef439ea11b30866fb1ed79b420801d", "2700445cd84c08546f4d003f8aa54d2099a006b7", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/compare.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}
