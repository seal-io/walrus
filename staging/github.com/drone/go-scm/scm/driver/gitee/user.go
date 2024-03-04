// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(user)
	res, err := s.client.do(ctx, "GET", "user", nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	user, res, err := s.Find(ctx)
	if err != nil {
		return "", nil, err
	}
	return user.Email, res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := fmt.Sprintf("users/%s", login)
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertUser(out), res, err
}

func (s *userService) ListEmail(context.Context, scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type user struct {
	ID                int       `json:"id"`
	Login             string    `json:"login"`
	Name              string    `json:"name"`
	AvatarURL         string    `json:"avatar_url"`
	URL               string    `json:"url"`
	HtmlURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	Blog              string    `json:"blog"`
	Weibo             string    `json:"weibo"`
	Bio               string    `json:"bio"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	Stared            int       `json:"stared"`
	Watched           int       `json:"watched"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Email             string    `json:"email"`
	Remark            string    `json:"remark"`
}

func convertUser(from *user) *scm.User {
	return &scm.User{
		Avatar:  from.AvatarURL,
		Email:   from.Email,
		Login:   from.Login,
		Name:    from.Name,
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
	}
}
