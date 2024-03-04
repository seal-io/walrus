// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type linker struct {
	base string
}

// Resource returns a link to the resource.
func (l *linker) Resource(ctx context.Context, repo string, ref scm.Reference) (string, error) {
	return "", scm.ErrNotSupported

}

// Diff returns a link to the diff.
func (l *linker) Diff(ctx context.Context, repo string, source, target scm.Reference) (string, error) {
	return "", scm.ErrNotSupported

}
