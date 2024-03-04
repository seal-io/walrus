// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type linker struct {
	base string
}

// Resource returns a link to the resource.
func (l *linker) Resource(ctx context.Context, repo string, ref scm.Reference) (string, error) {
	switch {
	case scm.IsTag(ref.Path):
		t := scm.TrimRef(ref.Path)
		return fmt.Sprintf("%s%s/tree/%s", l.base, repo, t), nil
	case scm.IsPullRequest(ref.Path):
		d := scm.ExtractPullRequest(ref.Path)
		return fmt.Sprintf("%s%s/pull/%d", l.base, repo, d), nil
	case ref.Sha == "":
		t := scm.TrimRef(ref.Path)
		return fmt.Sprintf("%s%s/tree/%s", l.base, repo, t), nil
	default:
		return fmt.Sprintf("%s%s/commit/%s", l.base, repo, ref.Sha), nil
	}
}

// Diff returns a link to the diff.
func (l *linker) Diff(ctx context.Context, repo string, source, target scm.Reference) (string, error) {
	if scm.IsPullRequest(target.Path) {
		d := scm.ExtractPullRequest(target.Path)
		return fmt.Sprintf("%s%s/pull/%d/files", l.base, repo, d), nil
	}

	s := source.Sha
	t := target.Sha
	if s == "" {
		s = scm.TrimRef(source.Path)
	}
	if t == "" {
		t = scm.TrimRef(target.Path)
	}

	return fmt.Sprintf("%s%s/compare/%s...%s", l.base, repo, s, t), nil
}
