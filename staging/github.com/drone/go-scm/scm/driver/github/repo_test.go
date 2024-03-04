// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestRepositoryFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/repo.json")

	client := NewDefault()
	got, res, err := client.Repositories.Find(context.Background(), "octocat/hello-world")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/repo.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryPerms(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/repo.json")

	client := NewDefault()
	got, res, err := client.Repositories.FindPerms(context.Background(), "octocat/hello-world")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/repo.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want.Perm); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryNotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
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
	if got, want := err.Error(), "Not Found"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
}

func TestRepositoryList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/user/repos").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/repos.json")

	client := NewDefault()
	got, res, err := client.Repositories.List(context.Background(), scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestRepositoryListV2(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/search/repositories").
		MatchParam("q", "testRepoin:name").
		MatchParam("q", "user:user123").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/repos_filter.json")

	client := NewDefault()
	got, res, err := client.Repositories.ListV2(context.Background(), scm.RepoListOptions{
		ListOptions: scm.ListOptions{Page: 1, Size: 30},
		RepoSearchTerm: scm.RepoSearchTerm{
			RepoName: "testRepo",
			User:     "user123",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos_filter.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGithubAppInstallationList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/installation/repositories").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/github_app_repos.json")

	client := NewDefault()
	got, res, err := client.Repositories.(*RepositoryService).ListByInstallation(context.Background(), scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestStatusList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/statuses.json")

	client := NewDefault()
	got, res, err := client.Repositories.ListStatus(context.Background(), "octocat/hello-world", "6dcb09b5b57875f334f61aebed695e2e4193db5e", scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Status{}
	raw, _ := ioutil.ReadFile("testdata/statuses.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestStatusCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/status.json")

	in := &scm.StatusInput{
		Desc:   "Build has completed successfully",
		Label:  "continuous-integration/drone",
		State:  scm.StateSuccess,
		Target: "https://ci.example.com/1000/output",
	}

	client := NewDefault()
	got, res, err := client.Repositories.CreateStatus(context.Background(), "octocat/hello-world", "6dcb09b5b57875f334f61aebed695e2e4193db5e", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Status)
	raw, _ := ioutil.ReadFile("testdata/status.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestDeployStatusCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/deployments/1/status").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/deployment.json")

	in := &scm.DeployStatus{
		Number:         1,
		Desc:           "Build has completed successfully",
		State:          scm.StateSuccess,
		Target:         "https://ci.example.com/1000/output",
		Environment:    "production",
		EnvironmentURL: "https://example.netlify.com",
	}

	client := NewDefault()
	got, res, err := client.Repositories.(*RepositoryService).CreateDeployStatus(context.Background(), "octocat/hello-world", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.DeployStatus)
	raw, _ := ioutil.ReadFile("testdata/deployment.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/hooks/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	client := NewDefault()
	got, res, err := client.Repositories.FindHook(context.Background(), "octocat/hello-world", "1")
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

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/hooks").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/hooks.json")

	client := NewDefault()
	got, res, err := client.Repositories.ListHooks(context.Background(), "octocat/hello-world", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Hook{}
	raw, _ := ioutil.ReadFile("testdata/hooks.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestRepositoryHookCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/hooks").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	in := &scm.HookInput{
		Name:       "drone",
		Target:     "https://example.com",
		Secret:     "topsecret",
		SkipVerify: true,
	}

	client := NewDefault()
	got, res, err := client.Repositories.CreateHook(context.Background(), "octocat/hello-world", in)
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

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Patch("/repos/octocat/hello-world/hooks/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	in := &scm.HookInput{
		Name:       "drone",
		Target:     "https://example.com",
		Secret:     "topsecret",
		SkipVerify: true,
	}

	client := NewDefault()
	got, res, err := client.Repositories.UpdateHook(context.Background(), "octocat/hello-world", "1", in)
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

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Delete("/repos/octocat/hello-world/hooks/1").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Repositories.DeleteHook(context.Background(), "octocat/hello-world", "1")
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 204; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestConvertState(t *testing.T) {
	tests := []struct {
		src string
		dst scm.State
	}{
		{
			src: "failure",
			dst: scm.StateFailure,
		},
		{
			src: "error",
			dst: scm.StateError,
		},
		{
			src: "pending",
			dst: scm.StatePending,
		},
		{
			src: "success",
			dst: scm.StateSuccess,
		},
		{
			src: "invalid",
			dst: scm.StateUnknown,
		},
	}
	for _, test := range tests {
		if got, want := convertState(test.src), test.dst; got != want {
			t.Errorf("Want state %s converted to %v", test.src, test.dst)
		}
	}
}

func TestConvertFromState(t *testing.T) {
	tests := []struct {
		src scm.State
		dst string
	}{
		{
			src: scm.StateCanceled,
			dst: "error",
		},
		{
			src: scm.StateError,
			dst: "error",
		},
		{
			src: scm.StateFailure,
			dst: "failure",
		},
		{
			src: scm.StatePending,
			dst: "pending",
		},
		{
			src: scm.StateRunning,
			dst: "pending",
		},
		{
			src: scm.StateSuccess,
			dst: "success",
		},
		{
			src: scm.StateUnknown,
			dst: "error",
		},
	}
	for _, test := range tests {
		if got, want := convertFromState(test.src), test.dst; got != want {
			t.Errorf("Want state %v converted to %s", test.src, test.dst)
		}
	}
}

func TestHookEvents(t *testing.T) {
	tests := []struct {
		in  scm.HookEvents
		out []string
	}{
		{
			in:  scm.HookEvents{Push: true},
			out: []string{"push"},
		},
		{
			in:  scm.HookEvents{Branch: true},
			out: []string{"create", "delete"},
		},
		{
			in:  scm.HookEvents{IssueComment: true},
			out: []string{"issue_comment"},
		},
		{
			in:  scm.HookEvents{PullRequestComment: true},
			out: []string{"pull_request_review_comment", "issue_comment"},
		},
		{
			in:  scm.HookEvents{Issue: true},
			out: []string{"issues"},
		},
		{
			in:  scm.HookEvents{PullRequest: true},
			out: []string{"pull_request"},
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
			out: []string{"push", "pull_request", "pull_request_review_comment", "issues", "issue_comment", "create", "delete", "deployment"},
		},
	}
	for i, test := range tests {
		got, want := convertFromHookEvents(test.in), test.out
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Unexpected Results at index %d", i)
			t.Log(diff)
		}
	}
}

func Test_convertRepository(t *testing.T) {
	permissionsStruct := &repository{
		Name: "bla",
	}
	permissionsStruct.Permissions.Admin = true
	tests := []struct {
		name string
		from *repository
		want *scm.Repository
	}{
		{
			name: "Simple",
			from: &repository{
				Name:        "hello-world",
				ID:          1,
				Permissions: permissionsStruct.Permissions,
			},
			want: &scm.Repository{
				Name: "hello-world",
				ID:   "1",
				Perm: &scm.Perm{
					Admin: true,
				},
			},
		},
		{
			name: "null",
			from: nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertRepository(tt.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepositoryListWithNull(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/user/repos").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/reposWithNull.json")

	client := NewDefault()
	got, res, err := client.Repositories.List(context.Background(), scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}
