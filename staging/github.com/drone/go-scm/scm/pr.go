// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"context"
	"time"
)

type (
	// PullRequest represents a repository pull request.
	PullRequest struct {
		Number  int
		Title   string
		Body    string
		Sha     string
		Ref     string
		Source  string
		Target  string
		Fork    string
		Link    string
		Diff    string
		Closed  bool
		Merged  bool
		Merge   string
		Base    Reference
		Head    Reference
		Author  User
		Created time.Time
		Updated time.Time
		Labels  []Label
	}

	// PullRequestInput provides the input fields required for creating a pull request.
	PullRequestInput struct {
		Title  string
		Body   string
		Source string
		Target string
	}

	// PullRequestListOptions provides options for querying
	// a list of repository merge requests.
	PullRequestListOptions struct {
		Page   int
		Size   int
		Open   bool
		Closed bool
	}

	// Change represents a changed file.
	Change struct {
		Path         string
		Added        bool
		Renamed      bool
		Deleted      bool
		Sha          string
		BlobID       string
		PrevFilePath string
	}

	Label struct {
		Name  string
		Color string
	}

	// Milestone the milestone
	Milestone struct {
		Number      int
		ID          int
		Title       string
		Description string
		Link        string
		State       string
		DueDate     time.Time
	}

	// PullRequestService provides access to pull request resources.
	PullRequestService interface {
		// Find returns the repository pull request by number.
		Find(context.Context, string, int) (*PullRequest, *Response, error)

		// FindComment returns the pull request comment by id.
		FindComment(context.Context, string, int, int) (*Comment, *Response, error)

		// Find returns the repository pull request list.
		List(context.Context, string, PullRequestListOptions) ([]*PullRequest, *Response, error)

		// ListChanges returns the pull request changeset.
		ListChanges(context.Context, string, int, ListOptions) ([]*Change, *Response, error)

		// ListComments returns the pull request comment list.
		ListComments(context.Context, string, int, ListOptions) ([]*Comment, *Response, error)

		// ListCommits returns the pull request commit list.
		ListCommits(context.Context, string, int, ListOptions) ([]*Commit, *Response, error)

		// Merge merges the repository pull request.
		Merge(context.Context, string, int) (*Response, error)

		// Close closes the repository pull request.
		Close(context.Context, string, int) (*Response, error)

		// Create creates a new pull request.
		Create(context.Context, string, *PullRequestInput) (*PullRequest, *Response, error)

		// CreateComment creates a new pull request comment.
		CreateComment(context.Context, string, int, *CommentInput) (*Comment, *Response, error)

		// DeleteComment deletes an pull request comment.
		DeleteComment(context.Context, string, int, int) (*Response, error)
	}
)
