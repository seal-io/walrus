// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/orgs/%s", name)
	out := new(org)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertOrg(out), res, err
}

func (s *organizationService) FindMembership(ctx context.Context, name, username string) (*scm.Membership, *scm.Response, error) {
	membership := new(membership)
	membership.Active = s.checkMembership(ctx, name, username)
	out := new(permissions)
	path := fmt.Sprintf("api/v1/users/%s/orgs/%s/permissions", username, name)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	membership.Permissions = out
	return convertMembership(membership), res, err
}

func (s *organizationService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/user/orgs?%s", encodeListOptions(opts))
	out := []*org{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertOrgList(out), res, err
}

type permissions struct {
	IsOwner             bool `json:"is_owner"`
	IsAdmin             bool `json:"is_admin"`
	CanWrite            bool `json:"can_write"`
	CanRead             bool `json:"can_read"`
	CanCreateRepository bool `json:"can_create_repository"`
}
type membership struct {
	Permissions *permissions
	Active      bool
}

//
// native data structures
//

type org struct {
	Name   string `json:"username"`
	Avatar string `json:"avatar_url"`
}

//
// native data structure conversion
//

func convertOrgList(from []*org) []*scm.Organization {
	to := []*scm.Organization{}
	for _, v := range from {
		to = append(to, convertOrg(v))
	}
	return to
}

func convertOrg(from *org) *scm.Organization {
	return &scm.Organization{
		Name:   from.Name,
		Avatar: from.Avatar,
	}
}

func (s *organizationService) checkMembership(ctx context.Context, name, username string) bool {
	path := fmt.Sprintf("api/v1/orgs/%s/members/%s", name, username)
	res, err := s.client.do(ctx, "GET", path, nil, nil)
	if err != nil {
		return false
	}
	return res.Status == 204
}

func convertMembership(from *membership) *scm.Membership {
	to := new(scm.Membership)
	to.Active = from.Active
	isAdmin := from.Permissions.IsAdmin
	isMember := from.Permissions.CanRead || from.Permissions.CanWrite || from.Permissions.CanCreateRepository
	if isAdmin {
		to.Role = scm.RoleAdmin
	} else if isMember {
		to.Role = scm.RoleMember
	} else {
		to.Role = scm.RoleUndefined
	}
	return to
}
