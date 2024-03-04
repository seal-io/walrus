// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/h2non/gock"
)

func TestRepositoryFind(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/demo/+").
			Reply(200).
			Type("application/json").
			File("testdata/repo.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Repositories.Find(context.Background(), "demo")
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
}

func TestRepositoryList(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/spaces/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/+/repos").
			MatchParam("page", "1").
			MatchParam("limit", "20").
			MatchParam("sort", "path").
			MatchParam("order", "asc").
			Reply(200).
			Type("application/json").
			File("testdata/repos.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Repositories.List(context.Background(), scm.ListOptions{Page: 1, Size: 20})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	_ = json.Unmarshal(raw, &want)

	if harnessPAT != "" && len(got) > 0 {
		// pass when running against a live harness instance and we get more than one repo
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryHookList(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/webhooks").
			MatchParam("page", "1").
			MatchParam("limit", "30").
			MatchParam("sort", "display_name").
			MatchParam("order", "asc").
			Reply(200).
			Type("application/json").
			File("testdata/hooks.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Repositories.ListHooks(context.Background(), harnessRepo, scm.ListOptions{Page: 1, Size: 30})
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
}

func TestRepositoryFindHook(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/webhooks/6").
			Reply(200).
			Type("application/json").
			File("testdata/hook.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Repositories.FindHook(context.Background(), harnessRepo, "6")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryHookCreateDelete(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Post("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/webhooks").
			Reply(200).
			Type("application/json").
			File("testdata/hook_create.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	in := &scm.HookInput{
		Name:       "drone",
		Target:     "https://example.com",
		Secret:     "topsecret",
		SkipVerify: true,
	}
	got, _, err := client.Repositories.CreateHook(context.Background(), harnessRepo, in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook_create.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want, cmpopts.IgnoreFields(scm.Hook{}, "ID")); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
	// delete webhook
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Delete(fmt.Sprintf("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/webhooks/%s", got.ID)).
			Reply(204)
	}
	client, _ = New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}

	_, deleteErr := client.Repositories.DeleteHook(context.Background(), harnessRepo, got.ID)
	if deleteErr != nil {
		t.Error(deleteErr)
		return
	}
}
