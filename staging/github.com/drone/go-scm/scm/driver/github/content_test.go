// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/contents/README").
		MatchParam("ref", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	client := NewDefault()
	got, res, err := client.Contents.Find(
		context.Background(),
		"octocat/hello-world",
		"README",
		"7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/contents/README").
		MatchParam("ref", "b1&b2").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	client = NewDefault()
	got, res, err = client.Contents.Find(
		context.Background(),
		"octocat/hello-world",
		"README",
		"b1&b2",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want = new(scm.Content)
	raw, _ = ioutil.ReadFile("testdata/content.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Put("/repos/octocat/hello-world/contents/test/hello").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
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
		"test/hello",
		params,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 201 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Put("/repos/octocat/hello-world/contents/test/hello").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
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
		"test/hello",
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

func TestContentUpdateBadBlobID(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Put("/repos/octocat/hello-world/contents/test/hello").
		Reply(401).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_update.json.fail")

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
	_, err := client.Contents.Update(
		context.Background(),
		"octocat/hello-world",
		"test/hello",
		params,
	)
	if err.Error() != "newfile does not match" {
		t.Errorf("Expecting 'newfile does not match'")
	}
}

func TestContentDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Delete("/repos/octocat/hello-world/contents/test/hello").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_delete.json")

	params := &scm.ContentParams{
		Message: "a new commit message",
		BlobID:  "95b966ae1c166bd92f8ae7d1c313e738c731dfc3",
		Signature: scm.Signature{
			Name:  "Monalisa Octocat",
			Email: "octocat@github.com",
		},
	}

	client := NewDefault()
	res, err := client.Contents.Delete(
		context.Background(),
		"octocat/hello-world",
		"test/hello",
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

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/drone/go-scm/contents/scm/driver/github").
		MatchParam("ref", "master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_list.json")

	client := NewDefault()
	got, res, err := client.Contents.List(
		context.Background(),
		"drone/go-scm",
		"scm/driver/github",
		"master",
		scm.ListOptions{},
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.ContentInfo{}
	raw, _ := ioutil.ReadFile("testdata/content_list.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
