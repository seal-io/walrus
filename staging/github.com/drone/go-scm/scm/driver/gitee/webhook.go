// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/drone/go-scm/scm"
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
	switch req.Header.Get("X-Gitee-Event") {
	case "Push Hook":
		hook, err = s.parsePushHook(data)
	case "Merge Request Hook":
		hook, err = s.parseMergeRequestHook(data)
	case "Issue Hook":
		hook, err = s.parseIssueHook(data)
	case "Note Hook":
		hook, err = s.parseNoteHook(data)
	case "Tag Push Hook":
		hook, err = s.parseTagPushHook(data)
	default:
		return nil, scm.ErrUnknownEvent
	}
	if err != nil {
		return nil, err
	}

	key, err := fn(hook)
	if err != nil {
		return hook, err
	} else if key == "" {
		return hook, nil
	}

	agent := req.Header.Get("User-Agent")
	if agent != "git-oschina-hook" {
		return hook, &Error{
			Message: "hook's user-agent is not git-oschina-hook",
		}
	}
	timestamp := req.Header.Get("X-Gitee-Timestamp")
	signature := req.Header.Get("X-Gitee-Token")

	if !validateSignature(signature, key, timestamp) {
		return hook, scm.ErrSignatureInvalid
	}
	return hook, nil
}

func (s *webhookService) parsePushHook(data []byte) (scm.Webhook, error) {
	dst := new(pushOrTagPushHook)
	err := json.Unmarshal(data, dst)
	return convertPushHook(dst), err
}

func (s *webhookService) parseMergeRequestHook(data []byte) (scm.Webhook, error) {
	dst := new(mergeRequestHook)
	err := json.Unmarshal(data, dst)
	return convertPullRequestHook(dst), err
}

func (s *webhookService) parseIssueHook(data []byte) (scm.Webhook, error) {
	dst := new(issueHook)
	err := json.Unmarshal(data, dst)
	return convertIssueHook(dst), err
}

func (s *webhookService) parseTagPushHook(data []byte) (scm.Webhook, error) {
	dst := new(pushOrTagPushHook)
	err := json.Unmarshal(data, dst)
	return convertTagPushHook(dst), err
}

func (s *webhookService) parseNoteHook(data []byte) (scm.Webhook, error) {
	dst := new(noteHook)
	err := json.Unmarshal(data, dst)
	return convertNoteHook(dst), err
}

// validateSignature
// see https://gitee.com/help/articles/4290#article-header3
func validateSignature(signature, key, timestamp string) bool {
	stringToSign := timestamp + "\n" + key
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(stringToSign))
	computedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return computedSignature == signature
}

