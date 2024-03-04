// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gitee implements a Gitee client.
package gitee

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/drone/go-scm/scm"
)

// New returns a new Gitee API client.
func New(uri string) (*scm.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	client := &wrapper{new(scm.Client)}
	client.BaseURL = base
	// initialize services
	client.Driver = scm.DriverGitee
	client.Linker = &linker{websiteAddress(base)}
	client.Contents = &contentService{client}
	client.Git = &gitService{client}
	client.Issues = &issueService{client}
	client.Organizations = &organizationService{client}
	client.PullRequests = &pullService{client}
	client.Repositories = &RepositoryService{client}
	client.Reviews = &reviewService{client}
	client.Users = &userService{client}
	client.Webhooks = &webhookService{client}
	return client.Client, nil
}

// NewDefault returns a new Gitee API client using the
// default gitee.com/api/v5 address.
func NewDefault() *scm.Client {
	client, _ := New("https://gitee.com/api/v5")
	return client
}

// wrapper wraps the Client to provide high level helper functions
// for making http requests and unmarshaling the response.
type wrapper struct {
	*scm.Client
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*scm.Response, error) {
	req := &scm.Request{
		Method: method,
		Path:   path,
	}
	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(in)
		req.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		req.Body = buf
	}
	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// parse the gitee request id.
	res.ID = res.Header.Get("X-Request-Id")

	// gitee pageValues
	populatePageValues(req, res)

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status > 300 {
		err := new(Error)
		json.NewDecoder(res.Body).Decode(err)
		return res, err
	}

	if out == nil {
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	return res, json.NewDecoder(res.Body).Decode(out)
}

// Error represents a Gitee error.
type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// helper function converts the Gitee API url to
// the website url.
func websiteAddress(u *url.URL) string {
	host, proto := u.Host, u.Scheme
	switch host {
	case "gitee.com/api/v5":
		return "https://gitee.com/"
	}
	return proto + "://" + host + "/"
}

// populatePageValues parses the HTTP Link response headers
// and populates the various pagination link values in the
// Response.
// response header: total_page, total_count
func populatePageValues(req *scm.Request, resp *scm.Response) {
	// get last
	last, totalError := strconv.Atoi(resp.Header.Get("total_page"))
	if totalError != nil {
		return
	}
	// get curren page
	reqURL, err := url.Parse(req.Path)
	if err != nil {
		return
	}
	currentPageStr := reqURL.Query().Get("page")
	var current int
	if currentPageStr == "" {
		current = 1
	} else {
		currentPage, currentError := strconv.Atoi(currentPageStr)
		if currentError != nil {
			return
		}
		current = currentPage
	}

	// first, prev
	if current <= 1 {
		resp.Page.First = 0
		resp.Page.Prev = 0
	} else {
		resp.Page.First = 1
		resp.Page.Prev = current - 1
	}
	// last, next
	if current >= last {
		resp.Page.Last = 0
		resp.Page.Next = 0
	} else {
		resp.Page.Last = last
		resp.Page.Next = current + 1
	}
}
