package gitee

import (
	"fmt"
	"net/url"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitee"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/vcs/driver"
	"github.com/seal-io/walrus/pkg/vcs/options"
)

const (
	Driver = types.GitDriverGitee
)

// NewClient creates a new gitee client.
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
		client = gitee.NewDefault()
	} else {
		client, err = gitee.New(rawURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create gitee client: %w", err)
		}
	}

	options.SetClientOptions(client, append(opts, options.WithToken(token))...)

	return client, nil
}

// NewClientFromURL creates a new gitee client from url.
func NewClientFromURL(rawURL string, opts ...options.ClientOption) (*scm.Client, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	var client *scm.Client

	switch u.Host {
	case "gitee.com":
		client = gitee.NewDefault()
	default:
		client, err = gitee.New(fmt.Sprintf("%s/api/v5", u.Scheme+"://"+u.Host))
		if err != nil {
			return nil, err
		}
	}

	options.SetClientOptions(client, opts...)

	return client, nil
}
