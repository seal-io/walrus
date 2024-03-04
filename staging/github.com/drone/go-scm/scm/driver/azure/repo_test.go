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

func TestRepositoryList(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories").
		Reply(200).
		Type("application/json").
		File("testdata/repos.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Repositories.List(context.Background(), scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	jsonErr := json.Unmarshal(raw, &want)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryHookCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/_apis/projects").
		Reply(201).
		Type("application/json").
		File("testdata/projects.json")

	gock.New("https:/dev.azure.com/").
		Post("/ORG/_apis/hooks/subscriptions").
		Reply(201).
		Type("application/json").
		File("testdata/hook.json")

	in := &scm.HookInput{
		Name:         "web",
		NativeEvents: []string{"git.push"},
		Target:       "http://www.example.com/webhook",
	}

	client := NewDefault("ORG", "test_project")
	got, _, err := client.Repositories.CreateHook(context.Background(), "test_project", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestHooksList(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/_apis/projects").
		Reply(201).
		Type("application/json").
		File("testdata/projects.json")

	gock.New("https:/dev.azure.com/").
		Get("/ORG/_apis/hooks/subscriptions").
		Reply(200).
		Type("application/json").
		File("testdata/hooks.json")

	client := NewDefault("ORG", "test_project")
	repoID := "fde2d21f-13b9-4864-a995-83329045289a"

	got, _, err := client.Repositories.ListHooks(context.Background(), repoID, scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Hook{}
	raw, _ := ioutil.ReadFile("testdata/hooks.json.golden")
	jsonErr := json.Unmarshal(raw, &want)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryHookDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Delete("/ORG/_apis/hooks/subscriptions").
		Reply(204).
		Type("application/json")

	client := NewDefault("ORG", "PROJ")
	res, err := client.Repositories.DeleteHook(context.Background(), "", "test-project")
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 204; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}

}

func TestRepositoryFind(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/test_project").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Repositories.Find(context.Background(), "test_project")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/repo.json.golden")
	jsonErr := json.Unmarshal(raw, &want)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}
