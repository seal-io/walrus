// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

const (
	gockOrigin     = "https://qa.harness.io/gateway/code"
	harnessOrg     = "px7xd_BFRCi-pfWPYXVjvw"
	harnessAccount = "default"
	harnessProject = "codeciintegration"
	harnessRepo    = "thomas"
	harnessPAT     = ""
)

func TestContentFind(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/content/README.md").
			Reply(200).
			Type("plain/text").
			File("testdata/content.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	result, _, err := client.Contents.Find(
		context.Background(),
		harnessRepo,
		"README.md",
		"98189d5cf2a751a6246c24a72945ba70839f1b20",
	)
	if err != nil {
		t.Error(err)
	}

	if got, want := result.Path, "README.md"; got != want {
		t.Errorf("Want file Path %q, got %q", want, got)
	}
	if !strings.Contains(string(result.Data), "project") {
		t.Errorf("Want file Data %q, must contain 'project'", result.Data)
	}
}

func TestContentCreate(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Post("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/commits").
			Reply(200).
			Type("plain/text").
			BodyString("{\"commit_id\":\"20ecde1f8c277da0e91750bef9f3b88f228d86db\"}")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	result, err := client.Contents.Create(
		context.Background(),
		harnessRepo,
		"README.2",
		&scm.ContentParams{
			Data:    []byte("hello world"),
			Message: "create README.2",
			Branch:  "main",
		},
	)
	if err != nil {
		t.Error(err)
	}

	if result.Status != 200 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentUpdate(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Post("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/commits").
			Reply(200).
			Type("plain/text").
			BodyString("{\"commit_id\":\"20ecde1f8c277da0e91750bef9f3b88f228d86db\"}")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	result, err := client.Contents.Update(
		context.Background(),
		harnessRepo,
		"README.2",
		&scm.ContentParams{
			Data:    []byte("hello world 2"),
			Message: "update README.2",
			Branch:  "main",
			BlobID:  "95d09f2b10159347eece71399a7e2e907ea3df4f",
		},
	)
	if err != nil {
		t.Error(err)
	}

	if result.Status != 200 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentDelete(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Post("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/commits").
			Reply(200).
			Type("plain/text").
			BodyString("{\"commit_id\":\"20ecde1f8c277da0e91750bef9f3b88f228d86db\"}")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	result, err := client.Contents.Delete(
		context.Background(),
		harnessRepo,
		"README.2",
		&scm.ContentParams{
			Message: "delete README.2",
			Branch:  "main",
		},
	)
	if err != nil {
		t.Error(err)
	}

	if result.Status != 200 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New(gockOrigin).
		Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/content/docker").
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Contents.List(
		context.Background(),
		harnessRepo,
		"docker",
		"",
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
