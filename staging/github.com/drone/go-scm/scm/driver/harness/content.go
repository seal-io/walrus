// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	slug := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(slug)
	if err != nil {
		return nil, nil, err
	}
	endpoint := fmt.Sprintf("api/v1/repos/%s/content/%s?git_ref=%s&include_commit=true&%s", repoId, path, ref, queryParams)
	out := new(fileContent)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	// decode raw output content
	raw, _ := base64.StdEncoding.DecodeString(out.Content.Data)
	return &scm.Content{
		Path:   path,
		Sha:    out.LatestCommit.Sha,
		BlobID: out.Sha,
		Data:   raw,
	}, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	slug := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(slug)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("api/v1/repos/%s/commits?%s", repoId, queryParams)
	a := action{
		Action:   "CREATE",
		Path:     path,
		Payload:  string(params.Data),
		Encoding: "string",
	}
	in := editFile{
		Branch:  params.Branch,
		Message: params.Message,
		Title:   params.Message,
		Actions: []action{a},
	}

	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	slug := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(slug)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("api/v1/repos/%s/commits?%s", repoId, queryParams)
	a := action{
		Action:   "UPDATE",
		Path:     path,
		Payload:  string(params.Data),
		Encoding: "string",
		Sha:      params.BlobID,
	}
	in := editFile{
		Branch:  params.Branch,
		Message: params.Message,
		Title:   params.Message,
		Actions: []action{a},
	}

	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) Delete(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	slug := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(slug)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("api/v1/repos/%s/commits?%s", repoId, queryParams)
	a := action{
		Action:   "DELETE",
		Path:     path,
		Encoding: "string",
	}
	in := editFile{
		Branch:  params.Branch,
		Message: params.Message,
		Title:   params.Message,
		Actions: []action{a},
	}

	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, _ scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	slug := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(slug)
	if err != nil {
		return nil, nil, err
	}
	endpoint := fmt.Sprintf("api/v1/repos/%s/content/%s?git_ref=%s&include_commit=true&%s", repoId, path, ref, queryParams)
	out := new(contentList)
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertContentInfoList(out.Content.Entries), res, err
}

type editFile struct {
	Actions   []action `json:"actions"`
	Branch    string   `json:"branch"`
	Message   string   `json:"message"`
	NewBranch string   `json:"new_branch"`
	Title     string   `json:"title"`
}

type action struct {
	Action   string `json:"action"`
	Encoding string `json:"encoding"`
	Path     string `json:"path"`
	Payload  string `json:"payload"`
	Sha      string `json:"sha"`
}

type fileContent struct {
	Type         string `json:"type"`
	Sha          string `json:"sha"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	LatestCommit struct {
		Sha     string `json:"sha"`
		Title   string `json:"title"`
		Message string `json:"message"`
		Author  struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"committer"`
	} `json:"latest_commit"`
	Content struct {
		Encoding string `json:"encoding"`
		Data     string `json:"data"`
		Size     int    `json:"size"`
	} `json:"content"`
}

type contentList struct {
	Type         string `json:"type"`
	Sha          string `json:"sha"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	LatestCommit struct {
		Sha     string `json:"sha"`
		Title   string `json:"title"`
		Message string `json:"message"`
		Author  struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"committer"`
	} `json:"latest_commit"`
	Content struct {
		Entries []fileEntry `json:"entries"`
	} `json:"content"`
}

type fileEntry struct {
	Type         string `json:"type"`
	Sha          string `json:"sha"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	LatestCommit struct {
		Sha     string `json:"sha"`
		Title   string `json:"title"`
		Message string `json:"message"`
		Author  struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"committer"`
	} `json:"latest_commit"`
}

func convertContentInfoList(from []fileEntry) []*scm.ContentInfo {
	to := []*scm.ContentInfo{}
	for _, v := range from {
		to = append(to, convertContentInfo(v))
	}
	return to
}

func convertContentInfo(from fileEntry) *scm.ContentInfo {
	to := &scm.ContentInfo{
		Path:   from.Path,
		Sha:    from.LatestCommit.Sha,
		BlobID: from.Sha,
	}
	switch from.Type {
	case "file":
		to.Kind = scm.ContentKindFile
	case "dir":
		to.Kind = scm.ContentKindDirectory
	default:
		to.Kind = scm.ContentKindUnsupported
	}
	return to
}
