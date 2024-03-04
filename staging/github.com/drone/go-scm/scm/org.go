// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"context"
)

type (
	// Organization represents an organization account.
	Organization struct {
		Name   string
		Avatar string
	}

	// Membership represents an organization membership.
	Membership struct {
		Active bool
		Role   Role
	}

	// OrganizationService provides access to organization resources.
	OrganizationService interface {
		// Find returns the organization by name.
		Find(ctx context.Context, name string) (*Organization, *Response, error)

		// FindMembership returns the organization membership
		// by a given user account.
		FindMembership(ctx context.Context, name, username string) (*Membership, *Response, error)

		// List returns the user organization list.
		List(ctx context.Context, opts ListOptions) ([]*Organization, *Response, error)
	}
)
