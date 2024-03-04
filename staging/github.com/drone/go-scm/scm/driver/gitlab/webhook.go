// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
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
	switch req.Header.Get("X-Gitlab-Event") {
	case "Push Hook", "Tag Push Hook":
		hook, err = parsePushHook(data)
	case "Issue Hook":
		return nil, scm.ErrUnknownEvent
	case "Merge Request Hook":
		hook, err = parsePullRequestHook(data)
	case "Note Hook":
		hook, err = parseIssueCommentHook(data)
	default:
		return nil, scm.ErrUnknownEvent
	}
	if err != nil {
		return nil, err
	}

	// get the gitlab shared token to verify the payload
	// authenticity. If no key is provided, no validation
	// is performed.
	token, err := fn(hook)
	if err != nil {
		return hook, err
	} else if token == "" {
		return hook, nil
	}

	if token != req.Header.Get("X-Gitlab-Token") {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func parseIssueCommentHook(data []byte) (scm.Webhook, error) {
	src := new(commentHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst, err := convertCommentHook(src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func parsePushHook(data []byte) (scm.Webhook, error) {
	src := new(pushHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	switch {
	case src.ObjectKind == "push" && src.Before == "0000000000000000000000000000000000000000":
		// TODO we previously considered returning a
		// branch creation hook, however, the push hook
		// returns more metadata (commit details).
		return convertPushHook(src), nil
	case src.ObjectKind == "push" && src.After == "0000000000000000000000000000000000000000":
		return converBranchHook(src), nil
	case src.ObjectKind == "tag_push" && src.Before == "0000000000000000000000000000000000000000":
		// TODO we previously considered returning a
		// tag creation hook, however, the push hook
		// returns more metadata (commit details).
		return convertPushHook(src), nil
	case src.ObjectKind == "tag_push" && src.After == "0000000000000000000000000000000000000000":
		return convertTagHook(src), nil
	default:
		return convertPushHook(src), nil
	}
}

func parsePullRequestHook(data []byte) (scm.Webhook, error) {
	src := new(pullRequestHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	switch src.ObjectAttributes.Action {
	case "open", "close", "reopen", "merge", "update":
		// no-op
	default:
		return nil, scm.ErrUnknownEvent
	}
	switch {
	default:
		return convertPullRequestHook(src), nil
	}
}

func convertPushHook(src *pushHook) *scm.PushHook {
	var commits []scm.Commit
	for _, c := range src.Commits {
		commits = append(commits,
			scm.Commit{
				Sha:     c.ID,
				Message: c.Message,
				Link:    c.URL,
				Author: scm.Signature{
					Name:  c.Author.Name,
					Email: c.Author.Email,
					Date:  c.Timestamp.ValueOrZero(),
				},
				Committer: scm.Signature{
					Name:  c.Author.Name,
					Email: c.Author.Email,
					Date:  c.Timestamp.ValueOrZero(),
				},
			})
	}
	namespace, name := scm.Split(src.Project.PathWithNamespace)
	dst := &scm.PushHook{
		Ref:    scm.ExpandRef(src.Ref, "refs/heads/"),
		Before: src.Before,
		After:  src.After,
		Repo: scm.Repository{
			ID:        strconv.Itoa(src.Project.ID),
			Namespace: namespace,
			Name:      name,
			Clone:     src.Project.GitHTTPURL,
			CloneSSH:  src.Project.GitSSHURL,
			Link:      src.Project.WebURL,
			Branch:    src.Project.DefaultBranch,
			Private:   false, // TODO how do we correctly set Private vs Public?
		},
		Commit: scm.Commit{
			Sha:     src.CheckoutSha,
			Message: "", // NOTE this is set below
			Author: scm.Signature{
				Login:  src.UserUsername,
				Name:   src.UserName,
				Email:  src.UserEmail,
				Avatar: src.UserAvatar,
			},
			Committer: scm.Signature{
				Login:  src.UserUsername,
				Name:   src.UserName,
				Email:  src.UserEmail,
				Avatar: src.UserAvatar,
			},
			Link: "", // NOTE this is set below
		},
		Sender: scm.User{
			Login:  src.UserUsername,
			Name:   src.UserName,
			Email:  src.UserEmail,
			Avatar: src.UserAvatar,
		},
		Commits: commits,
	}
	if len(src.Commits) > 0 {
		// get the last commit (most recent)
		dst.Commit.Message = src.Commits[len(src.Commits)-1].Message
		dst.Commit.Link = src.Commits[len(src.Commits)-1].URL
	}
	return dst
}

func converBranchHook(src *pushHook) *scm.BranchHook {
	action := scm.ActionCreate
	commit := src.After
	if src.After == "0000000000000000000000000000000000000000" {
		action = scm.ActionDelete
		commit = src.Before
	}
	namespace, name := scm.Split(src.Project.PathWithNamespace)
	return &scm.BranchHook{
		Action: action,
		Ref: scm.Reference{
			Name: scm.TrimRef(src.Ref),
			Sha:  commit,
		},
		Repo: scm.Repository{
			ID:        strconv.Itoa(src.Project.ID),
			Namespace: namespace,
			Name:      name,
			Clone:     src.Project.GitHTTPURL,
			CloneSSH:  src.Project.GitSSHURL,
			Link:      src.Project.WebURL,
			Branch:    src.Project.DefaultBranch,
			Private:   false, // TODO how do we correctly set Private vs Public?
		},
		Sender: scm.User{
			Login:  src.UserUsername,
			Name:   src.UserName,
			Email:  src.UserEmail,
			Avatar: src.UserAvatar,
		},
	}
}

func convertCommentHook(src *commentHook) (*scm.IssueCommentHook, error) {
	var issue scm.Issue
	var comment scm.Comment

	switch src.ObjectAttributes.NoteableType {
	case "Commit", "Issue", "Snippet":
		return nil, scm.ErrUnknownEvent
	case "MergeRequest":
		pr := scm.PullRequest{
			Number:  src.MergeRequest.Iid,
			Title:   src.MergeRequest.Title,
			Body:    src.MergeRequest.Description,
			Sha:     src.MergeRequest.LastCommit.ID,
			Ref:     fmt.Sprintf("refs/merge-requests/%d/head", src.MergeRequest.Iid),
			Source:  src.MergeRequest.SourceBranch,
			Target:  src.MergeRequest.TargetBranch,
			Link:    src.MergeRequest.URL,
			Closed:  src.MergeRequest.State != "opened",
			Merged:  src.MergeRequest.State == "merged",
			Author:  *convertUser(&src.User),
			Created: parseTimeString(src.MergeRequest.CreatedAt),
			Updated: parseTimeString(src.MergeRequest.UpdatedAt),
		}
		for _, l := range src.MergeRequest.Labels {
			label := scm.Label{
				Name:  l.Title,
				Color: l.Color,
			}
			pr.Labels = append(pr.Labels, label)
		}
		issue = scm.Issue{
			Number:      src.MergeRequest.Iid,
			Title:       src.MergeRequest.Title,
			Body:        src.MergeRequest.Title,
			Link:        src.Project.WebURL,
			Author:      *convertUser(&src.User),
			PullRequest: pr,
			Created:     parseTimeString(src.ObjectAttributes.CreatedAt),
			Updated:     parseTimeString(src.ObjectAttributes.UpdatedAt),
		}
		comment = scm.Comment{
			ID:      src.ObjectAttributes.ID,
			Body:    src.ObjectAttributes.Note,
			Author:  *convertUser(&src.User),
			Created: parseTimeString(src.ObjectAttributes.CreatedAt),
			Updated: parseTimeString(src.ObjectAttributes.UpdatedAt),
		}
	default:
		return nil, scm.ErrUnknownEvent
	}

	namespace, _ := scm.Split(src.Project.PathWithNamespace)
	dst := scm.IssueCommentHook{
		Action: scm.ActionCreate,
		Repo: scm.Repository{
			ID:        strconv.Itoa(src.Project.ID),
			Namespace: namespace,
			Name:      src.Repository.Name,
			Clone:     src.Project.GitHTTPURL,
			CloneSSH:  src.Project.GitSSHURL,
			Link:      src.Project.WebURL,
			Branch:    src.Project.DefaultBranch,
			Private:   false, // TODO how do we correctly set Private vs Public?
		},
		Issue:   issue,
		Comment: comment,
		Sender:  *convertUser(&src.User),
	}
	return &dst, nil
}

func convertTagHook(src *pushHook) *scm.TagHook {
	action := scm.ActionCreate
	commit := src.After
	if src.After == "0000000000000000000000000000000000000000" {
		action = scm.ActionDelete
		commit = src.Before
	}
	namespace, name := scm.Split(src.Project.PathWithNamespace)
	return &scm.TagHook{
		Action: action,
		Ref: scm.Reference{
			Name: scm.TrimRef(src.Ref),
			Sha:  commit,
		},
		Repo: scm.Repository{
			ID:        strconv.Itoa(src.Project.ID),
			Namespace: namespace,
			Name:      name,
			Clone:     src.Project.GitHTTPURL,
			CloneSSH:  src.Project.GitSSHURL,
			Link:      src.Project.WebURL,
			Branch:    src.Project.DefaultBranch,
			Private:   false, // TODO how do we correctly set Private vs Public?
		},
		Sender: scm.User{
			Login:  src.UserUsername,
			Name:   src.UserName,
			Email:  src.UserEmail,
			Avatar: src.UserAvatar,
		},
	}
}

func convertPullRequestHook(src *pullRequestHook) *scm.PullRequestHook {
	action := scm.ActionSync
	switch src.ObjectAttributes.Action {
	case "open":
		action = scm.ActionOpen
	case "close":
		action = scm.ActionClose
	case "reopen":
		action = scm.ActionReopen
	case "merge":
		action = scm.ActionMerge
	case "update":
		action = scm.ActionSync
		if src.Changes.Draft.Previous.Bool == false && src.Changes.Draft.Current.Bool == true {
			action = scm.ActionReviewReady
		}
	}
	fork := scm.Join(
		src.ObjectAttributes.Source.Namespace,
		src.ObjectAttributes.Source.Name,
	)
	namespace, name := scm.Split(src.Project.PathWithNamespace)
	return &scm.PullRequestHook{
		Action: action,
		PullRequest: scm.PullRequest{
			Number: src.ObjectAttributes.Iid,
			Title:  src.ObjectAttributes.Title,
			Body:   src.ObjectAttributes.Description,
			Sha:    src.ObjectAttributes.LastCommit.ID,
			Ref:    fmt.Sprintf("refs/merge-requests/%d/head", src.ObjectAttributes.Iid),
			Source: src.ObjectAttributes.SourceBranch,
			Target: src.ObjectAttributes.TargetBranch,
			Fork:   fork,
			Link:   src.ObjectAttributes.URL,
			Closed: src.ObjectAttributes.State != "opened",
			Merged: src.ObjectAttributes.State == "merged",
			// Created   : src.ObjectAttributes.CreatedAt,
			// Updated  : src.ObjectAttributes.UpdatedAt, // 2017-12-10 17:01:11 UTC
			Author: scm.User{
				Login:  src.User.Username,
				Name:   src.User.Name,
				Email:  "", // TODO how do we get the pull request author email?
				Avatar: src.User.AvatarURL,
			},
		},
		Repo: scm.Repository{
			ID:        strconv.Itoa(src.Project.ID),
			Namespace: namespace,
			Name:      name,
			Clone:     src.Project.GitHTTPURL,
			CloneSSH:  src.Project.GitSSHURL,
			Link:      src.Project.WebURL,
			Branch:    src.Project.DefaultBranch,
			Private:   false, // TODO how do we correctly set Private vs Public?
		},
		Sender: scm.User{
			Login:  src.User.Username,
			Name:   src.User.Name,
			Email:  "", // TODO how do we get the pull request author email?
			Avatar: src.User.AvatarURL,
		},
	}
}

func parseTimeString(timeString string) time.Time {
	layout := "2006-01-02 15:04:05 UTC"
	// Returns zero value of time in case of an error 0001-01-01 00:00:00 +0000 UTC
	t, _ := time.Parse(layout, timeString)
	return t
}

type (
	pushHook struct {
		ObjectKind   string      `json:"object_kind"`
		EventName    string      `json:"event_name"`
		Before       string      `json:"before"`
		After        string      `json:"after"`
		Ref          string      `json:"ref"`
		CheckoutSha  string      `json:"checkout_sha"`
		Message      interface{} `json:"message"`
		UserID       int         `json:"user_id"`
		UserName     string      `json:"user_name"`
		UserUsername string      `json:"user_username"`
		UserEmail    string      `json:"user_email"`
		UserAvatar   string      `json:"user_avatar"`
		ProjectID    int         `json:"project_id"`
		Project      struct {
			ID                int         `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      interface{} `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"project"`
		Commits []struct {
			ID        string    `json:"id"`
			Message   string    `json:"message"`
			Timestamp null.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
			Added    []string      `json:"added"`
			Modified []interface{} `json:"modified"`
			Removed  []interface{} `json:"removed"`
		} `json:"commits"`
		TotalCommitsCount int `json:"total_commits_count"`
		Repository        struct {
			Name            string `json:"name"`
			URL             string `json:"url"`
			Description     string `json:"description"`
			Homepage        string `json:"homepage"`
			GitHTTPURL      string `json:"git_http_url"`
			GitSSHURL       string `json:"git_ssh_url"`
			VisibilityLevel int    `json:"visibility_level"`
		} `json:"repository"`
	}

	commentHook struct {
		ObjectKind string `json:"object_kind"`
		EventType  string `json:"event_type"`
		User       user   `json:"user"`
		ProjectID  int    `json:"project_id"`
		Project    struct {
			ID                int         `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      interface{} `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"project"`
		ObjectAttributes struct {
			ID               int         `json:"id"`
			Note             string      `json:"note"`
			NoteableType     string      `json:"noteable_type"`
			AuthorID         int         `json:"author_id"`
			CreatedAt        string      `json:"created_at"`
			UpdatedAt        string      `json:"updated_at"`
			ProjectID        int         `json:"project_id"`
			Attachment       interface{} `json:"attachment"`
			LineCode         string      `json:"line_code"`
			CommitID         string      `json:"commit_id"`
			NoteableID       int         `json:"noteable_id"`
			StDiff           interface{} `json:"st_diff"`
			System           bool        `json:"system"`
			ResolvedAt       interface{} `json:"resolved_at"`
			ResolvedByID     interface{} `json:"resolved_by_id"`
			ResolvedByPush   interface{} `json:"resolved_by_push"`
			DiscussionID     string      `json:"discussion_id"`
			URL              string      `json:"url"`
			Position         interface{} `json:"position"`
			OriginalPosition interface{} `json:"original_position"`
			ChangePosition   interface{} `json:"change_position"`
			Type             interface{} `json:"type"`
			Description      string      `json:"description"`
		} `json:"object_attributes"`
		Repository struct {
			Name        string `json:"name"`
			URL         string `json:"url"`
			Description string `json:"description"`
			Homepage    string `json:"homepage"`
		} `json:"repository"`
		MergeRequest struct {
			AssigneeID                 interface{} `json:"assignee_id"`
			AuthorID                   int         `json:"author_id"`
			CreatedAt                  string      `json:"created_at"`
			DeletedAt                  interface{} `json:"deleted_at"`
			Description                string      `json:"description"`
			HeadPipelineID             interface{} `json:"head_pipeline_id"`
			ID                         int         `json:"id"`
			Iid                        int         `json:"iid"`
			LastEditedAt               interface{} `json:"last_edited_at"`
			LastEditedByID             interface{} `json:"last_edited_by_id"`
			MergeCommitSha             interface{} `json:"merge_commit_sha"`
			MergeError                 interface{} `json:"merge_error"`
			MergeParams                interface{} `json:"-"`
			MergeStatus                string      `json:"merge_status"`
			MergeUserID                interface{} `json:"merge_user_id"`
			MergeWhenPipelineSucceeds  bool        `json:"merge_when_pipeline_succeeds"`
			MilestoneID                interface{} `json:"milestone_id"`
			SourceBranch               string      `json:"source_branch"`
			SourceProjectID            int         `json:"source_project_id"`
			StateID                    int         `json:"state_id"`
			State                      string      `json:"state"`
			TargetBranch               string      `json:"target_branch"`
			TargetProjectID            int         `json:"target_project_id"`
			TimeEstimate               int         `json:"time_estimate"`
			Title                      string      `json:"title"`
			UpdatedAt                  string      `json:"updated_at"`
			UpdatedByID                interface{} `json:"updated_by_id"`
			URL                        string      `json:"url"`
			WorkInProgress             bool        `json:"work_in_progress"`
			TimeChange                 int         `json:"time_change"`
			HumanTimeChange            int         `json:"human_time_change"`
			TotalTimeSpent             int         `json:"total_time_spent"`
			HumanTotalTimeSpent        interface{} `json:"human_total_time_spent"`
			HumanTimeEstimate          interface{} `json:"human_time_estimate"`
			Action                     string      `json:"action"`
			AssigneeIDs                interface{} `json:"assignee_ids"`
			BlockingDiscussionResolved bool        `json:"blocking_discussions_resolved"`
			Source                     struct {
				ID                int         `json:"id"`
				Name              string      `json:"name"`
				Description       string      `json:"description"`
				WebURL            string      `json:"web_url"`
				AvatarURL         interface{} `json:"avatar_url"`
				GitSSHURL         string      `json:"git_ssh_url"`
				GitHTTPURL        string      `json:"git_http_url"`
				Namespace         string      `json:"namespace"`
				VisibilityLevel   int         `json:"visibility_level"`
				PathWithNamespace string      `json:"path_with_namespace"`
				DefaultBranch     string      `json:"default_branch"`
				CiConfigPath      interface{} `json:"ci_config_path"`
				Homepage          string      `json:"homepage"`
				URL               string      `json:"url"`
				SSHURL            string      `json:"ssh_url"`
				HTTPURL           string      `json:"http_url"`
			} `json:"source"`
			Target struct {
				ID                int         `json:"id"`
				Name              string      `json:"name"`
				Description       string      `json:"description"`
				WebURL            string      `json:"web_url"`
				AvatarURL         interface{} `json:"avatar_url"`
				GitSSHURL         string      `json:"git_ssh_url"`
				GitHTTPURL        string      `json:"git_http_url"`
				Namespace         string      `json:"namespace"`
				VisibilityLevel   int         `json:"visibility_level"`
				PathWithNamespace string      `json:"path_with_namespace"`
				DefaultBranch     string      `json:"default_branch"`
				CiConfigPath      interface{} `json:"ci_config_path"`
				Homepage          string      `json:"homepage"`
				URL               string      `json:"url"`
				SSHURL            string      `json:"ssh_url"`
				HTTPURL           string      `json:"http_url"`
			} `json:"target"`
			LastCommit struct {
				ID        string `json:"id"`
				Message   string `json:"message"`
				Timestamp string `json:"timestamp"`
				URL       string `json:"url"`
				Author    struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"author"`
			} `json:"last_commit"`
			Labels []struct {
				ID          int         `json:"id"`
				Title       string      `json:"title"`
				Color       string      `json:"color"`
				ProjectID   int         `json:"project_id"`
				CreatedAt   string      `json:"created_at"`
				UpdatedAt   string      `json:"updated_at"`
				Template    bool        `json:"template"`
				Description string      `json:"description"`
				Type        string      `json:"type"`
				GroupID     interface{} `json:"group_id"`
			} `json:"labels"`
		} `json:"merge_request"`
	}

	tagHook struct {
		ObjectKind   string      `json:"object_kind"`
		EventName    string      `json:"event_name"`
		Before       string      `json:"before"`
		After        string      `json:"after"`
		Ref          string      `json:"ref"`
		CheckoutSha  string      `json:"checkout_sha"`
		Message      interface{} `json:"message"`
		UserID       int         `json:"user_id"`
		UserName     string      `json:"user_name"`
		UserUsername string      `json:"user_username"`
		UserEmail    string      `json:"user_email"`
		UserAvatar   string      `json:"user_avatar"`
		ProjectID    int         `json:"project_id"`
		Project      struct {
			ID                int         `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      interface{} `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"project"`
		Commits []struct {
			ID        string `json:"id"`
			Message   string `json:"message"`
			Timestamp string `json:"timestamp"`
			URL       string `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
			Added    []string      `json:"added"`
			Modified []interface{} `json:"modified"`
			Removed  []interface{} `json:"removed"`
		} `json:"commits"`
		TotalCommitsCount int `json:"total_commits_count"`
		Repository        struct {
			Name            string `json:"name"`
			URL             string `json:"url"`
			Description     string `json:"description"`
			Homepage        string `json:"homepage"`
			GitHTTPURL      string `json:"git_http_url"`
			GitSSHURL       string `json:"git_ssh_url"`
			VisibilityLevel int    `json:"visibility_level"`
		} `json:"repository"`
	}

	issueHook struct {
		ObjectKind string `json:"object_kind"`
		User       struct {
			Name      string `json:"name"`
			Username  string `json:"username"`
			AvatarURL string `json:"avatar_url"`
		} `json:"user"`
		Project struct {
			ID                int         `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      interface{} `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"project"`
		ObjectAttributes struct {
			AssigneeID          interface{}   `json:"assignee_id"`
			AuthorID            int           `json:"author_id"`
			BranchName          interface{}   `json:"branch_name"`
			ClosedAt            interface{}   `json:"closed_at"`
			Confidential        bool          `json:"confidential"`
			CreatedAt           string        `json:"created_at"`
			DeletedAt           interface{}   `json:"deleted_at"`
			Description         string        `json:"description"`
			DueDate             interface{}   `json:"due_date"`
			ID                  int           `json:"id"`
			Iid                 int           `json:"iid"`
			LastEditedAt        string        `json:"last_edited_at"`
			LastEditedByID      int           `json:"last_edited_by_id"`
			MilestoneID         interface{}   `json:"milestone_id"`
			MovedToID           interface{}   `json:"moved_to_id"`
			ProjectID           int           `json:"project_id"`
			RelativePosition    int           `json:"relative_position"`
			State               string        `json:"state"`
			TimeEstimate        int           `json:"time_estimate"`
			Title               string        `json:"title"`
			UpdatedAt           string        `json:"updated_at"`
			UpdatedByID         int           `json:"updated_by_id"`
			URL                 string        `json:"url"`
			TotalTimeSpent      int           `json:"total_time_spent"`
			HumanTotalTimeSpent interface{}   `json:"human_total_time_spent"`
			HumanTimeEstimate   interface{}   `json:"human_time_estimate"`
			AssigneeIds         []interface{} `json:"assignee_ids"`
			Action              string        `json:"action"`
		} `json:"object_attributes"`
		Labels []struct {
			ID          int         `json:"id"`
			Title       string      `json:"title"`
			Color       string      `json:"color"`
			ProjectID   int         `json:"project_id"`
			CreatedAt   string      `json:"created_at"`
			UpdatedAt   string      `json:"updated_at"`
			Template    bool        `json:"template"`
			Description string      `json:"description"`
			Type        string      `json:"type"`
			GroupID     interface{} `json:"group_id"`
		} `json:"labels"`
		Changes struct {
			Labels struct {
				Previous []interface{} `json:"previous"`
				Current  []struct {
					ID          int         `json:"id"`
					Title       string      `json:"title"`
					Color       string      `json:"color"`
					ProjectID   int         `json:"project_id"`
					CreatedAt   string      `json:"created_at"`
					UpdatedAt   string      `json:"updated_at"`
					Template    bool        `json:"template"`
					Description string      `json:"description"`
					Type        string      `json:"type"`
					GroupID     interface{} `json:"group_id"`
				} `json:"current"`
			} `json:"labels"`
		} `json:"changes"`
		Repository struct {
			Name        string `json:"name"`
			URL         string `json:"url"`
			Description string `json:"description"`
			Homepage    string `json:"homepage"`
		} `json:"repository"`
	}

	pullRequestHook struct {
		ObjectKind string `json:"object_kind"`
		User       struct {
			Name      string `json:"name"`
			Username  string `json:"username"`
			AvatarURL string `json:"avatar_url"`
		} `json:"user"`
		Project struct {
			ID                int         `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      interface{} `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"project"`
		ObjectAttributes struct {
			AssigneeID                interface{} `json:"assignee_id"`
			AuthorID                  int         `json:"author_id"`
			CreatedAt                 string      `json:"created_at"`
			DeletedAt                 interface{} `json:"deleted_at"`
			Description               string      `json:"description"`
			HeadPipelineID            interface{} `json:"head_pipeline_id"`
			ID                        int         `json:"id"`
			Iid                       int         `json:"iid"`
			LastEditedAt              interface{} `json:"last_edited_at"`
			LastEditedByID            interface{} `json:"last_edited_by_id"`
			MergeCommitSha            interface{} `json:"merge_commit_sha"`
			MergeError                interface{} `json:"merge_error"`
			MergeParams               interface{} `json:"-"`
			MergeStatus               string      `json:"merge_status"`
			MergeUserID               interface{} `json:"merge_user_id"`
			MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
			MilestoneID               interface{} `json:"milestone_id"`
			SourceBranch              string      `json:"source_branch"`
			SourceProjectID           int         `json:"source_project_id"`
			State                     string      `json:"state"`
			TargetBranch              string      `json:"target_branch"`
			TargetProjectID           int         `json:"target_project_id"`
			TimeEstimate              int         `json:"time_estimate"`
			Title                     string      `json:"title"`
			UpdatedAt                 string      `json:"updated_at"`
			UpdatedByID               interface{} `json:"updated_by_id"`
			URL                       string      `json:"url"`
			Source                    struct {
				ID                int         `json:"id"`
				Name              string      `json:"name"`
				Description       string      `json:"description"`
				WebURL            string      `json:"web_url"`
				AvatarURL         interface{} `json:"avatar_url"`
				GitSSHURL         string      `json:"git_ssh_url"`
				GitHTTPURL        string      `json:"git_http_url"`
				Namespace         string      `json:"namespace"`
				VisibilityLevel   int         `json:"visibility_level"`
				PathWithNamespace string      `json:"path_with_namespace"`
				DefaultBranch     string      `json:"default_branch"`
				CiConfigPath      interface{} `json:"ci_config_path"`
				Homepage          string      `json:"homepage"`
				URL               string      `json:"url"`
				SSHURL            string      `json:"ssh_url"`
				HTTPURL           string      `json:"http_url"`
			} `json:"source"`
			Target struct {
				ID                int         `json:"id"`
				Name              string      `json:"name"`
				Description       string      `json:"description"`
				WebURL            string      `json:"web_url"`
				AvatarURL         interface{} `json:"avatar_url"`
				GitSSHURL         string      `json:"git_ssh_url"`
				GitHTTPURL        string      `json:"git_http_url"`
				Namespace         string      `json:"namespace"`
				VisibilityLevel   int         `json:"visibility_level"`
				PathWithNamespace string      `json:"path_with_namespace"`
				DefaultBranch     string      `json:"default_branch"`
				CiConfigPath      interface{} `json:"ci_config_path"`
				Homepage          string      `json:"homepage"`
				URL               string      `json:"url"`
				SSHURL            string      `json:"ssh_url"`
				HTTPURL           string      `json:"http_url"`
			} `json:"target"`
			LastCommit struct {
				ID        string `json:"id"`
				Message   string `json:"message"`
				Timestamp string `json:"timestamp"`
				URL       string `json:"url"`
				Author    struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"author"`
			} `json:"last_commit"`
			WorkInProgress      bool        `json:"work_in_progress"`
			TotalTimeSpent      int         `json:"total_time_spent"`
			HumanTotalTimeSpent interface{} `json:"human_total_time_spent"`
			HumanTimeEstimate   interface{} `json:"human_time_estimate"`
			Action              string      `json:"action"`
		} `json:"object_attributes"`
		Labels  []interface{} `json:"labels"`
		Changes struct {
			Draft struct {
				Previous null.Bool `json:"previous"`
				Current  null.Bool `json:"current"`
			} `json:"draft"`
		} `json:"changes"`
		Repository struct {
			Name        string `json:"name"`
			URL         string `json:"url"`
			Description string `json:"description"`
			Homepage    string `json:"homepage"`
		} `json:"repository"`
	}
)
