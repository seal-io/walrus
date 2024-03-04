// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	urlEncodedRef := url.QueryEscape(ref)
	endpoint := fmt.Sprintf("/2.0/repositories/%s/src/%s/%s", repo, urlEncodedRef, path)
	out := new(bytes.Buffer)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	content := &scm.Content{
		Path: path,
		Data: out.Bytes(),
	}
	if err != nil {
		return content, res, err
	}
	metaEndpoint := fmt.Sprintf("/2.0/repositories/%s/src/%s/%s?format=meta", repo, urlEncodedRef, path)
	metaOut := new(metaContent)
	metaRes, metaErr := s.client.do(ctx, "GET", metaEndpoint, nil, metaOut)
	if metaErr == nil {
		content.Sha = metaOut.Commit.Hash
		return content, metaRes, metaErr
	} else {
		// do not risk that returning an error if getting the meta fails.
		return content, res, err
	}
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("/2.0/repositories/%s/src", repo)
	in := &contentCreateUpdate{
		Files:   path,
		Message: params.Message,
		Branch:  params.Branch,
		Content: params.Data,
		Author:  fmt.Sprintf("%s <%s>", params.Signature.Name, params.Signature.Email),
	}
	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	// https://jira.atlassian.com/browse/BCLOUD-20424?error=login_required&error_description=Login+required&state=196d85f7-a181-4b63-babe-0b567858d8f5 ugh :(
	endpoint := fmt.Sprintf("/2.0/repositories/%s/src", repo)
	in := &contentCreateUpdate{
		Files:   path,
		Message: params.Message,
		Branch:  params.Branch,
		Content: params.Data,
		Sha:     params.Sha,
		Author:  fmt.Sprintf("%s <%s>", params.Signature.Name, params.Signature.Email),
	}
	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) Delete(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	author := fmt.Sprintf("%s <%s>", params.Signature.Name, params.Signature.Email)
	endpoint := fmt.Sprintf("/2.0/repositories/%s/src", repo)
	in := &contentDelete{
		File:    path,
		Branch:  params.Branch,
		Message: params.Message,
		Sha:     params.Sha,
		Author:  author,
	}
	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, opts scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	endpoint := fmt.Sprintf("/2.0/repositories/%s/src/%s/%s?%s", repo, ref, path, encodeListOptions(opts))
	if opts.URL != "" {
		endpoint = opts.URL
	}

	out := new(contents)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	copyPagination(out.pagination, res)
	return convertContentInfoList(out), res, err
}

type contents struct {
	pagination
	Values []*content `json:"values"`
}

type content struct {
	Path   string `json:"path"`
	Type   string `json:"type"`
	Commit struct {
		Hash string `json:"hash"`
	} `json:"commit"`
	Attributes []string `json:"attributes"`
}

type metaContent struct {
	Path   string `json:"path"`
	Commit struct {
		Hash string `json:"hash"`
	} `json:"commit"`
}

type contentCreateUpdate struct {
	Files   string `json:"files"`
	Branch  string `json:"branch"`
	Message string `json:"message"`
	Content []byte `json:"content"`
	Sha     string `json:"sha"`
	Author  string `json:"author"`
}

type contentDelete struct {
	File    string `json:"file"`
	Branch  string `json:"branch"`
	Message string `json:"message"`
	Sha     string `json:"sha"`
	Author  string `json:"author"`
}

func convertContentInfoList(from *contents) []*scm.ContentInfo {
	to := []*scm.ContentInfo{}
	for _, v := range from.Values {
		to = append(to, convertContentInfo(v))
	}
	return to
}

func convertContentInfo(from *content) *scm.ContentInfo {
	to := &scm.ContentInfo{
		Path: from.Path,
		Sha:  from.Commit.Hash,
	}
	switch from.Type {
	case "commit_file":
		to.Kind = func() scm.ContentKind {
			for _, attr := range from.Attributes {
				switch attr {
				case "link":
					return scm.ContentKindSymlink
				case "subrepository":
					return scm.ContentKindGitlink
				}
			}
			return scm.ContentKindFile
		}()
	case "commit_directory":
		to.Kind = scm.ContentKindDirectory
	default:
		to.Kind = scm.ContentKindUnsupported
	}
	return to
}
