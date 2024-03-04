package gogs

import (
	"context"
	"github.com/drone/go-scm/scm"
)

type releaseService struct {
	client *wrapper
}

func (s *releaseService) Find(ctx context.Context, repo string, id int) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) FindByTag(ctx context.Context, repo string, tag string) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) List(ctx context.Context, repo string, opts scm.ReleaseListOptions) ([]*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) Create(ctx context.Context, repo string, input *scm.ReleaseInput) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) Delete(ctx context.Context, repo string, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *releaseService) DeleteByTag(ctx context.Context, repo string, tag string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *releaseService) Update(ctx context.Context, repo string, id int, input *scm.ReleaseInput) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) UpdateByTag(ctx context.Context, repo string, tag string, input *scm.ReleaseInput) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}
