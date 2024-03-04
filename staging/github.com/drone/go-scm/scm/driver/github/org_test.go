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

func TestOrganizationFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/orgs/github").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/org.json")

	client := NewDefault()
	got, res, err := client.Organizations.Find(context.Background(), "github")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Organization)
	raw, _ := ioutil.ReadFile("testdata/org.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestOrganizationFindMembership(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/orgs/github/memberships/octocat").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/membership.json")

	client := NewDefault()
	got, res, err := client.Organizations.FindMembership(context.Background(), "github", "octocat")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Membership)
	raw, _ := ioutil.ReadFile("testdata/membership.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestOrganizationList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/user/orgs").
		MatchParam("per_page", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/orgs.json")

	client := NewDefault()
	got, res, err := client.Organizations.List(context.Background(), scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Organization{}
	raw, _ := ioutil.ReadFile("testdata/orgs.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}
