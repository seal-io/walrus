// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestRepositoryFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/repo.json")

	client := NewDefault()
	got, res, err := client.Repositories.Find(context.Background(), "kit101/drone-yml-test")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/repo.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestRepositoryHookFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/hooks/787341").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	client := NewDefault()
	got, res, err := client.Repositories.FindHook(context.Background(), "kit101/drone-yml-test", "787341")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestRepositoryPerms(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/repo.json")

	client := NewDefault()
	got, res, err := client.Repositories.FindPerms(context.Background(), "kit101/drone-yml-test")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/repo.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want.Perm); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestRepositoryList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/user/repos").
		MatchParam("page", "1").
		MatchParam("per_page", "3").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/repos.json")

	client := NewDefault()
	got, res, err := client.Repositories.List(context.Background(), scm.ListOptions{Page: 1, Size: 3})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Page", testPage(res))
}

func TestRepositoryListHook(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/hooks").
		MatchParam("page", "1").
		MatchParam("per_page", "3").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/hooks.json")

	client := NewDefault()
	got, res, err := client.Repositories.ListHooks(context.Background(), "kit101/drone-yml-test", scm.ListOptions{Page: 1, Size: 3})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Hook{}
	raw, _ := ioutil.ReadFile("testdata/hooks.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Page", testPage(res))
}

func TestRepositoryListStatus(t *testing.T) {
	_, _, err := NewDefault().Repositories.ListStatus(context.Background(), "", "", scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestRepositoryCreateHook(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/drone-yml-test/hooks").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	in := &scm.HookInput{
		Target: "http://test.kit101.com/webhook",
		Secret: "123asdas123",
	}

	client := NewDefault()
	got, res, err := client.Repositories.CreateHook(context.Background(), "kit101/drone-yml-test", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestRepositoryCreateStatus(t *testing.T) {
	_, _, err := NewDefault().Repositories.CreateStatus(context.Background(), "", "", &scm.StatusInput{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestRepositoryUpdateHook(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Patch("/repos/kit101/drone-yml-test/hooks/787341").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	in := &scm.HookInput{
		Target: "http://test.kit101.com/webhook",
		Secret: "123asdas123",
	}

	client := NewDefault()
	got, res, err := client.Repositories.UpdateHook(context.Background(), "kit101/drone-yml-test", "787341", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestRepositoryDeleteHook(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Delete("/repos/kit101/drone-yml-test/hooks/787341").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Repositories.DeleteHook(context.Background(), "kit101/drone-yml-test", "787341")
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 204; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
}

func TestRepositoryNotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/dev/null").
		Reply(404).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/error.json")

	client := NewDefault()
	_, _, err := client.Repositories.Find(context.Background(), "dev/null")
	if err == nil {
		t.Errorf("Expect Not Found error")
		return
	}
	if got, want := err.Error(), "404 Project Not Found"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
}

func TestHookEvents(t *testing.T) {
	tests := []struct {
		in  scm.HookEvents
		out *hook
	}{
		{
			in:  scm.HookEvents{Push: true},
			out: &hook{PushEvents: true},
		},
		{
			in:  scm.HookEvents{Branch: true},
			out: &hook{},
		},
		{
			in:  scm.HookEvents{IssueComment: true},
			out: &hook{NoteEvents: true},
		},
		{
			in:  scm.HookEvents{PullRequestComment: true},
			out: &hook{NoteEvents: true},
		},
		{
			in:  scm.HookEvents{Issue: true},
			out: &hook{IssuesEvents: true},
		},
		{
			in:  scm.HookEvents{PullRequest: true},
			out: &hook{MergeRequestsEvents: true},
		},
		{
			in: scm.HookEvents{
				Branch:             true,
				Deployment:         true,
				Issue:              true,
				IssueComment:       true,
				PullRequest:        true,
				PullRequestComment: true,
				Push:               true,
				ReviewComment:      true,
				Tag:                true,
			},
			out: &hook{
				IssuesEvents:        true,
				MergeRequestsEvents: true,
				NoteEvents:          true,
				PushEvents:          true,
				TagPushEvents:       true,
			},
		},
	}

	for i, test := range tests {
		fmt.Println(test, i)
		got := new(hook)
		convertFromHookEvents(test.in, got)
		want := test.out
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Unexpected Results at index %d", i)
			t.Log(diff)
		}
	}
}