type (
	pushOrTagPushHook struct {
		Action             string         `json:"action"`
		HookName           string         `json:"hook_name"`
		Password           string         `json:"password"`
		HookID             int            `json:"hook_id"`
		HookURL            string         `json:"hook_url"`
		Timestamp          string         `json:"timestamp"`
		Sign               string         `json:"sign"`
		Ref                string         `json:"ref"`
		Before             string         `json:"before"`
		After              string         `json:"after"`
		Created            bool           `json:"created"`
		Deleted            bool           `json:"deleted"`
		Compare            string         `json:"compare"`
		Commits            []hookCommit   `json:"commits"`
		HeadCommit         hookCommit     `json:"head_commit"`
		TotalCommitsCount  int            `json:"total_commits_count"`
		CommitsMoreThanTen bool           `json:"commits_more_than_ten"`
		Repository         hookRepository `json:"repository"`
		Sender             user           `json:"sender"`
		Enterprise         enterprise     `json:"enterprise"`
	}
	mergeRequestHook struct {
		Action string `json:"action"`

		Number      int            `json:"number"`
		Title       string         `json:"title"`
		Project     hookRepository `json:"project"`
		PullRequest pr             `json:"pull_request"`

		Iid            int            `json:"iid"`
		ActionDesc     string         `json:"action_desc"`
		Author         user           `json:"author"`
		Body           string         `json:"body"`
		Enterprise     enterprise     `json:"enterprise"`
		Languages      []string       `json:"languages"`
		MergeCommitSha string         `json:"merge_commit_sha"`
		MergeStatus    string         `json:"merge_status"`
		Password       string         `json:"password"`
		Repository     hookRepository `json:"repository"`
		Sender         user           `json:"sender"`
		SourceBranch   string         `json:"source_branch"`

		SourceRepo struct {
			Project    hookRepository `json:"project"`
			Repository hookRepository `json:"repository"`
		} `json:"source_repo"`
		TargetRepo struct {
			Project    hookRepository `json:"project"`
			Repository hookRepository `json:"repository"`
		} `json:"target_repo"`
		State        string `json:"state"`
		TargetBranch string `json:"target_branch"`
		TargetUser   user   `json:"target_user"`
		Timestamp    string `json:"timestamp"`
		UpdatedBy    user   `json:"updated_by"`
		URL          string `json:"url"`
	}
	noteHook struct {
		Action       string         `json:"action"`
		HookName     string         `json:"hook_name"`
		Password     string         `json:"password"`
		HookID       int            `json:"hook_id"`
		HookURL      string         `json:"hook_url"`
		Timestamp    string         `json:"timestamp"`
		Sign         string         `json:"sign"`
		Comment      hookComment    `json:"comment"`
		NoteableType string         `json:"noteable_type"`
		Issue        issue          `json:"issue"`
		PullRequest  pr             `json:"pull_request"`
		Repository   hookRepository `json:"repository"`
		Sender       user           `json:"sender"`
		Enterprise   enterprise     `json:"enterprise"`
	}
	issueHook struct {
		Action     string         `json:"action"`
		HookName   string         `json:"hook_name"`
		Password   string         `json:"password"`
		HookID     int            `json:"hook_id"`
		HookURL    string         `json:"hook_url"`
		Timestamp  string         `json:"timestamp"`
		Sign       string         `json:"sign"`
		Issue      issue          `json:"issue"`
		Repository hookRepository `json:"repository"`
		Sender     user           `json:"sender"`
		Enterprise enterprise     `json:"enterprise"`
	}

	hookAuthorOrCommitter struct {
		Time     time.Time `json:"time"`
		Name     string    `json:"name"`
		Email    string    `json:"email"`
		Username string    `json:"username"`
		UserName string    `json:"user_name"`
		URL      string    `json:"url"`
	}
	hookCommit struct {
		ID        string                `json:"id"`
		TreeID    string                `json:"tree_id"`
		Distinct  bool                  `json:"distinct"`
		Message   string                `json:"message"`
		Timestamp time.Time             `json:"timestamp"`
		URL       string                `json:"url"`
		Author    hookAuthorOrCommitter `json:"author"`
		Committer hookAuthorOrCommitter `json:"committer"`
		Added     interface{}           `json:"added"`
		Removed   interface{}           `json:"removed"`
		Modified  []string              `json:"modified"`
	}
	hookComment struct {
		HtmlURL   string    `json:"html_url"`
		ID        int       `json:"id"`
		Body      string    `json:"body"`
		User      user      `json:"user"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	hookRepository struct {
		CloneURL          string      `json:"clone_url"`
		CreatedAt         time.Time   `json:"created_at"`
		DefaultBranch     string      `json:"default_branch"`
		Description       string      `json:"description"`
		Fork              bool        `json:"fork"`
		ForksCount        int         `json:"forks_count"`
		FullName          string      `json:"full_name"`
		GitHttpURL        string      `json:"git_http_url"`
		GitSshURL         string      `json:"git_ssh_url"`
		GitSvnURL         string      `json:"git_svn_url"`
		GitURL            string      `json:"git_url"`
		HasIssues         bool        `json:"has_issues"`
		HasPages          bool        `json:"has_pages"`
		HasWiki           bool        `json:"has_wiki"`
		Homepage          string      `json:"homepage"`
		HtmlURL           string      `json:"html_url"`
		ID                int         `json:"id"`
		Language          interface{} `json:"language"`
		License           interface{} `json:"license"`
		Name              string      `json:"name"`
		NameWithNamespace string      `json:"name_with_namespace"`
		Namespace         string      `json:"namespace"`
		OpenIssuesCount   int         `json:"open_issues_count"`
		Owner             user        `json:"owner"`
		Path              string      `json:"path"`
		PathWithNamespace string      `json:"path_with_namespace"`
		Private           bool        `json:"private"`
		PushedAt          time.Time   `json:"pushed_at"`
		SSHURL            string      `json:"ssh_url"`
		StargazersCount   int         `json:"stargazers_count"`
		SvnURL            string      `json:"svn_url"`
		UpdatedAt         time.Time   `json:"updated_at"`
		URL               string      `json:"url"`
		WatchersCount     int         `json:"watchers_count"`
	}
	enterprise struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
)

func convertPushHook(src *pushOrTagPushHook) *scm.PushHook {
	var commits []scm.Commit
	if &src.Commits != nil {
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
						Date:  c.Timestamp,
					},
					Committer: scm.Signature{
						Login: c.Committer.Username,
						Email: c.Committer.Email,
						Name:  c.Committer.Name,
						Date:  c.Timestamp,
					},
				})
		}
	}

	dst := &scm.PushHook{
		Ref:    src.Ref,
		Repo:   *convertHookRepository(&src.Repository),
		Before: src.Before,
		After:  src.After,
		Commit: scm.Commit{
			Sha:  src.After,
			Link: src.Compare,
		},
		Sender:  *convertUser(&src.Sender),
		Commits: commits,
	}
	if &src.HeadCommit != nil {
		dst.Commit.Message = src.HeadCommit.Message
		dst.Commit.Author = scm.Signature{
			Login: src.HeadCommit.Author.Username,
			Email: src.HeadCommit.Author.Email,
			Name:  src.HeadCommit.Author.Name,
			Date:  src.HeadCommit.Timestamp,
		}
		dst.Commit.Committer = scm.Signature{
			Login: src.HeadCommit.Committer.Username,
			Email: src.HeadCommit.Committer.Email,
			Name:  src.HeadCommit.Committer.Name,
			Date:  src.HeadCommit.Timestamp,
		}
	}
	return dst
}

func convertPullRequestHook(src *mergeRequestHook) *scm.PullRequestHook {
	dst := &scm.PullRequestHook{
		Repo:        *convertHookRepository(&src.Repository),
		PullRequest: *convertPullRequest(&src.PullRequest),
		Sender:      *convertUser(&src.Sender),
	}
	switch src.Action {
	case "update":
		if src.ActionDesc == "update_label" {
			if len(src.PullRequest.Labels) == 0 {
				dst.Action = scm.ActionUnlabel
			} else {
				dst.Action = scm.ActionLabel
			}
		} else if src.ActionDesc == "source_branch_changed" {
			// Gitee does not provide a synchronize action.
			// But when action_desc is 'source_branch_changed',
			// what happens is the same as GitHub's synchronize
			dst.Action = scm.ActionSync
		}
	case "open":
		dst.Action = scm.ActionOpen
	case "close":
		dst.Action = scm.ActionClose
	case "merge":
		dst.Action = scm.ActionMerge
	case "test", "tested", "assign", "approved":
		dst.Action = scm.ActionUnknown
	default:
		dst.Action = scm.ActionUnknown
	}
	return dst
}

func convertIssueHook(src *issueHook) *scm.IssueHook {
	dst := &scm.IssueHook{
		Repo:   *convertHookRepository(&src.Repository),
		Issue:  *convertIssue(&src.Issue),
		Sender: *convertUser(&src.Sender),
	}
	switch src.Action {
	case "open":
		dst.Action = scm.ActionOpen
	case "delete":
		dst.Action = scm.ActionClose
	case "state_change":
		switch src.Issue.State {
		case "open", "progressing":
			dst.Action = scm.ActionOpen
		case "close", "rejected":
			dst.Action = scm.ActionClose
		}
	case "assign":
		dst.Action = scm.ActionUpdate
	default:
		dst.Action = scm.ActionUnknown
	}
	return dst
}

func convertTagPushHook(src *pushOrTagPushHook) scm.Webhook {
	dst := &scm.TagHook{
		Ref: scm.Reference{
			Name: scm.TrimRef(src.Ref),
			Sha:  src.HeadCommit.ID,
		},
		Repo:   *convertHookRepository(&src.Repository),
		Sender: *convertUser(&src.Sender),
	}
	if src.Created {
		dst.Action = scm.ActionCreate
	} else if src.Deleted {
		dst.Action = scm.ActionDelete
		dst.Ref.Sha = ""
	} else {
		dst.Action = scm.ActionUnknown
	}
	return dst
}

func convertNoteHook(src *noteHook) scm.Webhook {
	convertHookComment := func(comment *hookComment) *scm.Comment {
		return &scm.Comment{
			ID:      comment.ID,
			Body:    comment.Body,
			Author:  *convertUser(&comment.User),
			Created: comment.CreatedAt,
			Updated: comment.UpdatedAt,
		}
	}
	convertCommentAction := func(src string) (action scm.Action) {
		switch src {
		case "comment":
			return scm.ActionCreate
		case "edited":
			return scm.ActionEdit
		case "deleted":
			return scm.ActionDelete
		default:
			return scm.ActionUnknown
		}
	}

	if src.NoteableType == "Issue" {
		return &scm.IssueCommentHook{
			Action:  convertCommentAction(src.Action),
			Repo:    *convertHookRepository(&src.Repository),
			Issue:   *convertIssue(&src.Issue),
			Comment: *convertHookComment(&src.Comment),
			Sender:  *convertUser(&src.Sender),
		}
	}

	if src.NoteableType == "PullRequest" {
		// not support review comment
		return &scm.PullRequestCommentHook{
			Action:      convertCommentAction(src.Action),
			Repo:        *convertHookRepository(&src.Repository),
			PullRequest: *convertPullRequest(&src.PullRequest),
			Comment:     *convertHookComment(&src.Comment),
			Sender:      *convertUser(&src.Sender),
		}
	}
	return nil
}

func convertHookRepository(from *hookRepository) *scm.Repository {
	return &scm.Repository{
		ID:        fmt.Sprint(from.ID),
		Namespace: from.Namespace,
		Name:      from.Name,
		Branch:    from.DefaultBranch,
		Private:   from.Private,
		Clone:     from.CloneURL,
		CloneSSH:  from.GitSshURL,
		Link:      from.HtmlURL,
		Created:   from.CreatedAt,
		Updated:   from.UpdatedAt,
	}
}
