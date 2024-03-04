// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"crypto/sha256"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/hmac"
)

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (scm.Webhook, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}

	var hook scm.Webhook
	switch req.Header.Get("X-Harness-Trigger") {
	// case "create":
	// 	hook, err = s.parseCreateHook(data)
	// case "delete":
	// 	hook, err = s.parseDeleteHook(data)
	// case "issues":
	// 	hook, err = s.parseIssueHook(data)
	case "branch_created", "branch_updated":
		hook, err = s.parsePushHook(data)
	case "pullreq_created", "pullreq_reopened", "pullreq_branch_updated", "pullreq_closed", "pullreq_merged":
		hook, err = s.parsePullRequestHook(data)
	case "pullreq_comment_created":
		hook, err = s.parsePullRequestCommentHook(data)
	default:
		return nil, scm.ErrUnknownEvent
	}
	if err != nil {
		return nil, err
	}

	// get the gitea signature key to verify the payload
	// signature. If no key is provided, no validation
	// is performed.
	key, err := fn(hook)
	if err != nil {
		return hook, err
	} else if key == "" {
		return hook, nil
	}

	secret := req.FormValue("secret")
	signature := req.Header.Get("X-Harness-Signature")

	// fail if no signature passed
	if signature == "" && secret == "" {
		return hook, scm.ErrSignatureInvalid
	}

	// test signature if header not set and secret is in payload
	if signature == "" && secret != "" && secret != key {
		return hook, scm.ErrSignatureInvalid
	}

	// test signature using header
	if signature != "" && !hmac.Validate(sha256.New, data, []byte(key), signature) {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePullRequestHook(data []byte) (scm.Webhook, error) {
	dst := new(pullRequestHook)
	err := json.Unmarshal(data, dst)
	return convertPullRequestHook(dst), err
}

func (s *webhookService) parsePushHook(data []byte) (scm.Webhook, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	return convertPushHook(dst), err
}

func (s *webhookService) parsePullRequestCommentHook(data []byte) (scm.Webhook, error) {
	dst := new(pullRequestCommentHook)
	err := json.Unmarshal(data, dst)
	return convertPullRequestCommentHook(dst), err
}

// native data structures
type (
	repo struct {
		ID            int    `json:"id"`
		Path          string `json:"path"`
		UID           string `json:"uid"`
		DefaultBranch string `json:"default_branch"`
		GitURL        string `json:"git_url"`
	}
	principal struct {
		ID          int    `json:"id"`
		UID         string `json:"uid"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Type        string `json:"type"`
		Created     int64  `json:"created"`
		Updated     int64  `json:"updated"`
	}
	pullReq struct {
		Number        int         `json:"number"`
		State         string      `json:"state"`
		IsDraft       bool        `json:"is_draft"`
		Title         string      `json:"title"`
		SourceRepoID  int         `json:"source_repo_id"`
		SourceBranch  string      `json:"source_branch"`
		TargetRepoID  int         `json:"target_repo_id"`
		TargetBranch  string      `json:"target_branch"`
		MergeStrategy interface{} `json:"merge_strategy"`
		Author        principal   `json:"author"`
		PrURL         string      `json:"pr_url"`
	}
	targetRef struct {
		Name string `json:"name"`
		Repo struct {
			ID            int    `json:"id"`
			Path          string `json:"path"`
			UID           string `json:"uid"`
			DefaultBranch string `json:"default_branch"`
			GitURL        string `json:"git_url"`
		} `json:"repo"`
	}
	ref struct {
		Name string `json:"name"`
		Repo struct {
			ID            int    `json:"id"`
			Path          string `json:"path"`
			UID           string `json:"uid"`
			DefaultBranch string `json:"default_branch"`
			GitURL        string `json:"git_url"`
		} `json:"repo"`
	}
	hookCommit struct {
		Sha     string `json:"sha"`
		Message string `json:"message"`
		Author  struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When string `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When string `json:"when"`
		} `json:"committer"`
		Added    []string `json:"added"`
		Modified []string `json:"modified"`
		Removed  []string `json:"removed"`
	}
	comment struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
	}
	// harness pull request webhook payload
	pullRequestHook struct {
		Trigger           string       `json:"trigger"`
		Repo              repo         `json:"repo"`
		Principal         principal    `json:"principal"`
		PullReq           pullReq      `json:"pull_req"`
		TargetRef         targetRef    `json:"target_ref"`
		Ref               ref          `json:"ref"`
		Sha               string       `json:"sha"`
		HeadCommit        hookCommit   `json:"head_commit"`
		Commits           []hookCommit `json:"commits"`
		TotalCommitsCount int64        `json:"total_commits_count"`
	}
	// harness push webhook payload
	pushHook struct {
		Trigger           string       `json:"trigger"`
		Repo              repo         `json:"repo"`
		Principal         principal    `json:"principal"`
		Ref               ref          `json:"ref"`
		HeadCommit        hookCommit   `json:"head_commit"`
		Sha               string       `json:"sha"`
		OldSha            string       `json:"old_sha"`
		Forced            bool         `json:"forced"`
		Commits           []hookCommit `json:"commits"`
		TotalCommitsCount int64        `json:"total_commits_count"`
	}
	// harness pull request comment webhook payload
	pullRequestCommentHook struct {
		Trigger    string     `json:"trigger"`
		Repo       repo       `json:"repo"`
		Principal  principal  `json:"principal"`
		PullReq    pullReq    `json:"pull_req"`
		TargetRef  targetRef  `json:"target_ref"`
		Ref        ref        `json:"ref"`
		Sha        string     `json:"sha"`
		HeadCommit hookCommit `json:"head_commit"`
		Comment    comment    `json:"comment"`
	}
)

