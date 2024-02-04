package vcs

import (
	"fmt"
	"net/url"

	"github.com/drone/go-scm/scm"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/vcs/driver/gitee"
	"github.com/seal-io/walrus/pkg/vcs/driver/github"
	"github.com/seal-io/walrus/pkg/vcs/driver/gitlab"
	"github.com/seal-io/walrus/pkg/vcs/options"
)

func NewClient(conn *model.Connector, opts ...options.ClientOption) (*scm.Client, error) {
	var (
		client *scm.Client
		err    error
	)

	switch conn.Type {
	case github.Driver:
		client, err = github.NewClient(conn, opts...)
		if err != nil {
			return nil, err
		}
	case gitlab.Driver:
		client, err = gitlab.NewClient(conn, opts...)
		if err != nil {
			return nil, err
		}
	case gitee.Driver:
		client, err = gitee.NewClient(conn, opts...)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported SCM driver %q", conn.Type)
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewClientFromURL(driver, rawURL string, opts ...options.ClientOption) (*scm.Client, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	// TODO support reverse proxy for self-hosted server.
	server := u.Scheme + "://" + u.Host

	switch driver {
	case github.Driver:
		return github.NewClientFromURL(server, opts...)
	case gitlab.Driver:
		return gitlab.NewClientFromURL(server, opts...)
	case gitee.Driver:
		return gitee.NewClientFromURL(server, opts...)
	}

	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("unsupported SCM driver %q", driver)
}
