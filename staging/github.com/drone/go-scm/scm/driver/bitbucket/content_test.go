// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

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

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/src/425863f9dbe56d70c8dcdbf2e4e0805e85591fcc/README").
		Reply(200).
		Type("text/plain").
		File("testdata/content.txt")

	gock.New("https://api.bitbucket.org").
		MatchParam("format", "meta").
		Get("/2.0/repositories/atlassian/atlaskit/src/425863f9dbe56d70c8dcdbf2e4e0805e85591fcc/README").
		Reply(200).
		Type("application/json").
		File("testdata/content.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Contents.Find(context.Background(), "atlassian/atlaskit", "README", "425863f9dbe56d70c8dcdbf2e4e0805e85591fcc")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestContentFindNoMeta(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/src/425863f9dbe56d70c8dcdbf2e4e0805e85591fcc/README").
		Reply(200).
		Type("text/plain").
		File("testdata/content.txt")

	gock.New("https://api.bitbucket.org").
		MatchParam("format", "meta").
		Get("/2.0/repositories/atlassian/atlaskit/src/425863f9dbe56d70c8dcdbf2e4e0805e85591fcc/README").
		Reply(404).
		Type("application/json").
		File("testdata/content_fail.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Contents.Find(context.Background(), "atlassian/atlaskit", "README", "425863f9dbe56d70c8dcdbf2e4e0805e85591fcc")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.fail")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/atlassian/atlaskit/src").
		Reply(201).
		Type("application/json")

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
		"atlassian/atlaskit",
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

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/atlassian/atlaskit/src").
		Reply(201).
		Type("application/json")

	params := &scm.ContentParams{
		Message: "my commit message",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
		Signature: scm.Signature{
			Name:  "Monalisa Octocat",
			Email: "octocat@github.com",
		},
	}

	client := NewDefault()
	res, err := client.Contents.Update(
		context.Background(),
		"atlassian/atlaskit",
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

func TestContentUpdateBadCommitID(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/atlassian/atlaskit/src").
		Reply(400).
		Type("application/json").
		File("testdata/content_update.json.fail")

	params := &scm.ContentParams{
		Message: "my commit message",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
		Sha:     "bad commit",
		Signature: scm.Signature{
			Name:  "Monalisa Octocat",
			Email: "octocat@github.com",
		},
	}

	client := NewDefault()
	_, err := client.Contents.Update(
		context.Background(),
		"atlassian/atlaskit",
		"test/hello",
		params,
	)
	if err.Error() != "parents: Commit not found: 1a7eba6c-d4fe-47b7-b767-859abc660efc" {
		t.Errorf("Expecting 'parents: Commit not found: 1a7eba6c-d4fe-47b7-b767-859abc660efc'")
	}
}

func TestContentDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/atlassian/atlaskit/src").
		Reply(201).
		Type("application/json")

	params := &scm.ContentParams{
		Message: "my commit message",
		Signature: scm.Signature{
			Name:  "Monalisa Octocat",
			Email: "octocat@github.com",
		},
	}

	client := NewDefault()
	res, err := client.Contents.Update(
		context.Background(),
		"atlassian/atlaskit",
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

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/src/master/packages/activity").
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Contents.List(context.Background(), "atlassian/atlaskit", "packages/activity", "master", scm.ListOptions{})
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

func TestContentListWithUrlInput(t *testing.T) {
	defer gock.Off()

	mockNextPageUri := "https://api.bitbucket.org/2.0/repositories/atlassian/atlaskit/src/master/packages/activity?pageLen=3&page=RPfL"

	gock.New(mockNextPageUri).
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Contents.List(context.Background(), "atlassian/atlaskit", "packages/activity", "master", scm.ListOptions{URL: mockNextPageUri})
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
