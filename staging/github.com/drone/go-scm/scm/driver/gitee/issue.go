// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
)

type issueService struct {
	client *wrapper
}

func (s *issueService) Find(ctx context.Context, repo string, number int) (*scm.Issue, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/issues/%s", repo, decodeNumber(number))
	out := new(issue)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertIssue(out), res, err
}

func (s *issueService) FindComment(ctx context.Context, repo string, number, id int) (*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/issues/comments/%d", repo, id)
	out := new(issueComment)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertIssueComment(out), res, err
}

func (s *issueService) List(ctx context.Context, repo string, opts scm.IssueListOptions) ([]*scm.Issue, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/issues?%s", repo, encodeIssueListOptions(opts))
	out := []*issue{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertIssueList(out), res, err
}

func (s *issueService) ListComments(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/issues/%s/comments?%s", repo, decodeNumber(number), encodeListOptions(opts))
	out := []*issueComment{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertIssueCommentList(out), res, err
}

func (s *issueService) Create(ctx context.Context, repo string, input *scm.IssueInput) (*scm.Issue, *scm.Response, error) {
	owner, repoName := scm.Split(repo)
	path := fmt.Sprintf("repos/%s/issues", owner)
	in := &issueInput{
		Repo:  repoName,
		Title: input.Title,
		Body:  input.Body,
	}
	out := new(issue)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertIssue(out), res, err
}

func (s *issueService) CreateComment(ctx context.Context, repo string, number int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/issues/%s/comments", repo, decodeNumber(number))
	in := &issueCommentInput{
		Body: input.Body,
	}
	out := new(issueComment)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertIssueComment(out), res, err
}

func (s *issueService) DeleteComment(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/issues/comments/%d", repo, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *issueService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	owner, repoName := scm.Split(repo)
	path := fmt.Sprintf("repos/%s/issues/%s", owner, decodeNumber(number))
	data := map[string]string{
		"repo":  repoName,
		"state": "closed",
	}
	out := new(issue)
	res, err := s.client.do(ctx, "PATCH", path, &data, out)
	return res, err
}

func (s *issueService) Lock(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *issueService) Unlock(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

type issue struct {
	ID               int              `json:"id"`
	URL              string           `json:"url"`
	RepositoryURL    string           `json:"repository_url"`
	LabelsURL        string           `json:"labels_url"`
	CommentsURL      string           `json:"comments_url"`
	HtmlURL          string           `json:"html_url"`
	ParentURL        interface{}      `json:"parent_url"`
	Number           string           `json:"number"`
	ParentID         int              `json:"parent_id"`
	Depth            int              `json:"depth"`
	State            string           `json:"state"`
	Title            string           `json:"title"`
	Body             string           `json:"body"`
	User             user             `json:"user"`
	Labels           []label          `json:"labels"`
	Assignee         assignee         `json:"assignee"`
	Collaborators    []interface{}    `json:"collaborators"`
	Repository       issueRepository  `json:"repository"`
	Milestone        milestone        `json:"milestone"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	PlanStartedAt    interface{}      `json:"plan_started_at"`
	Deadline         interface{}      `json:"deadline"`
	FinishedAt       interface{}      `json:"finished_at"`
	ScheduledTime    float64          `json:"scheduled_time"`
	Comments         int              `json:"comments"`
	Priority         int              `json:"priority"`
	IssueType        string           `json:"issue_type"`
	SecurityHole     bool             `json:"security_hole"`
	IssueState       string           `json:"issue_state"`
	IssueTypeDetail  issueTypeDetail  `json:"issue_type_detail"`
	IssueStateDetail issueStateDetail `json:"issue_state_detail"`
}
type issueRepository struct {
	ID                  int           `json:"id"`
	FullName            string        `json:"full_name"`
	HumanName           string        `json:"human_name"`
	URL                 string        `json:"url"`
	Namespace           namespace     `json:"namespace"`
	Path                string        `json:"path"`
	Name                string        `json:"name"`
	Owner               user          `json:"owner"`
	Assigner            assignee      `json:"assigner"`
	Description         string        `json:"description"`
	Private             bool          `json:"private"`
	Public              bool          `json:"public"`
	Internal            bool          `json:"internal"`
	Fork                bool          `json:"fork"`
	HTMLURL             string        `json:"html_url"`
	SSHURL              string        `json:"ssh_url"`
	ForksURL            string        `json:"forks_url"`
	KeysURL             string        `json:"keys_url"`
	CollaboratorsURL    string        `json:"collaborators_url"`
	HooksURL            string        `json:"hooks_url"`
	BranchesURL         string        `json:"branches_url"`
	TagsURL             string        `json:"tags_url"`
	BlobsURL            string        `json:"blobs_url"`
	StargazersURL       string        `json:"stargazers_url"`
	ContributorsURL     string        `json:"contributors_url"`
	CommitsURL          string        `json:"commits_url"`
	CommentsURL         string        `json:"comments_url"`
	IssueCommentURL     string        `json:"issue_comment_url"`
	IssuesURL           string        `json:"issues_url"`
	PullsURL            string        `json:"pulls_url"`
	MilestonesURL       string        `json:"milestones_url"`
	NotificationsURL    string        `json:"notifications_url"`
	LabelsURL           string        `json:"labels_url"`
	ReleasesURL         string        `json:"releases_url"`
	Recommend           bool          `json:"recommend"`
	Gvp                 bool          `json:"gvp"`
	Homepage            string        `json:"homepage"`
	Language            interface{}   `json:"language"`
	ForksCount          int           `json:"forks_count"`
	StargazersCount     int           `json:"stargazers_count"`
	WatchersCount       int           `json:"watchers_count"`
	DefaultBranch       string        `json:"default_branch"`
	OpenIssuesCount     int           `json:"open_issues_count"`
	HasIssues           bool          `json:"has_issues"`
	HasWiki             bool          `json:"has_wiki"`
	IssueComment        bool          `json:"issue_comment"`
	CanComment          bool          `json:"can_comment"`
	PullRequestsEnabled bool          `json:"pull_requests_enabled"`
	HasPage             bool          `json:"has_page"`
	License             interface{}   `json:"license"`
	Outsourced          bool          `json:"outsourced"`
	ProjectCreator      string        `json:"project_creator"`
	Members             []string      `json:"members"`
	PushedAt            time.Time     `json:"pushed_at"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	Parent              interface{}   `json:"parent"`
	Paas                interface{}   `json:"paas"`
	AssigneesNumber     int           `json:"assignees_number"`
	TestersNumber       int           `json:"testers_number"`
	Assignee            []assignee    `json:"assignee"`
	Testers             []tester      `json:"testers"`
	Status              string        `json:"status"`
	EmptyRepo           bool          `json:"empty_repo"`
	Programs            []interface{} `json:"programs"`
	Enterprise          interface{}   `json:"enterprise"`
}
type issueTypeDetail struct {
	ID        int         `json:"id"`
	Title     string      `json:"title"`
	Template  interface{} `json:"template"`
	Ident     string      `json:"ident"`
	Color     string      `json:"color"`
	IsSystem  bool        `json:"is_system"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
type issueStateDetail struct {
	ID        int         `json:"id"`
	Title     string      `json:"title"`
	Color     string      `json:"color"`
	Icon      string      `json:"icon"`
	Command   interface{} `json:"command"`
	Serial    int         `json:"serial"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type issueInput struct {
	Repo  string `json:"repo"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type issueComment struct {
	ID      int    `json:"id"`
	HTMLURL string `json:"html_url"`
	User    struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"user"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type issueCommentInput struct {
	Body string `json:"body"`
}

func convertIssueList(from []*issue) []*scm.Issue {
	to := []*scm.Issue{}
	for _, v := range from {
		to = append(to, convertIssue(v))
	}
	return to
}
func convertIssue(from *issue) *scm.Issue {
	return &scm.Issue{
		Number: encodeNumber(from.Number),
		Title:  from.Title,
		Body:   from.Body,
		Link:   from.HtmlURL,
		Labels: convertLabels(from),
		Closed: from.State == "closed",
		Author: scm.User{
			Login:  from.User.Login,
			Avatar: from.User.AvatarURL,
		},
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
	}
}

func convertIssueCommentList(from []*issueComment) []*scm.Comment {
	to := []*scm.Comment{}
	for _, v := range from {
		to = append(to, convertIssueComment(v))
	}
	return to
}
func convertIssueComment(from *issueComment) *scm.Comment {
	return &scm.Comment{
		ID:   from.ID,
		Body: from.Body,
		Author: scm.User{
			Login:  from.User.Login,
			Avatar: from.User.AvatarURL,
		},
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
	}
}

func convertLabels(from *issue) []string {
	var labels []string
	for _, label := range from.Labels {
		labels = append(labels, label.Name)
	}
	return labels
}

// The issue number of gitee consists of 6 uppercase letters or numbers.
// The ASCII of uppercase letters or numbers is between 48 and 90, so encoded issue number(max:9090909090) less than the maximum value of int.
func encodeNumber(giteeIssueNumber string) int {
	runes := []rune(giteeIssueNumber)
	encodedNumber := ""
	for i := 0; i < len(runes); i++ {
		encodedNumber += strconv.Itoa(int(runes[i]))
	}
	scmNumber, err := strconv.Atoi(encodedNumber)
	if err != nil {
		return 0
	}
	return scmNumber
}
func decodeNumber(scmIssueNumber int) string {
	issueNumberStr := strconv.Itoa(scmIssueNumber)
	giteeNumber := ""
	for i := 0; i < len(issueNumberStr)-1; i += 2 {
		numberStr, err := strconv.Atoi(issueNumberStr[i : i+2])
		if err != nil {
			return ""
		}
		giteeNumber += string(rune(numberStr))
	}
	return giteeNumber
}
