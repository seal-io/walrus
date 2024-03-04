// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/null"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(user)
	res, err := s.client.do(ctx, "GET", "user", nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := fmt.Sprintf("users/%s", login)
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	out, res, err := s.ListEmail(ctx, scm.ListOptions{})
	return returnPrimaryEmail(out), res, err
}

func (s *userService) ListEmail(ctx context.Context, opts scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	path := fmt.Sprintf("user/emails?%s", encodeListOptions(opts))
	out := []*email{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertEmailList(out), res, err
}

type user struct {
	ID      int         `json:"id"`
	Login   string      `json:"login"`
	Name    string      `json:"name"`
	Email   null.String `json:"email"`
	Avatar  string      `json:"avatar_url"`
	Created time.Time   `json:"created_at"`
	Updated time.Time   `json:"updated_at"`
}

type email struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func convertUser(from *user) *scm.User {
	return &scm.User{
		Avatar:  from.Avatar,
		Email:   from.Email.String,
		Login:   from.Login,
		Name:    from.Name,
		Created: from.Created,
		Updated: from.Updated,
	}
}

func returnPrimaryEmail(from []*scm.Email) string {
	for _, v := range from {
		if v.Primary == true {
			return v.Value
		}
	}
	return ""
}

// helper function to convert from the github email list to
// the common email structure.
func convertEmailList(from []*email) []*scm.Email {
	to := []*scm.Email{}
	for _, v := range from {
		to = append(to, convertEmail(v))
	}
	return to
}

// helper function to convert from the github email structure to
// the common email structure.
func convertEmail(from *email) *scm.Email {
	return &scm.Email{
		Value:    from.Email,
		Primary:  from.Primary,
		Verified: from.Verified,
	}
}
