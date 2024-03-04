// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

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
			sha:  "2eac1cac02c325058cf959725c45b0612d3e8177",
			want: "https://gitee.com/kit101/drone-yml-test/commit/2eac1cac02c325058cf959725c45b0612d3e8177",
		},
		{
			path: "refs/pull/7/head",
			sha:  "00b76e8abd51ae6a96318b3450944b32995f9158",
			want: "https://gitee.com/kit101/drone-yml-test/pulls/7",
		},
		{
			path: "refs/tags/1.1",
			want: "https://gitee.com/kit101/drone-yml-test/tree/1.1",
		},
		{
			path: "refs/heads/master",
			want: "https://gitee.com/kit101/drone-yml-test/tree/master",
		},
	}

	for _, test := range tests {
		client := NewDefault()
		ref := scm.Reference{
			Path: test.path,
			Sha:  test.sha,
		}
		got, err := client.Linker.Resource(context.Background(), "kit101/drone-yml-test", ref)
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
			source: scm.Reference{Sha: "2eac1cac02c325058cf959725c45b0612d3e8177"},
			target: scm.Reference{Sha: "00b76e8abd51ae6a96318b3450944b32995f9158"},
			want:   "https://gitee.com/kit101/drone-yml-test/compare/2eac1cac02c325058cf959725c45b0612d3e8177...00b76e8abd51ae6a96318b3450944b32995f9158",
		},
		{
			source: scm.Reference{Path: "refs/heads/master"},
			target: scm.Reference{Sha: "00b76e8abd51ae6a96318b3450944b32995f9158"},
			want:   "https://gitee.com/kit101/drone-yml-test/compare/master...00b76e8abd51ae6a96318b3450944b32995f9158",
		},
		{
			source: scm.Reference{Sha: "00b76e8abd51ae6a96318b3450944b32995f9158"},
			target: scm.Reference{Path: "refs/heads/master"},
			want:   "https://gitee.com/kit101/drone-yml-test/compare/00b76e8abd51ae6a96318b3450944b32995f9158...master",
		},
		{
			target: scm.Reference{Path: "refs/pull/7/head"},
			want:   "https://gitee.com/kit101/drone-yml-test/pulls/7/files",
		},
	}

	for _, test := range tests {
		client := NewDefault()
		got, err := client.Linker.Diff(context.Background(), "kit101/drone-yml-test", test.source, test.target)
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

func TestLinkBase(t *testing.T) {
	if got, want := NewDefault().Linker.(*linker).base, "https://gitee.com/"; got != want {
		t.Errorf("Want url %s, got %s", want, got)
	}
}
