package vcs

import (
	"fmt"

	"github.com/drone/go-scm/scm"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/vcs/driver/github"
	"github.com/seal-io/seal/pkg/vcs/driver/gitlab"
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
