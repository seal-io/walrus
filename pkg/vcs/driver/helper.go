package driver

import (
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
)

// ParseConnector parses a connector and returns the driver type, raw url, and token for vcs driver.
func ParseConnector(conn *model.Connector) (rawURL, token, driverType string, err error) {
	var ok bool

	switch conn.ConfigVersion {
	default:
		return "", "", "", fmt.Errorf("unknown config version: %v", conn.ConfigVersion)
	case "v1":
	}

	driverType = conn.Type

	token, ok, err = conn.ConfigData["token"].GetString()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get token: %w", err)
	}

	if token == "" || !ok {
		return "", "", "", errors.New("token not found")
	}

	rawURL, ok, err = conn.ConfigData["base_url"].GetString()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get base url: %w", err)
	}

	if !ok {
		rawURL = ""
	}

	return
}
