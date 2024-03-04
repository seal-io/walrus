// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"fmt"
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/null"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(user)
	res, err := s.client.do(ctx, "GET", "api/v4/user", nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/users?search=%s", login)
	out := []*user{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if err != nil {
		return nil, nil, err
	}
	if len(out) != 1 || !strings.EqualFold(out[0].Username, login) {
		return nil, nil, scm.ErrNotFound
	}
	return convertUser(out[0]), res, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	user, res, err := s.Find(ctx)
	return user.Email, res, err
}

func (s *userService) ListEmail(ctx context.Context, opts scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/user/emails?%s", encodeListOptions(opts))
	out := []*email{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertEmailList(out), res, err
}

type user struct {
	ID       int         `json:"id"`
	Username string      `json:"username"`
	Name     string      `json:"name"`
	Email    null.String `json:"email"`
	Avatar   string      `json:"avatar_url"`
}

type email struct {
	Email     string      `json:"email"`
	Confirmed null.String `json:"confirmed_at"`
}

// helper function to convert from the gitlab user structure to
// the common user structure.
func convertUser(from *user) *scm.User {
	return &scm.User{
		Avatar: from.Avatar,
		Email:  from.Email.String,
		Login:  from.Username,
		Name:   from.Name,
	}
}

// helper function to convert from the gitlab email list to
// the common email structure.
func convertEmailList(from []*email) []*scm.Email {
	to := []*scm.Email{}
	for _, v := range from {
		to = append(to, convertEmail(v))
	}
	return to
}

// helper function to convert from the gitlab email structure to
// the common email structure.
func convertEmail(from *email) *scm.Email {
	return &scm.Email{
		Value:    from.Email,
		Verified: !from.Confirmed.IsZero(),
	}
}
