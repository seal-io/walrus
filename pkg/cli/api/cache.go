package api

import (
	"fmt"
	"os"
	"path"

	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/json"
)

const (
	cacheFileName = "apis.json"
)

var CacheFile = path.Join(config.GetConfigDir(), cacheFileName)

// getAPIFromCache load api from cache.
func getAPIFromCache() *API {
	var api API

	data, err := os.ReadFile(CacheFile)
	if err == nil {
		err = json.Unmarshal(data, &api)
		if err == nil {
			return &api
		}
	}

	return nil
}

// setAPIToCache set api to cache.
func setAPIToCache(api *API) error {
	b, err := json.Marshal(api)
	if err != nil {
		return fmt.Errorf("error marshal API cache: %w", err)
	}

	if err = os.WriteFile(CacheFile, b, 0o600); err != nil {
		return fmt.Errorf("error write API to cache: %w", err)
	}

	return nil
}
