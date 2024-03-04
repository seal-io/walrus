// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	path := "plugins/servlet/applinks/whoami"
	out := new(bytes.Buffer)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	if err != nil {
		return nil, res, err
	}
	login := out.String()
	login = strings.TrimSpace(login)
	return s.FindLogin(ctx, login)
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := fmt.Sprintf("rest/api/1.0/users/%s", login)
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	if err == nil {
		return convertUser(out), res, err
	}

	// HACK: the below code is a hack to account for the
	// fact that the above API call requires the slug,
	// but "plugins/servlet/applinks/whoami" may not return
	// the slug. When this happens we need to search
	// and find the matching user.

	// HACK: only use the below special logic for usernames
	// that contain an @ symbol.
	if !strings.Contains(login, "@") {
		return nil, res, err
	}

	path = fmt.Sprintf("/rest/api/1.0/users?filter=%s", login)
	filter := new(userFilter)
	res, err = s.client.do(ctx, "GET", path, nil, filter)
	if err != nil {
		return nil, res, err
	}

	// iterate through the search results and find
	// the username that is an exact match.
	for _, item := range filter.Values {
		// must be an exact match
		if item.Name == login {
			return convertUser(item), res, err
		}
	}
	return nil, res, scm.ErrNotFound
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	user, res, err := s.Find(ctx)
	var email string
	if err == nil {
		email = user.Email
	}
	return email, res, err
}

func (s *userService) ListEmail(context.Context, scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type user struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	Links        struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

type userFilter struct {
	Values []*user `json:"values"`
}

func convertUser(from *user) *scm.User {
	return &scm.User{
		Avatar: avatarLink(from.EmailAddress),
		Login:  from.Slug,
		Name:   from.DisplayName,
		Email:  from.EmailAddress,
	}
}

func avatarLink(email string) string {
	hasher := md5.New()
	hasher.Write([]byte(strings.ToLower(email)))
	emailHash := fmt.Sprintf("%v", hex.EncodeToString(hasher.Sum(nil)))
	avatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%s.jpg", emailHash)
	return avatarURL
}
