package github

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/transport"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/maps"
)

const (
	Driver = "Github"
)

func NewClient(conn *model.Connector) (*scm.Client, error) {
	var (
		client *scm.Client
		err    error
	)

	switch conn.ConfigVersion {
	default:
		return nil, fmt.Errorf("unknown config version: %v", conn.ConfigVersion)
	case "v1":
	}

	url := maps.GetString(conn.ConfigData, "base_url")
	if url == "" {
		client = github.NewDefault()
	} else {
		client, err = github.New(url)
		if err != nil {
			return nil, fmt.Errorf("failed to create github client: %w", err)
		}
	}

	token := maps.GetString(conn.ConfigData, "token")
	if token == "" {
		return nil, errors.New("token not found")
	}
	client.Client = &http.Client{
		Timeout: time.Second * 15,
		Transport: &transport.BearerToken{
			Token: token,
		},
	}
	return client, nil
}
