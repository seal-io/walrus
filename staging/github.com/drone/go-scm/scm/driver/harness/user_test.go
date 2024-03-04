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

func TestUsersFind(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		harnessUserOrigin := strings.Replace(gockOrigin, "code", "ng", 1)

		gock.New(harnessUserOrigin).
			Get("/gateway/ng/api/user/currentUser").
			Reply(200).
			Type("application/json").
			File("testdata/user.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.Users.Find(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	wantErr := json.Unmarshal(raw, &want)
	if wantErr != nil {
		t.Error(wantErr)
		return
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
