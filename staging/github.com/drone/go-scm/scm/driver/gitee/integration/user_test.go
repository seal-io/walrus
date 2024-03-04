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
// user sub-tests
//

func testUsers(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testUserFind(client))
	}
}

func testUserFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Users.FindLogin(context.Background(), "kit101")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("User", testUser(result))
	}
}

//
// struct sub-tests
//

func testUser(user *scm.User) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		if got, want := user.Login, "kit101"; got != want {
			t.Errorf("Want user Login %q, got %q", want, got)
		}
		if got, want := user.Name, "kit101"; got != want {
			t.Errorf("Want user Name %q, got %q", want, got)
		}
		if got, want := user.Avatar, "https://portrait.gitee.com/uploads/avatars/user/511/1535738_qkssk1711_1578953939.png"; got != want {
			t.Errorf("Want user Avatar %q, got %q", want, got)
		}
	}
}
