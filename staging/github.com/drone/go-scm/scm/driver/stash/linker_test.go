// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func TestLink(t *testing.T) {
	tests := []struct {
		path string
		sha  string
		want string
	}{
		{
			path: "refs/heads/master",
			sha:  "a7389057b0eb027e73b32a81e3c5923a71d01dde",
			want: "https://stash.acme.com/projects/PRJ/repos/hello-world/commits/a7389057b0eb027e73b32a81e3c5923a71d01dde",
		},
		{
			path: "refs/pull/42/head",
			sha:  "a7389057b0eb027e73b32a81e3c5923a71d01dde",
			want: "https://stash.acme.com/projects/PRJ/repos/hello-world/pull-requests/42/overview",
		},
		{
			path: "refs/tags/v1.0.0",
			want: "https://stash.acme.com/projects/PRJ/repos/hello-world/browse?at=refs/tags/v1.0.0",
		},
		{
			path: "refs/heads/master",
			want: "https://stash.acme.com/projects/PRJ/repos/hello-world/browse?at=refs/heads/master",
		},
	}

	for _, test := range tests {
		client, _ := New("https://stash.acme.com")
		ref := scm.Reference{
			Path: test.path,
			Sha:  test.sha,
		}
		got, err := client.Linker.Resource(context.Background(), "PRJ/hello-world", ref)
		if err != nil {
			t.Error(err)
			return
		}
		want := test.want
		if got != want {
			t.Errorf("Want link %q, got %q", want, got)
		}
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		source scm.Reference
		target scm.Reference
		want   string
	}{
		{
			source: scm.Reference{Path: "refs/heads/master"},
			target: scm.Reference{Path: "refs/heads/develop"},
			want:   "https://stash.acme.com/projects/PRJ/repos/hello-world/compare/diff?sourceBranch=refs/heads/master&targetBranch=refs/heads/develop",
		},
		{
			target: scm.Reference{Path: "refs/pull/12/head"},
			want:   "https://stash.acme.com/projects/PRJ/repos/hello-world/pull-requests/12/diff",
		},
	}

	for _, test := range tests {
		client, _ := New("https://stash.acme.com")
		got, err := client.Linker.Diff(context.Background(), "PRJ/hello-world", test.source, test.target)
		if err != nil {
			t.Error(err)
			return
		}
		want := test.want
		if got != want {
			t.Errorf("Want link %q, got %q", want, got)
		}
	}

	source := scm.Reference{}
	target := scm.Reference{}
	client, _ := New("https://stash.acme.com")
	_, err := client.Linker.Diff(context.Background(), "PRJ/hello-world", source, target)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect ErrNotSupported when refpath is empty")
	}
}
