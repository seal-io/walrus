// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package gitea implements a Gogs client.
package harness

import (
	"fmt"
	"testing"

	"github.com/drone/go-scm/scm"
)

func TestClient(t *testing.T) {
	client, err := New(gockOrigin, "", "", "")
	if err != nil {
		t.Error(err)
	}
	if got, want := client.BaseURL.String(), fmt.Sprintf("%s/", gockOrigin); got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Error(t *testing.T) {
	_, err := New("http://a b.com/", "", "", "")
	if err == nil {
		t.Errorf("Expect error when invalid URL")
	}
}

func testPage(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Page.Next, 2; got != want {
			t.Errorf("Want next page %d, got %d", want, got)
		}
		if got, want := res.Page.Prev, 1; got != want {
			t.Errorf("Want prev page %d, got %d", want, got)
		}
		if got, want := res.Page.First, 1; got != want {
			t.Errorf("Want first page %d, got %d", want, got)
		}
		if got, want := res.Page.Last, 5; got != want {
			t.Errorf("Want last page %d, got %d", want, got)
		}
	}
}
