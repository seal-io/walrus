// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/hmac"
	"github.com/drone/go-scm/scm/driver/internal/null"
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
	switch req.Header.Get("X-GitHub-Event") {
	case "push":
		hook, err = s.parsePushHook(data)
	case "create":
		hook, err = s.parseCreateHook(data)
	case "delete":
		hook, err = s.parseDeleteHook(data)
	case "pull_request":
		hook, err = s.parsePullRequestHook(data)
	case "deployment":
		hook, err = s.parseDeploymentHook(data)
	// case "pull_request_review_comment":
	// case "issues":
	case "issue_comment":
		hook, err = s.parseIssueCommentHook(data)
	case "release":
		hook, err = s.parseReleaseHook(data)
	default:
		return nil, scm.ErrUnknownEvent
	}
	if err != nil {
		return nil, err
	}

	// get the gogs signature key to verify the payload
	// signature. If no key is provided, no validation
	// is performed.
	key, err := fn(hook)
	if err != nil {
		return hook, err
	} else if key == "" {
		return hook, nil
	}

	sig := req.Header.Get("X-Hub-Signature-256")
	if sig == "" {
		sig = req.Header.Get("X-Hub-Signature")
	}
	if !hmac.ValidatePrefix(data, []byte(key), sig) {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePushHook(data []byte) (scm.Webhook, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	return convertPushHook(dst), err
}

func (s *webhookService) parseCreateHook(data []byte) (scm.Webhook, error) {
	src := new(createDeleteHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	if src.RefType == "branch" {
		dst := convertBranchHook(src)
		dst.Action = scm.ActionCreate
		return dst, nil
	}
	dst := convertTagHook(src)
	dst.Action = scm.ActionCreate
	return dst, nil
}

func (s *webhookService) parseDeleteHook(data []byte) (scm.Webhook, error) {
	src := new(createDeleteHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	if src.RefType == "branch" {
		dst := convertBranchHook(src)
		dst.Action = scm.ActionDelete
		return dst, nil
	}
	dst := convertTagHook(src)
	dst.Action = scm.ActionDelete
	return dst, nil
}

func (s *webhookService) parseDeploymentHook(data []byte) (scm.Webhook, error) {
	src := new(deploymentHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertDeploymentHook(src)
	return dst, nil
}

func (s *webhookService) parseIssueCommentHook(data []byte) (scm.Webhook, error) {
	src := new(issueCommentHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertIssueCommentHook(src)
	switch src.Action {
	case "created":
		dst.Action = scm.ActionCreate
	case "edited":
		dst.Action = scm.ActionEdit
	case "deleted":
		dst.Action = scm.ActionDelete
	default:
		dst.Action = scm.ActionUnknown
	}
	return dst, nil
}

func (s *webhookService) parsePullRequestHook(data []byte) (scm.Webhook, error) {
	src := new(pullRequestHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertPullRequestHook(src)
	switch src.Action {
	case "labeled":
		dst.Action = scm.ActionLabel
	case "unlabeled":
		dst.Action = scm.ActionUnlabel
	case "opened":
		dst.Action = scm.ActionOpen
	case "edited":
		dst.Action = scm.ActionUpdate
	case "closed":
		// TODO(bradrydzewski) github does not provide a merged action,
		// but this is provided by gitlab and bitbucket. Is it possible
		// to emulate the merge action?

		// if merged == true
		//    dst.Action = scm.ActionMerge
		dst.Action = scm.ActionClose
	case "reopened":
		dst.Action = scm.ActionReopen
	case "synchronize":
		dst.Action = scm.ActionSync
	case "ready_for_review":
		dst.Action = scm.ActionReviewReady
	case "assigned", "unassigned", "review_requested", "review_request_removed", "locked", "unlocked":
		dst.Action = scm.ActionUnknown
	default:
		dst.Action = scm.ActionUnknown
	}
	return dst, nil
}

func (s *webhookService) parseReleaseHook(data []byte) (scm.Webhook, error) {
	src := new(releaseHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertReleaseHook(src)
	switch src.Action {
	case "created":
		dst.Action = scm.ActionCreate
	case "edited":
		dst.Action = scm.ActionEdit
	case "deleted":
		dst.Action = scm.ActionDelete
	case "published":
		dst.Action = scm.ActionPublish
	case "unpublished":
		dst.Action = scm.ActionUnpublish
	case "prereleased":
		dst.Action = scm.ActionPrerelease
	case "released":
		dst.Action = scm.ActionRelease
	default:
		dst.Action = scm.ActionUnknown
	}
	return dst, nil
}

//
// native data structures
//

type (
	// github create webhook payload
	createDeleteHook struct {
		Ref        string     `json:"ref"`
		RefType    string     `json:"ref_type"`
		Repository repository `json:"repository"`
		Sender     user       `json:"sender"`
	}

	// github push webhook payload
	pushHook struct {
		Ref     string `json:"ref"`
		BaseRef string `json:"base_ref"`
		Before  string `json:"before"`
		After   string `json:"after"`
		Compare string `json:"compare"`
		Head    struct {
			ID        string    `json:"id"`
			TreeID    string    `json:"tree_id"`
			Distinct  bool      `json:"distinct"`
			Message   string    `json:"message"`
			Timestamp null.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"author"`
			Committer struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"committer"`
			Added    []interface{} `json:"added"`
			Removed  []interface{} `json:"removed"`
			Modified []string      `json:"modified"`
		} `json:"head_commit"`
		Commits []struct {
			ID        string    `json:"id"`
			TreeID    string    `json:"tree_id"`
			Distinct  bool      `json:"distinct"`
			Message   string    `json:"message"`
			Timestamp null.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"author"`
			Committer struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"committer"`
			Added    []interface{} `json:"added"`
			Removed  []interface{} `json:"removed"`
			Modified []string      `json:"modified"`
		} `json:"commits"`
		Repository struct {
			ID    int64 `json:"id"`
			Owner struct {
				Login     string `json:"login"`
				AvatarURL string `json:"avatar_url"`
			} `json:"owner"`
			Name          string `json:"name"`
			FullName      string `json:"full_name"`
			Private       bool   `json:"private"`
			Visibility    string `json:"visibility"`
			Fork          bool   `json:"fork"`
			HTMLURL       string `json:"html_url"`
			SSHURL        string `json:"ssh_url"`
			CloneURL      string `json:"clone_url"`
			DefaultBranch string `json:"default_branch"`
		} `json:"repository"`
		Pusher user `json:"pusher"`
		Sender user `json:"sender"`
	}

	pullRequestHook struct {
		Action      string     `json:"action"`
		Number      int        `json:"number"`
		PullRequest pr         `json:"pull_request"`
		Repository  repository `json:"repository"`
		Sender      user       `json:"sender"`
	}

	// github deployment webhook payload
	deploymentHook struct {
		Deployment struct {
			ID             int64       `json:"id"`
			Creator        user        `json:"creator"`
			Description    null.String `json:"description"`
			Environment    null.String `json:"environment"`
			EnvironmentURL null.String `json:"environment_url"`
			URL            null.String `json:"url"`
			Sha            null.String `json:"sha"`
			Ref            null.String `json:"ref"`
			Task           null.String `json:"task"`
			Payload        interface{} `json:"payload"`
		} `json:"deployment"`
		Repository repository `json:"repository"`
		Sender     user       `json:"sender"`
	}

	// github issue_comment webhook payload
	issueCommentHook struct {
		Action       string     `json:"action"`
		Issue        issue      `json:"issue"`
		Repository   repository `json:"repository"`
		Sender       user       `json:"sender"`
		Organization user       `json:"organization"`
		Comment      struct {
			ID      int       `json:"id"`
			Body    string    `json:"body"`
			User    user      `json:"user"`
			Created time.Time `json:"created_at"`
			Updated time.Time `json:"updated_at"`
		} `json:"comment"`
	}

	// github release webhook payload
	releaseHook struct {
		Action  string `json:"action"`
		Release struct {
			ID          int       `json:"id"`
			Title       string    `json:"name"`
			Description string    `json:"body"`
			Link        string    `json:"html_url,omitempty"`
			Tag         string    `json:"tag_name,omitempty"`
			Commitish   string    `json:"target_commitish,omitempty"`
			Draft       bool      `json:"draft"`
			Prerelease  bool      `json:"prerelease"`
			Created     time.Time `json:"created_at"`
			Published   time.Time `json:"published_at"`
		} `json:"release"`
		Repository repository `json:"repository"`
		Sender     user       `json:"sender"`
	}
)

//
// native data structure conversion
//

func convertPushHook(src *pushHook) *scm.PushHook {
	var commits []scm.Commit
	for _, c := range src.Commits {
		commits = append(commits,
			scm.Commit{
				Sha:     c.ID,
				Message: c.Message,
				Link:    c.URL,
				Author: scm.Signature{
					Login: c.Author.Username,
					Email: c.Author.Email,
					Name:  c.Author.Name,
					Date:  c.Timestamp.ValueOrZero(),
				},
				Committer: scm.Signature{
					Login: c.Committer.Username,
					Email: c.Committer.Email,
					Name:  c.Committer.Name,
					Date:  c.Timestamp.ValueOrZero(),
				},
			})
	}
	dst := &scm.PushHook{
		Ref:     src.Ref,
		BaseRef: src.BaseRef,
		Before:  src.Before,
		After:   src.After,
		Commit: scm.Commit{
			Sha:     src.After,
			Message: src.Head.Message,
			Link:    src.Compare,
			Author: scm.Signature{
				Login: src.Head.Author.Username,
				Email: src.Head.Author.Email,
				Name:  src.Head.Author.Name,
				Date:  src.Head.Timestamp.ValueOrZero(),
			},
			Committer: scm.Signature{
				Login: src.Head.Committer.Username,
				Email: src.Head.Committer.Email,
				Name:  src.Head.Committer.Name,
				Date:  src.Head.Timestamp.ValueOrZero(),
			},
		},
		Repo: scm.Repository{
			ID:         fmt.Sprint(src.Repository.ID),
			Namespace:  src.Repository.Owner.Login,
			Name:       src.Repository.Name,
			Branch:     src.Repository.DefaultBranch,
			Private:    src.Repository.Private,
			Visibility: convertVisibility(src.Repository.Visibility),
			Clone:      src.Repository.CloneURL,
			CloneSSH:   src.Repository.SSHURL,
			Link:       src.Repository.HTMLURL,
		},
		Sender:  *convertUser(&src.Sender),
		Commits: commits,
	}
	// fix https://github.com/drone/go-scm/issues/8
	if scm.IsTag(dst.Ref) && src.Head.ID != "" {
		dst.Commit.Sha = src.Head.ID
		dst.After = src.Head.ID
	}
	return dst
}

func convertBranchHook(src *createDeleteHook) *scm.BranchHook {
	return &scm.BranchHook{
		Ref: scm.Reference{
			Name: src.Ref,
		},
		Repo: scm.Repository{
			ID:         fmt.Sprint(src.Repository.ID),
			Namespace:  src.Repository.Owner.Login,
			Name:       src.Repository.Name,
			Branch:     src.Repository.DefaultBranch,
			Private:    src.Repository.Private,
			Visibility: convertVisibility(src.Repository.Visibility),
			Clone:      src.Repository.CloneURL,
			CloneSSH:   src.Repository.SSHURL,
			Link:       src.Repository.HTMLURL,
		},
		Sender: *convertUser(&src.Sender),
	}
}

func convertTagHook(src *createDeleteHook) *scm.TagHook {
	return &scm.TagHook{
		Ref: scm.Reference{
			Name: src.Ref,
		},
		Repo: scm.Repository{
			ID:         fmt.Sprint(src.Repository.ID),
			Namespace:  src.Repository.Owner.Login,
			Name:       src.Repository.Name,
			Branch:     src.Repository.DefaultBranch,
			Private:    src.Repository.Private,
			Visibility: convertVisibility(src.Repository.Visibility),
			Clone:      src.Repository.CloneURL,
			CloneSSH:   src.Repository.SSHURL,
			Link:       src.Repository.HTMLURL,
		},
		Sender: *convertUser(&src.Sender),
	}
}

func convertPullRequestHook(src *pullRequestHook) *scm.PullRequestHook {
	return &scm.PullRequestHook{
		// Action        Action
		Repo: scm.Repository{
			ID:         fmt.Sprint(src.Repository.ID),
			Namespace:  src.Repository.Owner.Login,
			Name:       src.Repository.Name,
			Branch:     src.Repository.DefaultBranch,
			Private:    src.Repository.Private,
			Visibility: convertVisibility(src.Repository.Visibility),
			Clone:      src.Repository.CloneURL,
			CloneSSH:   src.Repository.SSHURL,
			Link:       src.Repository.HTMLURL,
		},
		PullRequest: *convertPullRequest(&src.PullRequest),
		Sender:      *convertUser(&src.Sender),
	}
}

func convertDeploymentHook(src *deploymentHook) *scm.DeployHook {
	dst := &scm.DeployHook{
		Number: src.Deployment.ID,
		Data:   src.Deployment.Payload,
		Desc:   src.Deployment.Description.String,
		Ref: scm.Reference{
			Name: src.Deployment.Ref.String,
			Path: src.Deployment.Ref.String,
			Sha:  src.Deployment.Sha.String,
		},
		Repo: scm.Repository{
			ID:         fmt.Sprint(src.Repository.ID),
			Namespace:  src.Repository.Owner.Login,
			Name:       src.Repository.Name,
			Branch:     src.Repository.DefaultBranch,
			Private:    src.Repository.Private,
			Visibility: convertVisibility(src.Repository.Visibility),
			Clone:      src.Repository.CloneURL,
			CloneSSH:   src.Repository.SSHURL,
			Link:       src.Repository.HTMLURL,
		},
		Sender:    *convertUser(&src.Sender),
		Task:      src.Deployment.Task.String,
		Target:    src.Deployment.Environment.String,
		TargetURL: src.Deployment.EnvironmentURL.String,
	}
	// handle deployment events for commits which
	// use a different payload structure and lack
	// branch or reference details.
	if len(dst.Ref.Name) == 40 && dst.Ref.Name == dst.Ref.Sha {
		dst.Ref.Name = ""
		dst.Ref.Path = ""
		return dst
	}
	if tagRE.MatchString(dst.Ref.Name) {
		dst.Ref.Path = scm.ExpandRef(dst.Ref.Path, "refs/tags/")
	} else {
		dst.Ref.Path = scm.ExpandRef(dst.Ref.Path, "refs/heads/")
	}
	return dst
}

func convertIssueCommentHook(src *issueCommentHook) *scm.IssueCommentHook {
	dst := &scm.IssueCommentHook{
		Repo: scm.Repository{
			ID:         fmt.Sprint(src.Repository.ID),
			Namespace:  src.Repository.Owner.Login,
			Name:       src.Repository.Name,
			Branch:     src.Repository.DefaultBranch,
			Private:    src.Repository.Private,
			Visibility: convertVisibility(src.Repository.Visibility),
			Clone:      src.Repository.CloneURL,
			CloneSSH:   src.Repository.SSHURL,
			Link:       src.Repository.HTMLURL,
		},
		Issue: *convertIssue(&src.Issue),
		Comment: scm.Comment{
			ID:      src.Comment.ID,
			Body:    src.Comment.Body,
			Author:  *convertUser(&src.Comment.User),
			Created: src.Comment.Created,
			Updated: src.Comment.Updated,
		},
		Sender: *convertUser(&src.Sender),
	}
	return dst
}

func convertReleaseHook(src *releaseHook) *scm.ReleaseHook {
	dst := &scm.ReleaseHook{
		Release: scm.Release{
			ID:          src.Release.ID,
			Title:       src.Release.Title,
			Description: src.Release.Description,
			Link:        src.Release.Link,
			Tag:         src.Release.Tag,
			Commitish:   src.Release.Commitish,
			Draft:       src.Release.Draft,
			Prerelease:  src.Release.Prerelease,
			Created:     src.Release.Created,
			Published:   src.Release.Published,
		},
		Repo: scm.Repository{
			ID:         fmt.Sprint(src.Repository.ID),
			Namespace:  src.Repository.Owner.Login,
			Name:       src.Repository.Name,
			Branch:     src.Repository.DefaultBranch,
			Private:    src.Repository.Private,
			Visibility: convertVisibility(src.Repository.Visibility),
			Clone:      src.Repository.CloneURL,
			CloneSSH:   src.Repository.SSHURL,
			Link:       src.Repository.HTMLURL,
		},
		Sender: *convertUser(&src.Sender),
	}
	return dst
}

// regexp help determine if the named git object is a tag.
// this is not meant to be 100% accurate.
var tagRE = regexp.MustCompile("^v?(\\d+).(.+)")
