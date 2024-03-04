// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/drone/go-scm/scm"
)

// TODO(bradrydzewski) default repository branch is missing in push webhook payloads
// TODO(bradrydzewski) default repository branch is missing in pr webhook payloads
// TODO(bradrydzewski) default repository is_private is missing in pr webhook payloads

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
	switch req.Header.Get("x-event-key") {
	case "repo:push":
		hook, err = s.parsePushHook(data)
	case "pullrequest:created":
		hook, err = s.parsePullRequestHook(data)
	case "pullrequest:updated":
		hook, err = s.parsePullRequestHook(data)
		if hook != nil {
			hook.(*scm.PullRequestHook).Action = scm.ActionSync
		}
	case "pullrequest:fulfilled":
		hook, err = s.parsePullRequestHook(data)
		if hook != nil {
			hook.(*scm.PullRequestHook).Action = scm.ActionMerge
		}
	case "pullrequest:rejected":
		hook, err = s.parsePullRequestHook(data)
		if hook != nil {
			hook.(*scm.PullRequestHook).Action = scm.ActionClose
		}
	case "pullrequest:comment_created":
		hook, err = s.parsePullRequestCommentHook(data)
		if hook != nil {
			hook.(*scm.IssueCommentHook).Action = scm.ActionCreate
		}
	case "pullrequest:comment_updated":
		// Bitbucket PR Comment Update is unreliable and does not send events
		// most of the time https://github.com/iterative/cml/issues/817
		hook, err = s.parsePullRequestCommentHook(data)
		if hook != nil {
			hook.(*scm.IssueCommentHook).Action = scm.ActionEdit
		}
	case "pullrequest:comment_deleted":
		hook, err = s.parsePullRequestCommentHook(data)
		if hook != nil {
			hook.(*scm.IssueCommentHook).Action = scm.ActionDelete
		}
	}
	if err != nil {
		return nil, err
	}
	if hook == nil {
		return nil, nil
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

	if req.FormValue("secret") != key {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePullRequestCommentHook(data []byte) (scm.Webhook, error) {
	dst := new(prCommentHook)
	err := json.Unmarshal(data, dst)
	return convertPrCommentHook(dst), err
}

func (s *webhookService) parsePushHook(data []byte) (scm.Webhook, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	if err != nil {
		return nil, err
	}
	if len(dst.Push.Changes) == 0 {
		return nil, errors.New("Push hook has empty changeset")
	}
	change := dst.Push.Changes[0]
	switch {
	// case change.New.Type == "branch" && change.Created:
	// 	return convertBranchCreateHook(dst), nil
	case change.Old.Type == "branch" && change.Closed:
		return convertBranchDeleteHook(dst), nil
	// case change.New.Type == "tag" && change.Created:
	// 	return convertTagCreateHook(dst), nil
	case change.Old.Type == "tag" && change.Closed:
		return convertTagDeleteHook(dst), nil
	default:
		return convertPushHook(dst), err
	}
}

func (s *webhookService) parsePullRequestHook(data []byte) (*scm.PullRequestHook, error) {
	dst := new(webhook)
	err := json.Unmarshal(data, dst)
	if err != nil {
		return nil, err
	}

	switch {
	default:
		return convertPullRequestHook(dst), err
	}
}

//
// native data structures
//

type (
	pushHook struct {
		Push struct {
			Changes []struct {
				Forced bool `json:"forced"`
				Old    struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Links struct {
						Commits struct {
							Href string `json:"href"`
						} `json:"commits"`
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
					Target struct {
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
						Author struct {
							Raw  string `json:"raw"`
							Type string `json:"type"`
							User struct {
								Username    string `json:"username"`
								DisplayName string `json:"display_name"`
								AccountID   string `json:"account_id"`
								Links       struct {
									Self struct {
										Href string `json:"href"`
									} `json:"self"`
									HTML struct {
										Href string `json:"href"`
									} `json:"html"`
									Avatar struct {
										Href string `json:"href"`
									} `json:"avatar"`
								} `json:"links"`
								Type string `json:"type"`
								UUID string `json:"uuid"`
							} `json:"user"`
						} `json:"author"`
						Summary struct {
							Raw    string `json:"raw"`
							Markup string `json:"markup"`
							HTML   string `json:"html"`
							Type   string `json:"type"`
						} `json:"summary"`
						Parents []interface{} `json:"parents"`
						Date    time.Time     `json:"date"`
						Message string        `json:"message"`
						Type    string        `json:"type"`
					} `json:"target"`
				} `json:"old"`
				Links struct {
					Commits struct {
						Href string `json:"href"`
					} `json:"commits"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
					Diff struct {
						Href string `json:"href"`
					} `json:"diff"`
				} `json:"links"`
				Truncated bool `json:"truncated"`
				Commits   []struct {
					Hash  string `json:"hash"`
					Links struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						Comments struct {
							Href string `json:"href"`
						} `json:"comments"`
						Patch struct {
							Href string `json:"href"`
						} `json:"patch"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
						Diff struct {
							Href string `json:"href"`
						} `json:"diff"`
						Approve struct {
							Href string `json:"href"`
						} `json:"approve"`
						Statuses struct {
							Href string `json:"href"`
						} `json:"statuses"`
					} `json:"links"`
					Author struct {
						Raw  string `json:"raw"`
						Type string `json:"type"`
						User struct {
							Username    string `json:"username"`
							DisplayName string `json:"display_name"`
							AccountID   string `json:"account_id"`
							Links       struct {
								Self struct {
									Href string `json:"href"`
								} `json:"self"`
								HTML struct {
									Href string `json:"href"`
								} `json:"html"`
								Avatar struct {
									Href string `json:"href"`
								} `json:"avatar"`
							} `json:"links"`
							Type string `json:"type"`
							UUID string `json:"uuid"`
						} `json:"user"`
					} `json:"author"`
					Summary struct {
						Raw    string `json:"raw"`
						Markup string `json:"markup"`
						HTML   string `json:"html"`
						Type   string `json:"type"`
					} `json:"summary"`
					Parents []struct {
						Type  string `json:"type"`
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
					} `json:"parents"`
					Date    time.Time `json:"date"`
					Message string    `json:"message"`
					Type    string    `json:"type"`
				} `json:"commits"`
				Created bool `json:"created"`
				Closed  bool `json:"closed"`
				New     struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Links struct {
						Commits struct {
							Href string `json:"href"`
						} `json:"commits"`
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
					Target struct {
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
						Author struct {
							Raw  string `json:"raw"`
							Type string `json:"type"`
							User struct {
								Username    string `json:"username"`
								DisplayName string `json:"display_name"`
								AccountID   string `json:"account_id"`
								Links       struct {
									Self struct {
										Href string `json:"href"`
									} `json:"self"`
									HTML struct {
										Href string `json:"href"`
									} `json:"html"`
									Avatar struct {
										Href string `json:"href"`
									} `json:"avatar"`
								} `json:"links"`
								Type string `json:"type"`
								UUID string `json:"uuid"`
							} `json:"user"`
						} `json:"author"`
						Summary struct {
							Raw    string `json:"raw"`
							Markup string `json:"markup"`
							HTML   string `json:"html"`
							Type   string `json:"type"`
						} `json:"summary"`
						Parents []struct {
							Type  string `json:"type"`
							Hash  string `json:"hash"`
							Links struct {
								Self struct {
									Href string `json:"href"`
								} `json:"self"`
								HTML struct {
									Href string `json:"href"`
								} `json:"html"`
							} `json:"links"`
						} `json:"parents"`
						Date    time.Time `json:"date"`
						Message string    `json:"message"`
						Type    string    `json:"type"`
					} `json:"target"`
				} `json:"new"`
			} `json:"changes"`
		} `json:"push"`
		Repository webhookRepository `json:"repository"`
		Actor      webhookActor      `json:"actor"`
	}

	webhook struct {
		PullRequest pr                `json:"pullrequest"`
		Repository  webhookRepository `json:"repository"`
		Actor       webhookActor      `json:"actor"`
	}

	webhookRepository struct {
		Scm   string `json:"scm"`
		Name  string `json:"name"`
		Links struct {
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
		FullName string `json:"full_name"`
		Owner    struct {
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			AccountID   string `json:"account_id"`
			Links       struct {
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
			} `json:"links"`
			UUID string `json:"uuid"`
		} `json:"owner"`
		IsPrivate bool   `json:"is_private"`
		UUID      string `json:"uuid"`
	}

	webhookActor struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		AccountID   string `json:"account_id"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		UUID string `json:"uuid"`
	}

	prCommentInput struct {
		Content struct {
			Raw string `json:"raw"`
		} `json:"content"`
	}

	prComment struct {
		Links struct {
			Self link `json:"self"`
			HTML link `json:"html"`
		} `json:"links"`
		Deleted     bool `json:"deleted"`
		PullRequest struct {
			Type  string `json:"type"`
			ID    int    `json:"id"`
			Links struct {
				Self link `json:"self"`
				HTML link `json:"html"`
			} `json:"links"`
			Title string `json:"title"`
		}
		Content struct {
			Raw    string `json:"raw"`
			Markup string `json:"markup"`
			Html   string `json:"html"`
			Type   string `json:"type"`
		}
		CreatedOn time.Time         `json:"created_on"`
		User      prCommentHookUser `json:"user"`
		UpdatedOn time.Time         `json:"updated_on"`
		Type      string            `json:"type"`
		ID        int               `json:"id"`
	}

	prCommentHookRepo struct {
		Scm     string `json:"scm"`
		Website string `json:"website"`
		UUID    string `json:"uuid"`
		Links   struct {
			Self   link `json:"self"`
			HTML   link `json:"html"`
			Avatar link `json:"avatar"`
		} `json:"links"`
		Project struct {
			Links struct {
				Self   link `json:"self"`
				HTML   link `json:"html"`
				Avatar link `json:"avatar"`
			} `json:"links"`
			Type string `json:"type"`
			Name string `json:"name"`
			Key  string `json:"key"`
			UUID string `json:"uuid"`
		} `json:"project"`
		FullName  string            `json:"full_name"`
		Owner     prCommentHookUser `json:"owner"`
		Workspace struct {
			Slug  string `json:"slugg"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			Links struct {
				Self   link `json:"self"`
				HTML   link `json:"html"`
				Avatar link `json:"avatar"`
			} `json:"links"`
			UUID string `json:"uuid"`
		} `json:"workspace"`
		Type      string `json:"type"`
		IsPrivate bool   `json:"is_private"`
		Name      string `json:"name"`
	}

	prCommentHookUser struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self   link `json:"self"`
			HTML   link `json:"html"`
			Avatar link `json:"avatar"`
		} `json:"links"`
		Type      string `json:"type"`
		Nickname  string `json:"nickname"`
		AccountID string `json:"account_id"`
	}

	prCommentHookPullRequest struct {
		Rendered struct {
			Description struct {
				Raw    string `json:"raw"`
				Markup string `json:"markup"`
				Html   string `json:"html"`
				Type   string `json:"type"`
			} `json:"description"`
			Title struct {
				Raw    string `json:"raw"`
				Markup string `json:"markup"`
				Html   string `json:"html"`
				Type   string `json:"type"`
			} `json:"title"`
		} `json:"rendered"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Links       struct {
			Decline        link `json:"decline"`
			Diffstat       link `json:"diffstat"`
			Commits        link `json:"commits"`
			Self           link `json:"self"`
			Comments       link `json:"comments"`
			Merge          link `json:"merge"`
			Html           link `json:"html"`
			Activity       link `json:"activity"`
			RequestChanges link `json:"request-changes"`
			Diff           link `json:"diff"`
			Approve        link `json:"approve"`
			Statuses       link `json:"statuses"`
		} `json:"links"`
		Title             string        `json:"title"`
		CloseSourceBranch bool          `json:"close_source_branch"`
		Reviewers         []interface{} `json:"reviewers"`
		ID                int           `json:"id"`
		Destination       struct {
			Commit struct {
				Hash  string `json:"hash"`
				Type  string `json:"type"`
				Links struct {
					Self link `json:"self"`
					HTML link `json:"html"`
				} `json:"links"`
			}
			Repository struct {
				Links struct {
					Self   link `json:"self"`
					HTML   link `json:"html"`
					Avatar link `json:"avatar"`
				} `json:"links"`
				Type     string `json:"type"`
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				UUID     string `json:"uuid"`
			} `json:"repository"`
			Branch struct {
				Name string `json:"name"`
			} `json:"branch"`
		} `json:"destination"`
		CreatedOn time.Time `json:"created_on"`
		Summary   struct {
			Raw    string `json:"raw"`
			Markup string `json:"markup"`
			Html   string `json:"html"`
			Type   string `json:"type"`
		} `json:"summary"`
		Source struct {
			Commit struct {
				Hash  string `json:"hash"`
				Type  string `json:"type"`
				Links struct {
					Self link `json:"self"`
					HTML link `json:"html"`
				} `json:"links"`
			}
			Repository struct {
				Links struct {
					Self   link `json:"self"`
					HTML   link `json:"html"`
					Avatar link `json:"avatar"`
				} `json:"links"`
				Type     string `json:"type"`
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				UUID     string `json:"uuid"`
			} `json:"repository"`
			Branch struct {
				Name string `json:"name"`
			} `json:"branch"`
		} `json:"source"`
		CommentCount int               `json:"comment_count"`
		State        string            `json:"state"`
		TaskCount    int               `json:"task_count"`
		Participants []interface{}     `json:"participants"`
		Reason       string            `json:"reason"`
		UpdatedOn    time.Time         `json:"updated_on"`
		Author       prCommentHookUser `json:"author"`
		MergeCommit  interface{}       `json:"merge_commit"`
		ClosedBy     interface{}       `json:"closed_by"`
	}

	prCommentHook struct {
		Comment     prComment                `json:"comment"`
		PullRequest prCommentHookPullRequest `json:"pullRequest"`
		Repository  prCommentHookRepo        `json:"repository"`
		Actor       prCommentHookUser        `json:"actor"`
	}
)

//
// push hooks
//

func convertPushHook(src *pushHook) *scm.PushHook {
	change := src.Push.Changes[0]
	var commits []scm.Commit
	for _, c := range change.Commits {
		commits = append(commits,
			scm.Commit{
				Sha:     c.Hash,
				Message: c.Message,
				Link:    c.Links.HTML.Href,
				Author: scm.Signature{
					Login:  c.Author.User.Username,
					Email:  extractEmail(c.Author.Raw),
					Name:   c.Author.User.DisplayName,
					Avatar: c.Author.User.Links.Avatar.Href,
					Date:   c.Date,
				},
				Committer: scm.Signature{
					Login:  c.Author.User.Username,
					Email:  extractEmail(c.Author.Raw),
					Name:   c.Author.User.DisplayName,
					Avatar: c.Author.User.Links.Avatar.Href,
					Date:   c.Date,
				},
			})
	}
	namespace, name := scm.Split(src.Repository.FullName)
	dst := &scm.PushHook{
		Ref:    scm.ExpandRef(change.New.Name, "refs/heads/"),
		Before: change.Old.Target.Hash,
		After:  change.New.Target.Hash,
		Commit: scm.Commit{
			Sha:     change.New.Target.Hash,
			Message: change.New.Target.Message,
			Link:    change.New.Target.Links.HTML.Href,
			Author: scm.Signature{
				Login:  change.New.Target.Author.User.Username,
				Email:  extractEmail(change.New.Target.Author.Raw),
				Name:   change.New.Target.Author.User.DisplayName,
				Avatar: change.New.Target.Author.User.Links.Avatar.Href,
				Date:   change.New.Target.Date,
			},
			Committer: scm.Signature{
				Login:  change.New.Target.Author.User.Username,
				Email:  extractEmail(change.New.Target.Author.Raw),
				Name:   change.New.Target.Author.User.DisplayName,
				Avatar: change.New.Target.Author.User.Links.Avatar.Href,
				Date:   change.New.Target.Date,
			},
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: namespace,
			Name:      name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			ID:     src.Actor.UUID,
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
		Commits: commits,
	}
	if change.New.Type == "tag" {
		dst.Ref = scm.ExpandRef(change.New.Name, "refs/tags/")
	}
	return dst
}

func convertBranchCreateHook(src *pushHook) *scm.BranchHook {
	namespace, name := scm.Split(src.Repository.FullName)
	change := src.Push.Changes[0].New
	action := scm.ActionCreate
	return &scm.BranchHook{
		Action: action,
		Ref: scm.Reference{
			Name: change.Name,
			Sha:  change.Target.Hash,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: namespace,
			Name:      name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			ID:     src.Actor.UUID,
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

func convertBranchDeleteHook(src *pushHook) *scm.BranchHook {
	namespace, name := scm.Split(src.Repository.FullName)
	change := src.Push.Changes[0].Old
	action := scm.ActionDelete
	return &scm.BranchHook{
		Action: action,
		Ref: scm.Reference{
			Name: change.Name,
			Sha:  change.Target.Hash,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: namespace,
			Name:      name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			ID:     src.Actor.UUID,
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

func convertTagCreateHook(src *pushHook) *scm.TagHook {
	namespace, name := scm.Split(src.Repository.FullName)
	change := src.Push.Changes[0].New
	action := scm.ActionCreate
	return &scm.TagHook{
		Action: action,
		Ref: scm.Reference{
			Name: change.Name,
			Sha:  change.Target.Hash,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: namespace,
			Name:      name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			ID:     src.Actor.UUID,
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

func convertTagDeleteHook(src *pushHook) *scm.TagHook {
	namespace, name := scm.Split(src.Repository.FullName)
	change := src.Push.Changes[0].Old
	action := scm.ActionDelete
	return &scm.TagHook{
		Action: action,
		Ref: scm.Reference{
			Name: change.Name,
			Sha:  change.Target.Hash,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: namespace,
			Name:      name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			ID:     src.Actor.UUID,
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

//
// pull request hooks
//

func convertPullRequestHook(src *webhook) *scm.PullRequestHook {
	namespace, name := scm.Split(src.Repository.FullName)
	return &scm.PullRequestHook{
		Action: scm.ActionOpen,
		PullRequest: scm.PullRequest{
			Number: src.PullRequest.ID,
			Title:  src.PullRequest.Title,
			Body:   src.PullRequest.Description,
			Sha:    src.PullRequest.Source.Commit.Hash,
			Merge:  src.PullRequest.MergeCommit.Hash,
			Ref:    fmt.Sprintf("refs/pull-requests/%d/from", src.PullRequest.ID),
			Source: src.PullRequest.Source.Branch.Name,
			Target: src.PullRequest.Destination.Branch.Name,
			Fork:   src.PullRequest.Source.Repository.FullName,
			Link:   src.PullRequest.Links.HTML.Href,
			Closed: src.PullRequest.State != "OPEN",
			Merged: src.PullRequest.State == "MERGED",
			Author: scm.User{
				Login:  src.PullRequest.Author.Username,
				Name:   src.PullRequest.Author.DisplayName,
				Avatar: src.PullRequest.Author.Links.Avatar.Href,
			},
			Created: src.PullRequest.CreatedOn,
			Updated: src.PullRequest.UpdatedOn,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: namespace,
			Name:      name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			ID:     src.Actor.UUID,
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

func convertPrCommentHook(src *prCommentHook) *scm.IssueCommentHook {
	namespace, _ := scm.Split(src.Repository.FullName)
	dst := scm.IssueCommentHook{
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: namespace,
			Name:      src.Repository.Name,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
			Private:   src.Repository.IsPrivate,
		},
		Issue: scm.Issue{
			Number: src.PullRequest.ID,
			Title:  src.PullRequest.Title,
			Body:   src.PullRequest.Description,
			Link:   src.PullRequest.Links.Html.Href,
			Author: scm.User{
				Login:  src.PullRequest.Author.Username,
				Name:   src.PullRequest.Author.DisplayName,
				Avatar: src.PullRequest.Author.Links.Avatar.Href,
			},
			PullRequest: scm.PullRequest{
				Number: src.PullRequest.ID,
				Title:  src.PullRequest.Title,
				Body:   src.PullRequest.Description,
				Sha:    src.PullRequest.Source.Commit.Hash,
				// Bitbucket does not support PR Refs: https://jira.atlassian.com/browse/BCLOUD-5814
				Ref:    fmt.Sprintf("refs/pull-requests/%d/from", src.PullRequest.ID),
				Source: src.PullRequest.Source.Branch.Name,
				Target: src.PullRequest.Destination.Branch.Name,
				Fork:   src.PullRequest.Source.Repository.FullName,
				Link:   src.PullRequest.Links.Html.Href,
				Closed: src.PullRequest.State != "OPEN",
				Merged: src.PullRequest.State == "MERGED",
				Author: scm.User{
					Login:  src.PullRequest.Author.Username,
					Name:   src.PullRequest.Author.DisplayName,
					Avatar: src.PullRequest.Author.Links.Avatar.Href,
				},
				Created: src.PullRequest.CreatedOn,
				Updated: src.PullRequest.UpdatedOn,
			},
			Created: src.PullRequest.CreatedOn,
			Updated: src.PullRequest.UpdatedOn,
		},
		Comment: scm.Comment{
			ID:   src.Comment.ID,
			Body: src.Comment.Content.Raw,
			Author: scm.User{
				ID:     src.Comment.User.UUID,
				Login:  src.Comment.User.Username,
				Name:   src.Comment.User.DisplayName,
				Avatar: src.Comment.User.Links.Avatar.Href,
			},
			Created: src.Comment.CreatedOn,
			Updated: src.Comment.UpdatedOn,
		},
		Sender: scm.User{
			ID:     src.Actor.UUID,
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
	return &dst
}
