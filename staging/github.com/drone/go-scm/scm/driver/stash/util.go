// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/drone/go-scm/scm"
)

const (
	defaultLimit = 25
)

func encodeListOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page > 1 {
		params.Set("start", strconv.Itoa(
			(opts.Page-1)*opts.Size),
		)
	}
	if opts.Size != 0 {
		params.Set("limit", strconv.Itoa(opts.Size))
	}
	return params.Encode()
}

func encodeListOptionsV2(opts scm.ListOptions) string {
	params := url.Values{}
	limit := defaultLimit
	if opts.Size != 0 {
		limit = opts.Size
	}
	params.Set("limit", strconv.Itoa(limit))

	if opts.Page > 0 {
		params.Set("start", strconv.Itoa(
			(opts.Page-1)*limit),
		)
	}
	return params.Encode()
}

func encodeBranchListOptions(opts scm.BranchListOptions) string {
	params := url.Values{}
	if opts.SearchTerm != "" {
		params.Set("filterText", opts.SearchTerm)
	}
	if opts.PageListOptions != (scm.ListOptions{}) {
		if opts.PageListOptions.Page > 1 {
			params.Set("start", strconv.Itoa(
				(opts.PageListOptions.Page-1)*opts.PageListOptions.Size),
			)
		}
		if opts.PageListOptions.Size != 0 {
			params.Set("limit", strconv.Itoa(opts.PageListOptions.Size))
		}
	}
	return params.Encode()
}

func encodeListRoleOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page > 1 {
		params.Set("start", strconv.Itoa(
			(opts.Page-1)*opts.Size),
		)
	}
	if opts.Size != 0 {
		params.Set("limit", strconv.Itoa(opts.Size))
	}
	params.Set("permission", "REPO_READ")
	return params.Encode()
}

func encodeRepoListOptions(opts scm.RepoListOptions) string {
	params := url.Values{}
	if opts.RepoSearchTerm != (scm.RepoSearchTerm{}) {
		if opts.RepoSearchTerm.RepoName != "" {
			params.Set("name", opts.RepoSearchTerm.RepoName)
		}
	}
	if opts.ListOptions != (scm.ListOptions{}) {
		if opts.ListOptions.Page > 1 {
			params.Set("start", strconv.Itoa(
				(opts.ListOptions.Page-1)*opts.ListOptions.Size),
			)
		}
		if opts.ListOptions.Size != 0 {
			params.Set("limit", strconv.Itoa(opts.ListOptions.Size))
		}
	}
	params.Set("permission", "REPO_READ")
	return params.Encode()
}

func encodePullRequestListOptions(opts scm.PullRequestListOptions) string {
	params := url.Values{}
	if opts.Page > 1 {
		params.Set("start", strconv.Itoa(
			(opts.Page-1)*opts.Size),
		)
	}
	if opts.Size != 0 {
		params.Set("limit", strconv.Itoa(opts.Size))
	}
	if opts.Open && opts.Closed {
		params.Set("state", "all")
	} else if opts.Closed {
		params.Set("state", "closed")
	}
	return params.Encode()
}

func copyPagination(from pagination, to *scm.Response) error {
	if to == nil {
		return nil
	}
	to.Page.First = 1
	if from.LastPage.Bool {
		return nil
	}
	if from.Limit.Int64 == 0 {
		return errors.New("Unknown page limit")
	}
	to.Page.Next = int(from.NextPage.Int64/from.Limit.Int64 + 1)
	return nil
}
