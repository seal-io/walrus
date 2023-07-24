package telemetry

import (
	"fmt"
	"os"

	"github.com/posthog/posthog-go"
)

var (
	APIKey        = ""
	APIKeyEnvName = "SEAL_TELEMETRY_API_KEY"
	Endpoint      = "https://app.posthog.com"
)

func init() {
	if APIKey == "" {
		APIKey = os.Getenv(APIKeyEnvName)
	}
}

func PhClient() (posthog.Client, error) {
	if APIKey == "" {
		return nil, fmt.Errorf("%s is not set", APIKeyEnvName)
	}

	phCli, err := posthog.NewWithConfig(
		APIKey,
		posthog.Config{
			Endpoint: Endpoint,
			Logger:   wrapLogger,
		})
	if err != nil {
		return nil, err
	}

	return phCli, nil
}
