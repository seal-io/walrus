// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestUserFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/user").
		Reply(200).
		Type("application/json").
		SetHeader("X-Request-Id", "7cb049dbf1fafae67b4f0aa81ca7e870").
		SetHeader("X-Runtime", "0.040217").
		SetHeader("ETag", "W/\"c50ecf93b72a24474a76423d6d5c338c\"").
		File("testdata/user.json")

	client := NewDefault()
	got, res, err := client.Users.Find(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	// t.Run("Rate", testRate(res))
}

func TestUserFindEmail(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/user").
		Reply(200).
		Type("application/json").
		SetHeader("X-Request-Id", "7cb049dbf1fafae67b4f0aa81ca7e870").
		SetHeader("X-Runtime", "0.040217").
		SetHeader("ETag", "W/\"c50ecf93b72a24474a76423d6d5c338c\"").
		File("testdata/user.json")

	client := NewDefault()
	result, res, err := client.Users.FindEmail(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := result, "qkssk1711@163.com"; got != want {
		t.Errorf("Want user Email %q, got %q", want, got)
	}
	t.Run("Request", testRequest(res))
	// t.Run("Rate", testRate(res))
}

func TestUserFindLogin(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/users/kit101").
		Reply(200).
		Type("application/json").
		SetHeader("X-Request-Id", "7cb049dbf1fafae67b4f0aa81ca7e870").
		SetHeader("X-Runtime", "0.040217").
		File("testdata/user.json")

	client := NewDefault()
	got, res, err := client.Users.FindLogin(context.Background(), "kit101")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
		json.NewEncoder(os.Stdout).Encode(got)
	}

	t.Run("Request", testRequest(res))
	// t.Run("Rate", testRate(res))
}
