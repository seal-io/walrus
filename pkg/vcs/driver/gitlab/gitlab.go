package gitlab

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitlab"
	"github.com/drone/go-scm/scm/transport"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

const (
	Driver = types.GitDriverGitlab
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

	url, ok, err := conn.ConfigData["base_url"].GetString()
	if err != nil {
		return nil, fmt.Errorf("failed to get base url: %w", err)
	}

	if url == "" || !ok {
		client = gitlab.NewDefault()
	} else {
		client, err = gitlab.New(url)
		if err != nil {
			return nil, fmt.Errorf("failed to create github client: %w", err)
		}
	}

	token, ok, err := conn.ConfigData["token"].GetString()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	if token == "" || !ok {
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

func NewClientFromURL(rawURL, token string) (*scm.Client, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	client, err := gitlab.New(u.Scheme + "://" + u.Host)
	if err != nil {
		return nil, err
	}

	client.Client = &http.Client{
		Timeout: time.Second * 15,
	}

	if token != "" {
		client.Client.Transport = &transport.BearerToken{
			Token: token,
		}
	}

	return client, nil
}
