// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

//
// organization sub-tests
//

func testOrgs(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testOrgFind(client))
	}
}

func testOrgFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Organizations.Find(context.Background(), "gitee-community")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("Organization", testOrg(result))
	}
}

//
// struct sub-tests
//

func testOrg(organization *scm.Organization) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		if got, want := organization.Name, "gitee-community"; got != want {
			t.Errorf("Want organization Name %q, got %q", want, got)
		}
		if got, want := organization.Avatar, "https://portrait.gitee.com/uploads/avatars/namespace/1715/5146011_gitee-community_1607593452.png?is_link=true"; got != want {
			t.Errorf("Want organization Avatar %q, got %q", want, got)
		}
	}
}
