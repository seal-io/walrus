// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/drone/go-scm/scm/driver/bitbucket"
	"github.com/drone/go-scm/scm/transport"
)

var noContext = context.Background()

func TestIntegration(t *testing.T) {
	client := bitbucket.NewDefault()
	client.Client = &http.Client{
		Transport: &transport.BearerToken{
			Token: os.Getenv("BITBUCKET_TOKEN"),
		},
	}
}
