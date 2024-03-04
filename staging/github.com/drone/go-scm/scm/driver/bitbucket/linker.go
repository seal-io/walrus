// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"
	"strings"

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
		return fmt.Sprintf("%s%s/src/%s", l.base, repo, t), nil
	case scm.IsPullRequest(ref.Path):
		d := scm.ExtractPullRequest(ref.Path)
		return fmt.Sprintf("%s%s/pull-requests/%d", l.base, repo, d), nil
	case ref.Sha == "":
		t := scm.TrimRef(ref.Path)

		// Bitbucket has a bug where the "source view" link for
		// a branch which contains a slash results in a 404.
		// The link to the "branch view" works with names containing
		// a slash so we do this for branches with slashes in its names.
		// See https://jira.atlassian.com/browse/BCLOUD-14422 for more information
		if scm.IsBranch(ref.Path) && strings.Contains(t, "/") {
			return fmt.Sprintf("%s%s/branch/%s", l.base, repo, t), nil
		}

		return fmt.Sprintf("%s%s/src/%s", l.base, repo, t), nil
	default:
		return fmt.Sprintf("%s%s/commits/%s", l.base, repo, ref.Sha), nil
	}
}

// Diff returns a link to the diff.
func (l *linker) Diff(ctx context.Context, repo string, source, target scm.Reference) (string, error) {
	if scm.IsPullRequest(target.Path) {
		d := scm.ExtractPullRequest(target.Path)
		return fmt.Sprintf("%s%s/pull-requests/%d", l.base, repo, d), nil
	}

	s := source.Sha
	t := target.Sha
	if s == "" {
		s = scm.TrimRef(source.Path)
	}
	if t == "" {
		t = scm.TrimRef(target.Path)
	}

	return fmt.Sprintf("%s%s/compare/%s%%0D%s", l.base, repo, s, t), nil
}
