package driver

import (
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/property"
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

	token, ok, err = property.GetString(conn.ConfigData["token"].Value)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get token: %w", err)
	}

	if token == "" || !ok {
		return "", "", "", errors.New("token not found")
	}

	rawURL, ok, err = property.GetString(conn.ConfigData["base_url"].Value)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get base url: %w", err)
	}

	if !ok {
		rawURL = ""
	}

	return
}
