package harness

import (
	"context"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/null"
)

type releaseService struct {
	client *wrapper
}

func (s *releaseService) Find(ctx context.Context, repo string, id int) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) FindByTag(ctx context.Context, repo string, tag string) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) List(ctx context.Context, repo string, opts scm.ReleaseListOptions) ([]*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) Create(ctx context.Context, repo string, input *scm.ReleaseInput) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) Delete(ctx context.Context, repo string, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *releaseService) DeleteByTag(ctx context.Context, repo string, tag string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *releaseService) Update(ctx context.Context, repo string, id int, input *scm.ReleaseInput) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *releaseService) UpdateByTag(ctx context.Context, repo string, tag string, input *scm.ReleaseInput) (*scm.Release, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type ReleaseInput struct {
	TagName      string `json:"tag_name"`
	Target       string `json:"target_commitish"`
	Title        string `json:"name"`
	Note         string `json:"body"`
	IsDraft      bool   `json:"draft"`
	IsPrerelease bool   `json:"prerelease"`
}

// release represents a repository release
type release struct {
	ID           int64         `json:"id"`
	TagName      string        `json:"tag_name"`
	Target       string        `json:"target_commitish"`
	Title        string        `json:"name"`
	Note         string        `json:"body"`
	URL          string        `json:"url"`
	HTMLURL      string        `json:"html_url"`
	TarURL       string        `json:"tarball_url"`
	ZipURL       string        `json:"zipball_url"`
	IsDraft      bool          `json:"draft"`
	IsPrerelease bool          `json:"prerelease"`
	CreatedAt    null.Time     `json:"created_at"`
	PublishedAt  null.Time     `json:"published_at"`
	Publisher    *string       `json:"author"`
	Attachments  []*Attachment `json:"assets"`
}

type Attachment struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Size          int64     `json:"size"`
	DownloadCount int64     `json:"download_count"`
	Created       null.Time `json:"created_at"`
	UUID          string    `json:"uuid"`
	DownloadURL   string    `json:"browser_download_url"`
}

func convertRelease(src *release) *scm.Release {
	return &scm.Release{
		ID:          int(src.ID),
		Title:       src.Title,
		Description: src.Note,
		Link:        convertAPIURLToHTMLURL(src.URL, src.TagName),
		Tag:         src.TagName,
		Commitish:   src.Target,
		Draft:       src.IsDraft,
		Prerelease:  src.IsPrerelease,
		Created:     src.CreatedAt.ValueOrZero(),
		Published:   src.PublishedAt.ValueOrZero(),
	}
}

func convertReleaseList(src []*release) []*scm.Release {
	var dst []*scm.Release
	for _, v := range src {
		dst = append(dst, convertRelease(v))
	}
	return dst
}

func releaseListOptionsToGiteaListOptions(in scm.ReleaseListOptions) ListOptions {
	return ListOptions{
		Page:     in.Page,
		PageSize: in.Size,
	}
}
