// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type reviewService struct {
	client *wrapper
}

func (s *reviewService) Find(context.Context, string, int, int) (*scm.Review, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *reviewService) List(context.Context, string, int, scm.ListOptions) ([]*scm.Review, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *reviewService) Create(context.Context, string, int, *scm.ReviewInput) (*scm.Review, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *reviewService) Delete(context.Context, string, int, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}
