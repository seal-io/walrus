// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testContents(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testContentFind(client))
		t.Run("Find/Branch", testContentFindBranch(client))
	}
}

func testContentFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Contents.Find(context.Background(), "kit101/drone-yml-test", "main.py", "5e7876efb3468ff679410b82a72f7c002382d41e")
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := result.Path, "main.py"; got != want {
			t.Errorf("Got file path %q, want %q", got, want)
		}
		if got, want := string(result.Data), "if __name__ == '__main__':\r\n    print('Hello world.')"; got != want {
			t.Errorf("Got file data %q, want %q", got, want)
		}
	}
}

func testContentFindBranch(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Contents.Find(context.Background(), "kit101/drone-yml-test", "main.py", "feat-4")
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := result.Path, "main.py"; got != want {
			t.Errorf("Got file path %q, want %q", got, want)
		}
		if got, want := string(result.Data), "if __name__ == '__main__':\r\n    print('Hello world.')"; got != want {
			t.Errorf("Got file data %q, want %q", got, want)
		}
	}
}
