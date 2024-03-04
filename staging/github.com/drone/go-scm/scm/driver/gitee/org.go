// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("orgs/%s", name)
	out := new(organization)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertOrganization(out), res, err
}

func (s *organizationService) FindMembership(ctx context.Context, name, username string) (*scm.Membership, *scm.Response, error) {
	path := fmt.Sprintf("orgs/%s/memberships/%s", name, username)
	out := new(membership)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertMembership(out), res, err
}

func (s *organizationService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("user/orgs?%s", encodeListOptions(opts))
	out := []*organization{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertOrganizationList(out), res, err
}

type organization struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	AvatarURL   string `json:"avatar_url"`
	ReposURL    string `json:"repos_url"`
	EventsURL   string `json:"events_url"`
	MembersURL  string `json:"members_url"`
	Description string `json:"description"`
	FollowCount int    `json:"follow_count"`
}

type membership struct {
	URL             string       `json:"url"`
	Active          bool         `json:"active"`
	Remark          string       `json:"remark"`
	Role            string       `json:"role"`
	OrganizationURL string       `json:"organization_url"`
	Organization    organization `json:"organization"`
	User            struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		URL       string `json:"url"`
		HtmlURL   string `json:"html_url"`
		Remark    string `json:"remark"`
	} `json:"user"`
}

func convertOrganizationList(from []*organization) []*scm.Organization {
	to := []*scm.Organization{}
	for _, v := range from {
		to = append(to, convertOrganization(v))
	}
	return to
}

func convertOrganization(from *organization) *scm.Organization {
	return &scm.Organization{
		Name:   from.Login,
		Avatar: from.AvatarURL,
	}
}

func convertMembership(from *membership) *scm.Membership {
	to := new(scm.Membership)
	to.Active = from.Active
	switch from.Role {
	case "admin":
		to.Role = scm.RoleAdmin
	case "member":
		to.Role = scm.RoleMember
	default:
		to.Role = scm.RoleUndefined
	}
	return to
}
