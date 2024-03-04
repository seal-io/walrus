package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestReleaseFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/releases/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/release.json")

	client := NewDefault()
	got, res, err := client.Releases.Find(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Release)
	raw, _ := ioutil.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReleaseList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/releases").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("state", "all").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/releases.json")

	client := NewDefault()
	got, res, err := client.Releases.List(context.Background(), "octocat/hello-world", scm.ReleaseListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Release{}
	raw, _ := ioutil.ReadFile("testdata/releases.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReleaseCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/releases").
		File("testdata/release_create.json").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/release.json")

	client := NewDefault()
	input := &scm.ReleaseInput{
		Title:       "v1.0",
		Description: "Tracking release for version 1.0",
		Tag:         "v1.0",
	}

	got, res, err := client.Releases.Create(context.Background(), "octocat/hello-world", input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Release)
	raw, _ := ioutil.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReleaseUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Patch("/repos/octocat/hello-world/releases/1").
		File("testdata/release_create.json").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/release.json")

	client := NewDefault()
	input := &scm.ReleaseInput{
		Title:       "v1.0",
		Description: "Tracking release for version 1.0",
		Tag:         "v1.0",
	}

	got, res, err := client.Releases.Update(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Release)
	raw, _ := ioutil.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReleaseDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Delete("/repos/octocat/hello-world/releases/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Releases.Delete(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
