// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

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
			want: "https://try.gogs.io/octocat/hello-world/commit/a7389057b0eb027e73b32a81e3c5923a71d01dde",
		},
		{
			path: "refs/pull/42/head",
			sha:  "a7389057b0eb027e73b32a81e3c5923a71d01dde",
			want: "https://try.gogs.io/octocat/hello-world/pulls/42",
		},
		{
			path: "refs/tags/v1.0.0",
			want: "https://try.gogs.io/octocat/hello-world/src/v1.0.0",
		},
		{
			path: "refs/heads/master",
			want: "https://try.gogs.io/octocat/hello-world/src/master",
		},
	}

	for _, test := range tests {
		client, _ := New("https://try.gogs.io")
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
			want:   "https://try.gogs.io/octocat/hello-world/compare/a7389057b0eb027e73b32a81e3c5923a71d01dde...49bbaf4a113bbebfa21cf604cad9aa1503c3f04d",
		},
		{
			source: scm.Reference{Path: "refs/heads/master"},
			target: scm.Reference{Sha: "49bbaf4a113bbebfa21cf604cad9aa1503c3f04d"},
			want:   "https://try.gogs.io/octocat/hello-world/compare/master...49bbaf4a113bbebfa21cf604cad9aa1503c3f04d",
		},
		{
			source: scm.Reference{Sha: "a7389057b0eb027e73b32a81e3c5923a71d01dde"},
			target: scm.Reference{Path: "refs/heads/master"},
			want:   "https://try.gogs.io/octocat/hello-world/compare/a7389057b0eb027e73b32a81e3c5923a71d01dde...master",
		},
		{
			target: scm.Reference{Path: "refs/pull/12/head"},
			want:   "https://try.gogs.io/octocat/hello-world/pulls/12/files",
		},
	}

	for _, test := range tests {
		client, _ := New("https://try.gogs.io")
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

func TestLink_Base(t *testing.T) {
	client, _ := New("https://foo.bar.com/baz")
	ref := scm.Reference{
		Path: "refs/heads/master",
		Sha:  "a7389057b0eb027e73b32a81e3c5923a71d01dde",
	}
	got, err := client.Linker.Resource(context.Background(), "octocat/hello-world", ref)
	if err != nil {
		t.Error(err)
		return
	}
	want := "https://foo.bar.com/baz/octocat/hello-world/commit/a7389057b0eb027e73b32a81e3c5923a71d01dde"
	if got != want {
		t.Errorf("Want link %q, got %q", want, got)
	}
}
