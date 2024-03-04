// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import "context"

type (
	// Content stores the contents of a repository file.
	Content struct {
		Path   string
		Data   []byte
		Sha    string
		BlobID string
	}

	// ContentParams provide parameters for creating and
	// updating repository content.
	ContentParams struct {
		Ref       string
		Branch    string
		Message   string
		Data      []byte
		Sha       string
		BlobID    string
		Signature Signature
	}

	// ContentInfo stores the kind of any content in a repository.
	ContentInfo struct {
		Path   string
		Sha    string
		BlobID string
		Kind   ContentKind
	}

	// ContentService provides access to repositroy content.
	ContentService interface {
		// Find returns the repository file content by path.
		Find(ctx context.Context, repo, path, ref string) (*Content, *Response, error)

		// Create creates a new repositroy file.
		Create(ctx context.Context, repo, path string, params *ContentParams) (*Response, error)

		// Update updates a repository file.
		Update(ctx context.Context, repo, path string, params *ContentParams) (*Response, error)

		// Delete deletes a reository file.
		Delete(ctx context.Context, repo, path string, params *ContentParams) (*Response, error)

		// List returns a list of contents in a repository directory by path. It is
		// up to the driver to list the directory recursively or non-recursively,
		// but a robust driver should return a non-recursive list if possible.
		List(ctx context.Context, repo, path, ref string, opts ListOptions) ([]*ContentInfo, *Response, error)
	}
)
