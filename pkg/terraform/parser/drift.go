package parser

import (
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/json"
)

// ParseDriftOutput parse terraform output, get resource change info.
func ParseDriftOutput(input string) (*types.ResourceDrift, error) {
	resourceDrift := &types.ResourceDrift{}

	if err := json.Unmarshal([]byte(input), resourceDrift); err != nil {
		return nil, err
	}

	return resourceDrift, nil
}
