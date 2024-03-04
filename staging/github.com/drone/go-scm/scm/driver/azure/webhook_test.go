// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
)

func TestWebhooks(t *testing.T) {
	tests := []struct {
		before string
		after  string
		obj    interface{}
	}{
		// push hook
		{
			before: "testdata/webhooks/push.json",
			after:  "testdata/webhooks/push.json.golden",
			obj:    new(scm.PushHook),
		},
		// pull request events
		// pull request created
		{
			before: "testdata/webhooks/pr_created.json",
			after:  "testdata/webhooks/pr_created.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request updated
		{
			before: "testdata/webhooks/pr_updated.json",
			after:  "testdata/webhooks/pr_updated.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request merged
		{
			before: "testdata/webhooks/pr_merged.json",
			after:  "testdata/webhooks/pr_merged.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// issue comment create
		{
			before: "testdata/webhooks/issue_comment.json",
			after:  "testdata/webhooks/issue_comment.json.golden",
			obj:    new(scm.IssueCommentHook),
		},
		// issue comment edit
		{
			before: "testdata/webhooks/issue_comment_edit.json",
			after:  "testdata/webhooks/issue_comment_edit.json.golden",
			obj:    new(scm.IssueCommentHook),
		},
		// issue comment delete
		{
			before: "testdata/webhooks/issue_comment_delete.json",
			after:  "testdata/webhooks/issue_comment_delete.json.golden",
			obj:    new(scm.IssueCommentHook),
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
		r, _ := http.NewRequest("GET", "/?secret=71295b197fa25f4356d2fb9965df3f2379d903d7", buf)

		s := new(webhookService)
		o, err := s.Parse(r, secretFunc)
		if err != nil {
			t.Error(err)
			continue
		}

		err = json.Unmarshal(after, &test.obj)
		if err != nil {
			t.Error(err)
			continue
		}

		if diff := cmp.Diff(test.obj, o); diff != "" {
			t.Errorf("Error unmarshaling %s", test.before)
			t.Log(diff)

			// debug only. remove once implemented
			//	_ = json.NewEncoder(os.Stdout).Encode(o)
		}

		switch event := o.(type) {
		case *scm.PushHook:
			if !strings.HasPrefix(event.Ref, "refs/") {
				t.Errorf("Push hook reference must start with refs/")
			}
		case *scm.BranchHook:
			if strings.HasPrefix(event.Ref.Name, "refs/") {
				t.Errorf("Branch hook reference must not start with refs/")
			}
		case *scm.TagHook:
			if strings.HasPrefix(event.Ref.Name, "refs/") {
				t.Errorf("Branch hook reference must not start with refs/")
			}
		}
	}
}

func secretFunc(scm.Webhook) (string, error) {
	return "71295b197fa25f4356d2fb9965df3f2379d903d7", nil
}