// native data structure conversion
func convertPullRequestHook(src *pullRequestHook) *scm.PullRequestHook {
	return &scm.PullRequestHook{
		Action:      convertAction(src.Trigger),
		PullRequest: convertPullReq(src.PullReq, src.Ref, src.HeadCommit),
		Repo:        convertRepo(src.Repo),
		Sender:      convertUser(src.Principal),
	}
}

func convertPushHook(src *pushHook) *scm.PushHook {
	var commits []scm.Commit
	for _, c := range src.Commits {
		commits = append(commits, convertHookCommit(c))
	}
	return &scm.PushHook{
		Ref:     src.Ref.Name,
		Before:  src.OldSha,
		After:   src.Sha,
		Repo:    convertRepo(src.Repo),
		Commit:  convertHookCommit(src.HeadCommit),
		Sender:  convertUser(src.Principal),
		Commits: commits,
	}
}

func convertHookCommit(c hookCommit) scm.Commit {
	return scm.Commit{
		Sha:     c.Sha,
		Message: c.Message,
		Author: scm.Signature{
			Name:  c.Author.Identity.Name,
			Email: c.Author.Identity.Email,
		},
		Committer: scm.Signature{
			Name:  c.Committer.Identity.Name,
			Email: c.Committer.Identity.Email,
		},
	}
}

func convertPullRequestCommentHook(src *pullRequestCommentHook) *scm.PullRequestCommentHook {
	return &scm.PullRequestCommentHook{
		PullRequest: convertPullReq(src.PullReq, src.Ref, src.HeadCommit),
		Repo:        convertRepo(src.Repo),
		Comment: scm.Comment{
			Body: src.Comment.Text,
			ID:   src.Comment.ID,
		},
		Sender: convertUser(src.Principal),
	}
}

func convertAction(src string) (action scm.Action) {
	switch src {
	case "pullreq_created":
		return scm.ActionCreate
	case "pullreq_branch_updated":
		return scm.ActionUpdate
	case "pullreq_reopened":
		return scm.ActionReopen
	case "pullreq_closed":
		return scm.ActionClose
	case "pullreq_merged":
		return scm.ActionMerge
	default:
		return
	}
}

func convertPullReq(pr pullReq, ref ref, commit hookCommit) scm.PullRequest {
	return scm.PullRequest{
		Number: pr.Number,
		Title:  pr.Title,
		Closed: pr.State != "open",
		Source: pr.SourceBranch,
		Target: pr.TargetBranch,
		Merged: pr.State == "merged",
		Fork:   "fork",
		Link:   pr.PrURL,
		Sha:    commit.Sha,
		Ref:    ref.Name,
		Author: convertUser(pr.Author),
	}
}

func convertRepo(repo repo) scm.Repository {
	return scm.Repository{
		ID:     strconv.Itoa(repo.ID),
		Name:   repo.UID,
		Branch: repo.DefaultBranch,
		Link:   repo.GitURL,
		Clone:  repo.GitURL,
	}
}

func convertUser(principal principal) scm.User {
	return scm.User{
		Name:    principal.DisplayName,
		ID:      principal.UID,
		Login:   principal.UID,
		Email:   principal.Email,
		Created: time.Unix(0, principal.Created*int64(time.Millisecond)),
		Updated: time.Unix(0, principal.Updated*int64(time.Millisecond)),
	}
}
