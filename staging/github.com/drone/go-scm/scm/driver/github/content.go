// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	urlEncodedRef := url.QueryEscape(ref)
	endpoint := fmt.Sprintf("repos/%s/contents/%s?ref=%s", repo, path, urlEncodedRef)
	out := new(content)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	raw, _ := base64.StdEncoding.DecodeString(out.Content)
	return &scm.Content{
		Path: out.Path,
		Data: raw,
		// NB the sha returned for github rest api is the blob sha, not the commit sha
		BlobID: out.Sha,
	}, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s", repo, path)
	in := &contentCreateUpdate{
		Message: params.Message,
		Branch:  params.Branch,
		Content: params.Data,
		Committer: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
		Author: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
	}

	res, err := s.client.do(ctx, "PUT", endpoint, in, nil)
	return res, err
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s", repo, path)
	in := &contentCreateUpdate{
		Message: params.Message,
		Branch:  params.Branch,
		Content: params.Data,
		// NB the sha passed to github rest api is the blob sha, not the commit sha
		Sha: params.BlobID,
		Committer: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
		Author: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
	}
	res, err := s.client.do(ctx, "PUT", endpoint, in, nil)
	return res, err
}

func (s *contentService) Delete(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s", repo, path)
	in := &contentCreateUpdate{
		Message: params.Message,
		Branch:  params.Branch,
		// NB the sha passed to github rest api is the blob sha, not the commit sha
		Sha: params.BlobID,
		Committer: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
		Author: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
	}
	res, err := s.client.do(ctx, "DELETE", endpoint, in, nil)
	return res, err
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, _ scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s?ref=%s", repo, path, ref)
	out := []*content{}
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertContentInfoList(out), res, err
}

type content struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Sha     string `json:"sha"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type contentCreateUpdate struct {
	Branch    string       `json:"branch"`
	Message   string       `json:"message"`
	Content   []byte       `json:"content"`
	Sha       string       `json:"sha"`
	Author    commitAuthor `json:"author"`
	Committer commitAuthor `json:"committer"`
}

type commitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func convertContentInfoList(from []*content) []*scm.ContentInfo {
	to := []*scm.ContentInfo{}
	for _, v := range from {
		to = append(to, convertContentInfo(v))
	}
	return to
}

func convertContentInfo(from *content) *scm.ContentInfo {
	to := &scm.ContentInfo{
		Path:   from.Path,
		BlobID: from.Sha}
	switch from.Type {
	case "file":
		to.Kind = scm.ContentKindFile
	case "dir":
		to.Kind = scm.ContentKindDirectory
	case "symlink":
		to.Kind = scm.ContentKindSymlink
	case "submodule":
		to.Kind = scm.ContentKindGitlink
	default:
		to.Kind = scm.ContentKindUnsupported
	}
	return to
}
