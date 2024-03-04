// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
)

func TestWebhooks(t *testing.T) {
	tests := []struct {
		event  string
		before string
		after  string
		obj    interface{}
	}{
		//
		// branch events
		//
		// push branch create
		{
			event:  "branch_created",
			before: "testdata/webhooks/branch_create.json",
			after:  "testdata/webhooks/branch_create.json.golden",
			obj:    new(scm.PushHook),
		},
		// push branch update
		{
			event:  "branch_updated",
			before: "testdata/webhooks/branch_updated.json",
			after:  "testdata/webhooks/branch_updated.json.golden",
			obj:    new(scm.PushHook),
		},
		//
		// pull request events
		//
		// pull request opened
		{
			event:  "pullreq_created",
			before: "testdata/webhooks/pull_request_opened.json",
			after:  "testdata/webhooks/pull_request_opened.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request reopened
		{
			event:  "pullreq_reopened",
			before: "testdata/webhooks/pull_request_reopened.json",
			after:  "testdata/webhooks/pull_request_reopened.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request branch updated
		{
			event:  "pullreq_branch_updated",
			before: "testdata/webhooks/pull_request_branch_updated.json",
			after:  "testdata/webhooks/pull_request_branch_updated.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request comment created
		{
			event:  "pullreq_comment_created",
			before: "testdata/webhooks/pull_request_comment_created.json",
			after:  "testdata/webhooks/pull_request_comment_created.json.golden",
			obj:    new(scm.PullRequestCommentHook),
		},
		// pull request closed
		{
			event:  "pullreq_reopened",
			before: "testdata/webhooks/pull_request_closed.json",
			after:  "testdata/webhooks/pull_request_closed.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request merged
		{
			event:  "pullreq_reopened",
			before: "testdata/webhooks/pull_request_merged.json",
			after:  "testdata/webhooks/pull_request_merged.json.golden",
			obj:    new(scm.PullRequestHook),
		},
	}

	for _, test := range tests {
		before, err := ioutil.ReadFile(test.before)
		if err != nil {
			t.Error(err)
			continue
		}
		after, err := ioutil.ReadFile(test.after)
		if err != nil {
			t.Error(err)
			continue
		}

		buf := bytes.NewBuffer(before)
		r, _ := http.NewRequest("GET", "/", buf)
		r.Header.Set("X-Harness-Trigger", test.event)

		s := new(webhookService)
		o, err := s.Parse(r, secretFunc)
		if err != nil && err != scm.ErrSignatureInvalid {
			t.Error(err)
			continue
		}

		err = json.Unmarshal(after, test.obj)
		if err != nil {
			t.Error(err)
			continue
		}

		if diff := cmp.Diff(test.obj, o); diff != "" {
			t.Errorf("Error unmarshaling %s", test.before)
			t.Log(diff)

			// debug only. remove once implemented
			_ = json.NewEncoder(os.Stdout).Encode(o)

		}

		// switch event := o.(type) {
		// case *scm.PushHook:
		// 	if !strings.HasPrefix(event.Ref, "refs/") {
		// 		t.Errorf("Push hook reference must start with refs/")
		// 	}
		// case *scm.BranchHook:
		// 	if strings.HasPrefix(event.Ref.Name, "refs/") {
		// 		t.Errorf("Branch hook reference must not start with refs/")
		// 	}
		// case *scm.TagHook:
		// 	if strings.HasPrefix(event.Ref.Name, "refs/") {
		// 		t.Errorf("Branch hook reference must not start with refs/")
		// 	}
		// }
	}
}
func secretFunc(scm.Webhook) (string, error) {
	return "topsecret", nil
}
