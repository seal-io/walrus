// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

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
		event  string
		before string
		after  string
		obj    interface{}
	}{
		//
		// Push Hook events
		//

		// Push Hook
		{
			event:  "Push Hook",
			before: "testdata/webhooks/push.json",
			after:  "testdata/webhooks/push.json.golden",
			obj:    new(scm.PushHook),
		},
		// Push Hook create
		{
			event:  "Push Hook",
			before: "testdata/webhooks/push_branch_create.json",
			after:  "testdata/webhooks/push_branch_create.json.golden",
			obj:    new(scm.PushHook),
		},
		// Push Hook delete
		{
			event:  "Push Hook",
			before: "testdata/webhooks/push_branch_delete.json",
			after:  "testdata/webhooks/push_branch_delete.json.golden",
			obj:    new(scm.PushHook),
		},

		//
		// Push Tag Hook events
		//

		// Push Tag Hook create
		{
			event:  "Tag Push Hook",
			before: "testdata/webhooks/tag_create.json",
			after:  "testdata/webhooks/tag_create.json.golden",
			obj:    new(scm.TagHook),
		},
		// Push Tag Hook delete
		{
			event:  "Tag Push Hook",
			before: "testdata/webhooks/tag_delete.json",
			after:  "testdata/webhooks/tag_delete.json.golden",
			obj:    new(scm.TagHook),
		},

		//
		// Merge Request Hook events
		//

		// Merge Request Hook merge
		{
			event:  "Merge Request Hook",
			before: "testdata/webhooks/pr_merge.json",
			after:  "testdata/webhooks/pr_merge.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// Merge Request Hook open
		{
			event:  "Merge Request Hook",
			before: "testdata/webhooks/pr_open.json",
			after:  "testdata/webhooks/pr_open.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// Merge Request Hook close
		{
			event:  "Merge Request Hook",
			before: "testdata/webhooks/pr_close.json",
			after:  "testdata/webhooks/pr_close.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request update
		{
			event:  "Merge Request Hook",
			before: "testdata/webhooks/pr_update.json",
			after:  "testdata/webhooks/pr_update.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// Merge Request Hook labeled
		{
			event:  "Merge Request Hook",
			before: "testdata/webhooks/pr_labeled.json",
			after:  "testdata/webhooks/pr_labeled.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// Merge Request Hook unlabeled
		{
			event:  "Merge Request Hook",
			before: "testdata/webhooks/pr_unlabeled.json",
			after:  "testdata/webhooks/pr_unlabeled.json.golden",
			obj:    new(scm.PullRequestHook),
		},

		//
		// Issue Hook events
		//

		// Issue Hook open
		{
			event:  "Issue Hook",
			before: "testdata/webhooks/issue_hook_open.json",
			after:  "testdata/webhooks/issue_hook_open.json.golden",
			obj:    new(scm.IssueHook),
		},
		// Issue Hook delete
		{
			event:  "Issue Hook",
			before: "testdata/webhooks/issue_hook_delete.json",
			after:  "testdata/webhooks/issue_hook_delete.json.golden",
			obj:    new(scm.IssueHook),
		},
		// Issue Hook state_change
		{
			event:  "Issue Hook",
			before: "testdata/webhooks/issue_hook_state_change.json",
			after:  "testdata/webhooks/issue_hook_state_change.json.golden",
			obj:    new(scm.IssueHook),
		},
		// Issue Hook assign
		{
			event:  "Issue Hook",
			before: "testdata/webhooks/issue_hook_assign.json",
			after:  "testdata/webhooks/issue_hook_assign.json.golden",
			obj:    new(scm.IssueHook),
		},

		//
		// Note Hook events
		//

		// Note Hook issue comment
		{
			event:  "Note Hook",
			before: "testdata/webhooks/note_hook_issue_comment.json",
			after:  "testdata/webhooks/note_hook_issue_comment.json.golden",
			obj:    new(scm.IssueCommentHook),
		},
		// Note Hook pull request comment
		{
			event:  "Note Hook",
			before: "testdata/webhooks/note_hook_pr_comment.json",
			after:  "testdata/webhooks/note_hook_pr_comment.json.golden",
			obj:    new(scm.PullRequestCommentHook),
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
		r.Header.Set("X-Gitee-Event", test.event)
		r.Header.Set("X-Gitee-Token", "Xvh4YPVe6l31XpDRL9J2yeaEXabsckIoUUschpXiVck=")
		r.Header.Set("X-Gitee-Timestamp", "1633679083918")
		r.Header.Set("User-Agent", "git-oschina-hook")

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
		}

		switch event := o.(type) {
		case *scm.PushHook:
			if !strings.HasPrefix(event.Ref, "refs/") {
				t.Errorf("Push hook reference must start with refs/")
			}
		case *scm.TagHook:
			if strings.HasPrefix(event.Ref.Name, "refs/") {
				t.Errorf("Branch hook reference must not start with refs/")
			}
		}
	}
}

func TestWebhook_ErrUnknownEvent(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/webhooks/push.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != scm.ErrUnknownEvent {
		t.Errorf("Expect unknown event error, got %v", err)
	}
}

func TestWebhookInvalid(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/webhooks/push.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
	r.Header.Set("X-Gitee-Event", "Push Hook")
	r.Header.Set("X-Gitee-Token", "Xvh4YPVe6l31XpDRL9J2yeaEXabsckIoUUschpXiVck=")
	r.Header.Set("X-Gitee-Timestamp", "1633679083917")
	r.Header.Set("User-Agent", "git-oschina-hook")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != scm.ErrSignatureInvalid {
		t.Errorf("Expect invalid signature error, got %v", err)
	}
}

func TestWebhookValid(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/webhooks/push.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
	r.Header.Set("X-Gitee-Event", "Push Hook")
	r.Header.Set("X-Gitee-Token", "Xvh4YPVe6l31XpDRL9J2yeaEXabsckIoUUschpXiVck=")
	r.Header.Set("X-Gitee-Timestamp", "1633679083918")
	r.Header.Set("User-Agent", "git-oschina-hook")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != nil {
		t.Errorf("Expect valid signature, got %v", err)
	}
}

func secretFunc(scm.Webhook) (string, error) {
	return "bBg5lrt03VixkX85CNqYIcecC0SIGASE", nil
}
