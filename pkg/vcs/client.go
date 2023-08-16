package vcs

import (
	"fmt"

	"github.com/drone/go-scm/scm"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/vcs/driver/github"
	"github.com/seal-io/walrus/pkg/vcs/driver/gitlab"
)

func NewClient(conn *model.Connector) (*scm.Client, error) {
	var (
		client *scm.Client
		err    error
	)

	switch conn.Type {
	case github.Driver:
		client, err = github.NewClient(conn)
		if err != nil {
			return nil, err
		}
	case gitlab.Driver:
		client, err = gitlab.NewClient(conn)
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

func NewClientFromURL(driver, url, token string) (*scm.Client, error) {
	switch driver {
	case github.Driver:
		return github.NewClientFromURL(url, token)
	case gitlab.Driver:
		return gitlab.NewClientFromURL(url, token)
	}

	return nil, fmt.Errorf("unsupported SCM driver %q", driver)
}
