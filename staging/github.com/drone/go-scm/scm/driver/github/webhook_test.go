// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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
		// push events
		//
		// push hooks
		{
			event:  "push",
			before: "testdata/webhooks/push.json",
			after:  "testdata/webhooks/push.json.golden",
			obj:    new(scm.PushHook),
		},
		// push tag create hooks
		{
			event:  "push",
			before: "testdata/webhooks/push_tag.json",
			after:  "testdata/webhooks/push_tag.json.golden",
			obj:    new(scm.PushHook),
		},
		// push tag delete hooks
		{
			event:  "push",
			before: "testdata/webhooks/push_tag_delete.json",
			after:  "testdata/webhooks/push_tag_delete.json.golden",
			obj:    new(scm.PushHook),
		},
		// push branch create
		{
			event:  "push",
			before: "testdata/webhooks/push_branch_create.json",
			after:  "testdata/webhooks/push_branch_create.json.golden",
			obj:    new(scm.PushHook),
		},
		// push branch delete
		{
			event:  "push",
			before: "testdata/webhooks/push_branch_delete.json",
			after:  "testdata/webhooks/push_branch_delete.json.golden",
			obj:    new(scm.PushHook),
		},

		//
		// branch events
		//

		// push branch create
		{
			event:  "create",
			before: "testdata/webhooks/branch_create.json",
			after:  "testdata/webhooks/branch_create.json.golden",
			obj:    new(scm.BranchHook),
		},
		// push branch delete
		{
			event:  "delete",
			before: "testdata/webhooks/branch_delete.json",
			after:  "testdata/webhooks/branch_delete.json.golden",
			obj:    new(scm.BranchHook),
		},

		//
		// comment events
		//

		// issue_comment
		{
			event:  "issue_comment",
			before: "testdata/webhooks/comment.json",
			after:  "testdata/webhooks/comment.json.golden",
			obj:    new(scm.IssueCommentHook),
		},

		//
		// tag events
		//

		// push tag create
		{
			event:  "create",
			before: "testdata/webhooks/tag_create.json",
			after:  "testdata/webhooks/tag_create.json.golden",
			obj:    new(scm.TagHook),
		},
		// push tag delete
		{
			event:  "delete",
			before: "testdata/webhooks/tag_delete.json",
			after:  "testdata/webhooks/tag_delete.json.golden",
			obj:    new(scm.TagHook),
		},

		//
		// pull request events
		//

		// pull request synced
		{
			event:  "pull_request",
			before: "testdata/webhooks/pr_sync.json",
			after:  "testdata/webhooks/pr_sync.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request ready for review
		{
			event:  "pull_request",
			before: "testdata/webhooks/pr_ready_for_review.json",
			after:  "testdata/webhooks/pr_ready_for_review.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request opened
		{
			event:  "pull_request",
			before: "testdata/webhooks/pr_opened.json",
			after:  "testdata/webhooks/pr_opened.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request closed
		{
			event:  "pull_request",
			before: "testdata/webhooks/pr_closed.json",
			after:  "testdata/webhooks/pr_closed.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request reopened
		{
			event:  "pull_request",
			before: "testdata/webhooks/pr_reopened.json",
			after:  "testdata/webhooks/pr_reopened.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request edited
		{
			event:  "pull_request",
			before: "testdata/webhooks/pr_edited.json",
			after:  "testdata/webhooks/pr_edited.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request labeled
		{
			event:  "pull_request",
			before: "testdata/webhooks/pr_labeled.json",
			after:  "testdata/webhooks/pr_labeled.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request unlabeled
		{
			event:  "pull_request",
			before: "testdata/webhooks/pr_unlabeled.json",
			after:  "testdata/webhooks/pr_unlabeled.json.golden",
			obj:    new(scm.PullRequestHook),
		},

		//
		// deployment
		//

		{
			event:  "deployment",
			before: "testdata/webhooks/deployment.json",
			after:  "testdata/webhooks/deployment.json.golden",
			obj:    new(scm.DeployHook),
		},
		{
			event:  "deployment",
			before: "testdata/webhooks/deployment_commit.json",
			after:  "testdata/webhooks/deployment_commit.json.golden",
			obj:    new(scm.DeployHook),
		},
		//
		// release
		//
		{
			event:  "release",
			before: "testdata/webhooks/release_published.json",
			after:  "testdata/webhooks/release_published.json.golden",
			obj:    new(scm.ReleaseHook),
		},
		{
			event:  "release",
			before: "testdata/webhooks/release_unpublished.json",
			after:  "testdata/webhooks/release_unpublished.json.golden",
			obj:    new(scm.ReleaseHook),
		},
		{
			event:  "release",
			before: "testdata/webhooks/release_created.json",
			after:  "testdata/webhooks/release_created.json.golden",
			obj:    new(scm.ReleaseHook),
		},
		{
			event:  "release",
			before: "testdata/webhooks/release_edited.json",
			after:  "testdata/webhooks/release_edited.json.golden",
			obj:    new(scm.ReleaseHook),
		},
		{
			event:  "release",
			before: "testdata/webhooks/release_deleted.json",
			after:  "testdata/webhooks/release_deleted.json.golden",
			obj:    new(scm.ReleaseHook),
		},
		{
			event:  "release",
			before: "testdata/webhooks/release_prereleased.json",
			after:  "testdata/webhooks/release_prereleased.json.golden",
			obj:    new(scm.ReleaseHook),
		},
		{
			event:  "release",
			before: "testdata/webhooks/release_released.json",
			after:  "testdata/webhooks/release_released.json.golden",
			obj:    new(scm.ReleaseHook),
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
		r.Header.Set("X-GitHub-Event", test.event)
		r.Header.Set("X-Hub-Signature-256", "sha256=3bfbbc3bfc44498db2254f577b2e4bed201ece6163518ba91cb2c21f0f59d512")
		r.Header.Set("X-GitHub-Delivery", "f2467dea-70d6-11e8-8955-3c83993e0aef")

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
	r.Header.Set("X-GitHub-Event", "push")
	r.Header.Set("X-GitHub-Delivery", "ee8d97b4-1479-43f1-9cac-fbbd1b80da55")
	r.Header.Set("X-Hub-Signature-256", "sha256=3bfbbc3bfc44498db2254f577b2e4bed201ece6163518ba91cb2c21f0f59d512")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != scm.ErrSignatureInvalid {
		t.Errorf("Expect invalid signature error, got %v", err)
	}
}

func TestWebhookValid(t *testing.T) {
	// the sha can be recalculated with the below command
	// openssl dgst -sha256 -hmac <secret> <file>

	f, _ := ioutil.ReadFile("testdata/webhooks/push.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
	r.Header.Set("X-GitHub-Event", "push")
	r.Header.Set("X-GitHub-Delivery", "ee8d97b4-1479-43f1-9cac-fbbd1b80da55")
	r.Header.Set("X-Hub-Signature-256", "sha256=e3bfe744d4e2e29ed990bde8acfb8255ca51ef65f99657767989fb6349f32957")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != nil {
		t.Errorf("Expect valid signature, got %v", err)
	}
}

func TestWebhookSignatureFallback(t *testing.T) {
	// the sha can be recalculated with the below command
	// openssl dgst -sha1 -hmac <secret> <file>

	f, _ := ioutil.ReadFile("testdata/webhooks/push.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
	r.Header.Set("X-GitHub-Event", "push")
	r.Header.Set("X-GitHub-Delivery", "ee8d97b4-1479-43f1-9cac-fbbd1b80da55")
	r.Header.Set("X-Hub-Signature", "sha1=cf93f9ba3c8d3a789e61f91e1e5c6a360d036e98")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != nil {
		t.Errorf("Expect valid signature, got %v", err)
	}
}

func secretFunc(scm.Webhook) (string, error) {
	return "topsecret", nil
}
