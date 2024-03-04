// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

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
	namespace, name := scm.Split(repo)
	switch {
	case scm.IsTag(ref.Path):
		return fmt.Sprintf("%sprojects/%s/repos/%s/browse?at=%s", l.base, namespace, name, ref.Path), nil
	case scm.IsPullRequest(ref.Path):
		d := scm.ExtractPullRequest(ref.Path)
		return fmt.Sprintf("%sprojects/%s/repos/%s/pull-requests/%d/overview", l.base, namespace, name, d), nil
	case ref.Sha == "":
		return fmt.Sprintf("%sprojects/%s/repos/%s/browse?at=%s", l.base, namespace, name, ref.Path), nil
	default:
		return fmt.Sprintf("%sprojects/%s/repos/%s/commits/%s", l.base, namespace, name, ref.Sha), nil
	}
}

// Diff returns a link to the diff.
func (l *linker) Diff(ctx context.Context, repo string, source, target scm.Reference) (string, error) {
	namespace, name := scm.Split(repo)
	if scm.IsPullRequest(target.Path) {
		d := scm.ExtractPullRequest(target.Path)
		return fmt.Sprintf("%sprojects/%s/repos/%s/pull-requests/%d/diff", l.base, namespace, name, d), nil
	}
	if target.Path != "" && source.Path != "" {
		return fmt.Sprintf("%sprojects/%s/repos/%s/compare/diff?sourceBranch=%s&targetBranch=%s", l.base, namespace, name, source.Path, target.Path), nil
	}
	// TODO(bradrydzewski) bitbucket server does not appear to have
	// an endpoint for evaluating diffs of two commits.
	return "", scm.ErrNotSupported
}
