// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(user)
	res, err := s.client.do(ctx, "GET", "2.0/user", nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := fmt.Sprintf("2.0/users/%s", login)
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
    out := new(emails)
    res, err := s.client.do(ctx, "GET", "2.0/user/emails", nil, &out)
    return convertEmailList(out), res, err
}

func (s *userService) ListEmail(context.Context, scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func convertEmailList(from *emails) string {
	for _, v := range from.Values {
		if v.IsPrimary == true {
		return v.Email
		}
	}
	return ""
}

type user struct {
	// The `username` field is no longer available after 29 April 2019 in
	// accordance with GDPR regulations. See:
	// https://developer.atlassian.com/cloud/bitbucket/bitbucket-api-changes-gdpr/
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AccountID   string `json:"account_id"`
	Nickname    string `json:"nickname"`
	Links       struct {
		Self   link `json:"self"`
		HTML   link `json:"html"`
		Avatar link `json:"avatar"`
	} `json:"links"`
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

type email struct {
    Email string `json:"email"`
    IsPrimary bool `json:"is_primary"`
}

type emails struct {
    Values []*email `json:"values"`
}

func convertUser(from *user) *scm.User {
	return &scm.User{
		Avatar: fmt.Sprintf("https://bitbucket.org/account/%s/avatar/32/", from.Username),
		Login:  from.Username,
		Name:   from.DisplayName,
	}
}
