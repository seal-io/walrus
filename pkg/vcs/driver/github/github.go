package github

import (
	"fmt"
	"net/url"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/vcs/driver"
	"github.com/seal-io/walrus/pkg/vcs/options"
)

const (
	Driver = types.GitDriverGithub
)

// NewClient creates a new github client.
// Options connector token will overwrite options.WithToken in the client.
func NewClient(conn *model.Connector, opts ...options.ClientOption) (*scm.Client, error) {
	var (
		client *scm.Client
		err    error
	)

	rawURL, token, _, err := driver.ParseConnector(conn)
	if err != nil {
		return nil, err
	}

	if rawURL == "" {
		client = github.NewDefault()
	} else {
		client, err = github.New(rawURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create github client: %w", err)
		}
	}

	options.SetClientOptions(client, append(opts, options.WithToken(token))...)

	return client, nil
}

// NewClientFromURL creates a new github client from url.
func NewClientFromURL(rawURL string, opts ...options.ClientOption) (*scm.Client, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	var client *scm.Client

	switch u.Host {
	case "github.com":
		client = github.NewDefault()

	default:
		client, err = github.New(rawURL)
		if err != nil {
			return nil, err
		}
	}

	options.SetClientOptions(client, opts...)

	return client, nil
}
