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

func TestUserFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("plugins/servlet/applinks/whoami").
		Reply(200).
		Type("text/plain").
		BodyString("jcitizen")

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/jcitizen").
		Reply(200).
		Type("application/json").
		File("testdata/user.json")

	client, _ := New("https://bitbucket.example.com")
	got, _, err := client.Users.Find(context.Background())
	if err != nil {
		t.Error(err)
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestUserLoginFind_ViaSearch_NotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/jcitizen").
		Reply(404)

	client, _ := New("https://bitbucket.example.com")
	_, _, err := client.Users.FindLogin(context.Background(), "jcitizen")
	if err == nil {
		t.Errorf("Want ErrNotFound got nil error")
	}

	if !gock.IsDone() {
		t.Errorf("Pending mocks")
	}
}

func TestUserLoginFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/jcitizen").
		Reply(200).
		Type("application/json").
		File("testdata/user.json")

	client, _ := New("https://bitbucket.example.com")
	got, _, err := client.Users.FindLogin(context.Background(), "jcitizen")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	if !gock.IsDone() {
		t.Errorf("Pending mocks")
	}
}

func TestUserLoginFind_ViaSearch(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/jane@example").
		Reply(404)

	gock.New("https://bitbucket.example.com").
		Get("/rest/api/1.0/users").
		MatchParam("filter", "jane@example").
		Reply(200).
		Type("application/json").
		File("testdata/user_search.json")

	client, _ := New("https://bitbucket.example.com")
	got, _, err := client.Users.FindLogin(context.Background(), "jane@example")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user_search.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	if !gock.IsDone() {
		t.Errorf("Pending mocks")
	}
}

func TestUserLoginFind_ViaSearch_NoMatch(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/john@example").
		Reply(404)

	gock.New("https://bitbucket.example.com").
		Get("/rest/api/1.0/users").
		MatchParam("filter", "john@example").
		Reply(200).
		Type("application/json").
		File("testdata/user_search.json")

	client, _ := New("https://bitbucket.example.com")
	_, _, err := client.Users.FindLogin(context.Background(), "john@example")
	if err == nil {
		t.Errorf("Want ErrNotFound got nil error")
	}

	if !gock.IsDone() {
		t.Errorf("Pending mocks")
	}
}

func TestUserFindEmail(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("plugins/servlet/applinks/whoami").
		Reply(200).
		Type("text/plain").
		BodyString("jcitizen")

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/jcitizen").
		Reply(200).
		Type("application/json").
		File("testdata/user.json")

	client, _ := New("https://bitbucket.example.com")
	email, _, err := client.Users.FindEmail(context.Background())
	if err != nil {
		t.Error(err)
	}

	if got, want := email, "jane@example.com"; got != want {
		t.Errorf("Want email %s, got %s", want, got)
	}
}
