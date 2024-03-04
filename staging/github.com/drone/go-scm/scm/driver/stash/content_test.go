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

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/raw/README").
		MatchParam("at", "5c64a07cd6c0f21b753bf261ef059c7e7633c50a").
		Reply(200).
		Type("text/plain").
		File("testdata/content.txt")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Contents.Find(context.Background(), "PRJ/my-repo", "README", "5c64a07cd6c0f21b753bf261ef059c7e7633c50a")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/raw/README").
		MatchParam("at", "b1&b2").
		Reply(200).
		Type("text/plain").
		File("testdata/content.txt")

	client, _ = New("http://example.com:7990")
	got, _, err = client.Contents.Find(context.Background(), "PRJ/my-repo", "README", "b1&b2")
	if err != nil {
		t.Error(err)
	}

	want = new(scm.Content)
	raw, _ = ioutil.ReadFile("testdata/content.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost:7990").
		Put("/rest/api/1.0/projects/octocat/repos/hello-world/browse/README").
		Reply(200).
		Type("application/json").
		File("testdata/content_create.json")

	params := &scm.ContentParams{
		Message: "my commit message",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
		Signature: scm.Signature{
			Name:  "Monalisa Octocat",
			Email: "octocat@github.com",
		},
	}

	client := NewDefault()
	res, err := client.Contents.Create(
		context.Background(),
		"octocat/hello-world",
		"README",
		params,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 200 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost:7990").
		Put("/rest/api/1.0/projects/octocat/repos/hello-world/browse/README").
		Reply(200).
		Type("application/json").
		File("testdata/content_update.json")

	params := &scm.ContentParams{
		Message: "a new commit message",
		Data:    []byte("bXkgdXBkYXRlZCBmaWxlIGNvbnRlbnRz"),
		BlobID:  "95b966ae1c166bd92f8ae7d1c313e738c731dfc3",
		Signature: scm.Signature{
			Name:  "Monalisa Octocat",
			Email: "octocat@github.com",
		},
	}

	client := NewDefault()
	res, err := client.Contents.Update(
		context.Background(),
		"octocat/hello-world",
		"README",
		params,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 200 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentDelete(t *testing.T) {
	content := new(contentService)
	_, err := content.Delete(context.Background(), "atlassian/atlaskit", "README", &scm.ContentParams{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/drone/repos/go-scm/files/scm/driver/stash").
		MatchParam("at", "master").
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Contents.List(
		context.Background(),
		"drone/go-scm",
		"scm/driver/stash",
		"master",
		scm.ListOptions{},
	)
	if err != nil {
		t.Error(err)
	}

	want := []*scm.ContentInfo{}
	raw, _ := ioutil.ReadFile("testdata/content_list.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
