package harness

import (
	"context"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/null"
)

type milestoneService struct {
	client *wrapper
}

func (s *milestoneService) Find(ctx context.Context, repo string, id int) (*scm.Milestone, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported

}

func (s *milestoneService) List(ctx context.Context, repo string, opts scm.MilestoneListOptions) ([]*scm.Milestone, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported

}

func (s *milestoneService) Create(ctx context.Context, repo string, input *scm.MilestoneInput) (*scm.Milestone, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported

}

func (s *milestoneService) Delete(ctx context.Context, repo string, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported

}

func (s *milestoneService) Update(ctx context.Context, repo string, id int, input *scm.MilestoneInput) (*scm.Milestone, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// stateType issue state type
type stateType string

const (
	// stateOpen pr/issue is open
	stateOpen stateType = "open"
	// stateClosed pr/issue is closed
	stateClosed stateType = "closed"
	// stateAll is all
	stateAll stateType = "all"
)

type milestone struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        stateType `json:"state"`
	OpenIssues   int       `json:"open_issues"`
	ClosedIssues int       `json:"closed_issues"`
	Created      null.Time `json:"created_at"`
	Updated      null.Time `json:"updated_at"`
	Closed       null.Time `json:"closed_at"`
	Deadline     null.Time `json:"due_on"`
}

type milestoneInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       stateType `json:"state"`
	Deadline    null.Time `json:"due_on"`
}

func convertMilestoneList(src []*milestone) []*scm.Milestone {
	var dst []*scm.Milestone
	for _, v := range src {
		dst = append(dst, convertMilestone(v))
	}
	return dst
}

func convertMilestone(src *milestone) *scm.Milestone {
	if src == nil || src.Deadline.IsZero() {
		return nil
	}
	return &scm.Milestone{
		Number:      int(src.ID),
		ID:          int(src.ID),
		Title:       src.Title,
		Description: src.Description,
		State:       string(src.State),
		DueDate:     src.Deadline.ValueOrZero(),
	}
}
