package api

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/json"
	versionutil "github.com/seal-io/walrus/utils/version"
)

const (
	apiCacheFileName        = "api.json"
	apiVersionCacheFileName = "api-version.json"
)

var (
	apiCacheFile        = filepath.Join(config.GetConfigDir(), apiCacheFileName)
	apiVersionCacheFile = filepath.Join(config.GetConfigDir(), apiVersionCacheFileName)
)

// getAPIFromCache load api from cache.
func getAPIFromCache() *API {
	var api API

	data, err := os.ReadFile(apiCacheFile)
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
	// API cache.
	b, err := json.Marshal(api)
	if err != nil {
		return fmt.Errorf("error marshal API cache: %w", err)
	}

	if err = os.WriteFile(apiCacheFile, b, 0o600); err != nil {
		return fmt.Errorf("error write API to cache: %w", err)
	}

	// API version cache.
	v := Version{
		Version:      api.Version.Version,
		GitCommit:    api.Version.GitCommit,
		IsDevVersion: versionutil.IsDevVersionWith(api.Version.Version),
	}

	vb, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Errorf("error marshal API version cache: %w", err)
	}

	if err = os.WriteFile(apiVersionCacheFile, vb, 0o600); err != nil {
		return fmt.Errorf("error write API version to cache: %w", err)
	}

	return nil
}

// GetAPIVersionFromCache load api version from cache.
func GetAPIVersionFromCache() *Version {
	var v Version

	data, err := os.ReadFile(apiVersionCacheFile)
	if err == nil {
		err = json.Unmarshal(data, &v)
		if err == nil {
			return &v
		}
	}

	return nil
}
