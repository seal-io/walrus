// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"context"
	"time"
)

// EmptyCommit is an empty commit sha.
const EmptyCommit = "0000000000000000000000000000000000000000"

type (
	// Reference represents a git reference.
	Reference struct {
		Name string
		Path string
		Sha  string
	}

	// ReferenceInput provides a SHA for creating a reference.
	ReferenceInput struct {
		Name string
		Sha  string
	}

	// Commit represents a repository commit.
	Commit struct {
		Sha       string
		Message   string
		Author    Signature
		Committer Signature
		Link      string
	}

	// CommitListOptions provides options for querying a
	// list of repository commits.
	CommitListOptions struct {
		Ref  string
		Page int
		Size int
		Path string
	}

	// Signature identifies a git commit creator.
	Signature struct {
		Name  string
		Email string
		Date  time.Time

		// Fields are optional. The provider may choose to
		// include account information in the response.
		Login  string
		Avatar string
	}

	// GitService provides access to git resources.
	GitService interface {
		// CreateBranch creates a git branch by name given a sha.
		CreateBranch(ctx context.Context, repo string, params *ReferenceInput) (*Response, error)

		// FindBranch finds a git branch by name.
		FindBranch(ctx context.Context, repo, name string) (*Reference, *Response, error)

		// FindCommit finds a git commit by ref.
		FindCommit(ctx context.Context, repo, ref string) (*Commit, *Response, error)

		// FindTag finds a git tag by name.
		FindTag(ctx context.Context, repo, name string) (*Reference, *Response, error)

		// ListBranches returns a list of git branches.
		ListBranches(ctx context.Context, repo string, opts ListOptions) ([]*Reference, *Response, error)

		// ListBranchesV2 returns a list of git branches based on the searchTerm passed.
		ListBranchesV2(ctx context.Context, repo string, opts BranchListOptions) ([]*Reference, *Response, error)

		// ListCommits returns a list of git commits.
		ListCommits(ctx context.Context, repo string, opts CommitListOptions) ([]*Commit, *Response, error)

		// ListChanges returns the changeset of a commit.
		ListChanges(ctx context.Context, repo, ref string, opts ListOptions) ([]*Change, *Response, error)

		// ListTags returns a list of git tags.
		ListTags(ctx context.Context, repo string, opts ListOptions) ([]*Reference, *Response, error)

		// CompareChanges returns the changeset between two
		// commits. If the source commit is not an ancestor
		// of the target commit, it is up to the driver to
		// return a 2-way or 3-way diff changeset.
		CompareChanges(ctx context.Context, repo, source, target string, opts ListOptions) ([]*Change, *Response, error)
	}

	// CreateBranch is a type alias for upstream projects
	// that use the previous CreateBranch type name.
	CreateBranch = ReferenceInput
)
