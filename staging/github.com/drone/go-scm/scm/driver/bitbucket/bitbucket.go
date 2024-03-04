// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bitbucket implements a Bitbucket Cloud client.
package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/drone/go-scm/scm"
)

// New returns a new Bitbucket API client.
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
	client.Driver = scm.DriverBitbucket
	client.Linker = &linker{"https://bitbucket.org/"}
	client.Contents = &contentService{client}
	client.Git = &gitService{client}
	client.Issues = &issueService{client}
	client.Milestones = &milestoneService{client}
	client.Organizations = &organizationService{client}
	client.PullRequests = &pullService{client}
	client.Repositories = &repositoryService{client}
	client.Releases = &releaseService{client}
	client.Reviews = &reviewService{client}
	client.Users = &userService{client}
	client.Webhooks = &webhookService{client}
	return client.Client, nil
}

// NewDefault returns a new Bitbucket API client using the
// default api.bitbucket.org address.
func NewDefault() *scm.Client {
	client, _ := New("https://api.bitbucket.org")
	return client
}

// wraper wraps the Client to provide high level helper functions
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
		// create or update content
		switch content := in.(type) {
		case *contentCreateUpdate:
			// add the content to the multipart
			myReader := strings.NewReader(string(content.Content))
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			var fw io.Writer
			fw, _ = w.CreateFormFile(content.Files, "")
			_, _ = io.Copy(fw, myReader)
			// add the other fields
			if content.Message != "" {
				_ = w.WriteField("message", content.Message)
			}
			if content.Branch != "" {
				_ = w.WriteField("branch", content.Branch)
			}
			if content.Sha != "" {
				_ = w.WriteField("parents", content.Sha)
			}
			if content.Author != "" {
				_ = w.WriteField("author", content.Author)
			}
			w.Close()
			// write the multipart response to the body
			req.Body = &b
			// write the content type that contains the length of the multipart
			req.Header = map[string][]string{
				"Content-Type": {w.FormDataContentType()},
			}
		case *contentDelete:
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			fw, err := writer.CreateFormField("files")
			_, err = io.Copy(fw, strings.NewReader(content.File))
			if err != nil {
				return nil, err
			}
			fw, err = writer.CreateFormField("message")
			_, err = io.Copy(fw, strings.NewReader(content.Message))
			if err != nil {
				return nil, err
			}
			fw, err = writer.CreateFormField("author")
			_, err = io.Copy(fw, strings.NewReader(content.Author))
			if err != nil {
				return nil, err
			}
			if content.Branch != "" {
				fw, err = writer.CreateFormField("branch")
				_, err = io.Copy(fw, strings.NewReader(content.Branch))
				if err != nil {
					return nil, err
				}
			}
			writer.Close()
			req.Body = bytes.NewReader(body.Bytes())
			req.Header = map[string][]string{
				"Content-Type": {writer.FormDataContentType()},
			}
		default:
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(in)
			req.Header = map[string][]string{
				"Content-Type": {"application/json"},
			}
			req.Body = buf
		}
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status == 401 {
		return res, scm.ErrNotAuthorized
	} else if res.Status > 300 {
		err := new(Error)
		json.NewDecoder(res.Body).Decode(err)
		return res, err
	}

	if out == nil {
		return res, nil
	}

	// if raw output is expected, copy to the provided
	// buffer and exit.
	if w, ok := out.(io.Writer); ok {
		io.Copy(w, res.Body)
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	return res, json.NewDecoder(res.Body).Decode(out)
}

// pagination represents Bitbucket pagination properties
// embedded in list responses.
type pagination struct {
	PageLen int    `json:"pagelen"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Next    string `json:"next"`
}

// link represents Bitbucket link properties embedded in responses.
type link struct {
	Href string `json:"href"`
}

// Error represents a Bitbucket error.
type Error struct {
	Type string `json:"type"`
	Data struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *Error) Error() string {
	return e.Data.Message
}
