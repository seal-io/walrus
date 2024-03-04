// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package azure implements a azure client.
package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/drone/go-scm/scm"
)

// New returns a new azure API client.
func New(uri, owner, project string) (*scm.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	if owner == "" {
		return nil, fmt.Errorf("azure owner is required")
	}
	client := &wrapper{
		new(scm.Client),
		owner,
		project,
	}
	client.BaseURL = base
	// initialize services
	client.Driver = scm.DriverAzure
	client.Linker = &linker{base.String()}
	client.Contents = &contentService{client}
	client.Git = &gitService{client}
	client.Issues = &issueService{client}
	client.Organizations = &organizationService{client}
	client.PullRequests = &pullService{&issueService{client}}
	client.Repositories = &RepositoryService{client}
	client.Reviews = &reviewService{client}
	client.Users = &userService{client}
	client.Webhooks = &webhookService{client}
	return client.Client, nil
}

// NewDefault returns a new azure API client.
func NewDefault(owner, project string) *scm.Client {
	client, _ := New("https://dev.azure.com", owner, project)
	return client
}

// wrapper wraps the Client to provide high level helper functions for making http requests and unmarshaling the response.
type wrapper struct {
	*scm.Client
	owner   string
	project string
}

// do wraps the Client.Do function by creating the Request and unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*scm.Response, error) {
	req := &scm.Request{
		Method: method,
		Path:   path,
	}
	// if we are posting or putting data, we need to write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(in)
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

	// error response.
	if res.Status > 300 {
		err := new(Error)
		_ = json.NewDecoder(res.Body).Decode(err)
		return res, err
	}
	// the following is used for debugging purposes.
	// bytes, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(bytes))

	if out == nil {
		return res, nil
	}
	// if a json response is expected, parse and return the json response.
	decodeErr := json.NewDecoder(res.Body).Decode(out)
	// following line is used for debugging purposes.
	//_ = json.NewEncoder(os.Stdout).Encode(out)
	return res, decodeErr
}

// Error represents am Azure error.
type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func ProjectRequiredError() error {
	return errors.New("This API endpoint requires a project to be specified")
}

func SanitizeBranchName(name string) string {
	if strings.Contains(name, "/") {
		return name
	}
	return "refs/heads/" + name
}
