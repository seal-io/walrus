// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

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
			want: "https://bitbucket.org/octocat/hello-world/commits/a7389057b0eb027e73b32a81e3c5923a71d01dde",
		},
		{
			path: "refs/pull/42/head",
			sha:  "a7389057b0eb027e73b32a81e3c5923a71d01dde",
			want: "https://bitbucket.org/octocat/hello-world/pull-requests/42",
		},
		{
			path: "refs/tags/v1.0.0",
			want: "https://bitbucket.org/octocat/hello-world/src/v1.0.0",
		},
		{
			path: "refs/heads/master",
			want: "https://bitbucket.org/octocat/hello-world/src/master",
		},
		{
			path: "refs/heads/release/production",
			want: "https://bitbucket.org/octocat/hello-world/branch/release/production",
		},
		{
			path: "refs/heads/release/production",
			sha:  "a7389057b0eb027e73b32a81e3c5923a71d01dde",
			want: "https://bitbucket.org/octocat/hello-world/commits/a7389057b0eb027e73b32a81e3c5923a71d01dde",
		},
	}

	for _, test := range tests {
		client := NewDefault()
		ref := scm.Reference{
			Path: test.path,
			Sha:  test.sha,
		}
		got, err := client.Linker.Resource(context.Background(), "octocat/hello-world", ref)
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
			source: scm.Reference{Sha: "a7389057b0eb027e73b32a81e3c5923a71d01dde"},
			target: scm.Reference{Sha: "49bbaf4a113bbebfa21cf604cad9aa1503c3f04d"},
			want:   "https://bitbucket.org/octocat/hello-world/compare/a7389057b0eb027e73b32a81e3c5923a71d01dde%0D49bbaf4a113bbebfa21cf604cad9aa1503c3f04d",
		},
		{
			source: scm.Reference{Path: "refs/heads/master"},
			target: scm.Reference{Sha: "49bbaf4a113bbebfa21cf604cad9aa1503c3f04d"},
			want:   "https://bitbucket.org/octocat/hello-world/compare/master%0D49bbaf4a113bbebfa21cf604cad9aa1503c3f04d",
		},
		{
			source: scm.Reference{Sha: "a7389057b0eb027e73b32a81e3c5923a71d01dde"},
			target: scm.Reference{Path: "refs/heads/master"},
			want:   "https://bitbucket.org/octocat/hello-world/compare/a7389057b0eb027e73b32a81e3c5923a71d01dde%0Dmaster",
		},
		{
			target: scm.Reference{Path: "refs/pull/12/head"},
			want:   "https://bitbucket.org/octocat/hello-world/pull-requests/12",
		},
	}

	for _, test := range tests {
		client := NewDefault()
		got, err := client.Linker.Diff(context.Background(), "octocat/hello-world", test.source, test.target)
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
