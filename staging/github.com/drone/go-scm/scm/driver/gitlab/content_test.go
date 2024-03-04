// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

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

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/files/app/models/key.rb").
		MatchParam("ref", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	client := NewDefault()
	got, res, err := client.Contents.Find(
		context.Background(),
		"diaspora/diaspora",
		"app/models/key.rb",
		"7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/files/app/models/key.rb").
		MatchParam("ref", "b1&b2").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	client = NewDefault()
	got, res, err = client.Contents.Find(
		context.Background(),
		"diaspora/diaspora",
		"app/models/key.rb",
		"b1&b2",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want = new(scm.Content)
	raw, _ = ioutil.ReadFile("testdata/content.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Post("/api/v4/projects/diaspora/diaspora/repository/files/app/project.rb").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_create.json")

	client := NewDefault()
	params := &scm.ContentParams{
		Message: "create a new file",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
		Signature: scm.Signature{
			Name:  "Firstname Lastname",
			Email: "kubesphere@example.com",
		},
	}

	res, err := client.Contents.Create(context.Background(), "diaspora/diaspora", "app/project.rb", params)
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

	gock.New("https://gitlab.com").
		Put("/api/v4/projects/diaspora/diaspora/repository/files/app/project.rb").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_update.json")

	client := NewDefault()
	params := &scm.ContentParams{
		Message: "update file",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
		Signature: scm.Signature{
			Name:  "Firstname Lastname",
			Email: "kubesphere@example.com",
		},
	}

	res, err := client.Contents.Update(context.Background(), "diaspora/diaspora", "app/project.rb", params)
	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 200 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentUpdateBadCommitID(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Put("/api/v4/projects/diaspora/diaspora/repository/files/app/project.rb").
		Reply(400).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_update.json.fail")

	client := NewDefault()
	params := &scm.ContentParams{
		Message: "update file",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
		Sha:     "bad sha",
		Signature: scm.Signature{
			Name:  "Firstname Lastname",
			Email: "kubesphere@example.com",
		},
	}

	_, err := client.Contents.Update(context.Background(), "diaspora/diaspora", "app/project.rb", params)
	if err.Error() != "You are attempting to update a file that has changed since you started editing it." {
		t.Errorf("Expecting error 'You are attempting to update a file that has changed since you started editing it.'")
	}
}

func TestContentDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Delete("/api/v4/projects/diaspora/diaspora/repository/files/app/project.rb").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	params := &scm.ContentParams{
		Message: "update file",
		Signature: scm.Signature{
			Name:  "Firstname Lastname",
			Email: "kubesphere@example.com",
		},
	}

	res, err := client.Contents.Delete(context.Background(), "diaspora/diaspora", "app/project.rb", params)
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

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/gitlab-org/gitlab/repository/tree").
		MatchParam("path", "lib/gitlab/ci").
		MatchParam("ref", "master").
		Reply(200).
		SetHeaders(mockHeaders).
		File("testdata/content_list.json")

	client := NewDefault()
	got, res, err := client.Contents.List(
		context.Background(),
		"gitlab-org/gitlab",
		"lib/gitlab/ci",
		"master",
		scm.ListOptions{},
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.ContentInfo{}
	raw, _ := ioutil.ReadFile("testdata/content_list.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
