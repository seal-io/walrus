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

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/contents/README.md").
		MatchParam("ref", "d295a4c616d46fbcdfa3dfd1473c1337a1ec6f83").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	client := NewDefault()
	got, res, err := client.Contents.Find(
		context.Background(),
		"kit101/drone-yml-test",
		"README.md",
		"d295a4c616d46fbcdfa3dfd1473c1337a1ec6f83",
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

	t.Run("Request", testRequest(res))
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/drone-yml-test/contents/apitest/CreateByDroneGiteeProvider.md").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_create.json")

	params := &scm.ContentParams{
		Message: "my commit message",
		Data:    []byte("bY3JlYXRlIGJ5IGRyb25lLXNjbSBnaXRlZSBwcm92aWRlci4gMjAyMS0wOC0xNyAyMzowNToxNi4="),
		Signature: scm.Signature{
			Name:  "kit101",
			Email: "kit101@gitee.com",
		},
	}

	client := NewDefault()
	res, err := client.Contents.Create(
		context.Background(),
		"kit101/drone-yml-test",
		"apitest/CreateByDroneGiteeProvider.md",
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

	gock.New("https://gitee.com/api/v5").
		Put("/repos/kit101/drone-yml-test/contents/apitest/UpdateByDroneGiteeProvider.md").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_update.json")

	params := &scm.ContentParams{
		Message: "my commit message by update",
		Data:    []byte("bdXBkYXRlIGJ5IGRyb25lLXNjbSBnaXRlZSBwcm92aWRlci4gMjAyMS0wOC0xNyAyMzozMDozNy4="),
		Sha:     "9de0cf94e1e3c1cbe0a25c3865de4cc9ede7ad3e",
		Signature: scm.Signature{
			Name:  "kit101",
			Email: "kit101@gitee.com",
		},
	}

	client := NewDefault()
	res, err := client.Contents.Update(
		context.Background(),
		"kit101/drone-yml-test",
		"apitest/UpdateByDroneGiteeProvider.md",
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
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Delete("/repos/kit101/drone-yml-test/contents/apitest/DeleteByDroneGiteeProvider.md").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_delete.json")

	contentParams := scm.ContentParams{
		Branch:  "master",
		Sha:     "ae0653e4ab697cd77fc559cd798bdb30e8bb4d7e",
		Message: "delete commit message",
		Signature: scm.Signature{
			Name:  "kit101",
			Email: "kit101@gitee.com",
		},
	}
	client := NewDefault()
	res, err := client.Contents.Delete(
		context.Background(),
		"kit101/drone-yml-test",
		"apitest/DeleteByDroneGiteeProvider.md",
		&contentParams,
	)
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 204; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/contents/apitest").
		MatchParam("ref", "master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_list.json")

	client := NewDefault()
	got, res, err := client.Contents.List(
		context.Background(),
		"kit101/drone-yml-test",
		"apitest",
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
}
