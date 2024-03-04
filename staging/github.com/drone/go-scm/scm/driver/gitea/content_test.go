// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Get("/api/v1/repos/go-gitea/gitea/raw/f05f642b892d59a0a9ef6a31f6c905a24b5db13a/README.md").
		Reply(200).
		Type("plain/text").
		BodyString("Hello World\n")

	client, _ := New("https://try.gitea.io")
	result, _, err := client.Contents.Find(
		context.Background(),
		"go-gitea/gitea",
		"README.md",
		"f05f642b892d59a0a9ef6a31f6c905a24b5db13a",
	)
	if err != nil {
		t.Error(err)
	}

	if got, want := result.Path, "README.md"; got != want {
		t.Errorf("Want file Path %q, got %q", want, got)
	}
	if got, want := string(result.Data), "Hello World\n"; got != want {
		t.Errorf("Want file Data %q, got %q", want, got)
	}
}

func TestContentCreate(t *testing.T) {
	client, _ := New("https://try.gitea.io")
	_, err := client.Contents.Create(context.Background(), "go-gitea/gitea", "README.md", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentUpdate(t *testing.T) {
	client, _ := New("https://try.gitea.io")
	_, err := client.Contents.Update(context.Background(), "go-gitea/gitea", "README.md", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentDelete(t *testing.T) {
	client, _ := New("https://try.gitea.io")
	_, err := client.Contents.Delete(context.Background(), "go-gitea/gitea", "README.md", &scm.ContentParams{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Get("/api/v1/repos/go-gitea/gitea/contents/docs/content/doc").
		MatchParam("ref", "master").
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client, _ := New("https://try.gitea.io")
	got, _, err := client.Contents.List(
		context.Background(),
		"go-gitea/gitea",
		"docs/content/doc",
		"master",
		scm.ListOptions{},
	)
	if err != nil {
		t.Error(err)
	}

	want := []*scm.ContentInfo{}
	raw, _ := ioutil.ReadFile("testdata/content_list.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
