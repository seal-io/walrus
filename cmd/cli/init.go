package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/seal-io/seal/pkg/cli/api"
	"github.com/seal-io/seal/pkg/cli/config"
	"github.com/seal-io/seal/utils/json"
)

const (
	cacheFileName  = "apis.json"
	configFileName = "config.json"
)

// root represent the root command.
var root *cobra.Command

// Init define init steps.
func Init() error {
	err := initConfig()
	if err != nil {
		return err
	}

	root = NewRootCmd()

	return nil
}

// initConfig init config dir and load server context from cache.
func initConfig() error {
	// Config dir.
	configDir := getConfigDir()

	err := os.MkdirAll(configDir, 0o700)
	if err != nil {
		return err
	}

	// Config file.
	filename := path.Join(configDir, configFileName)

	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(filename, []byte("{}"), 0o600)
		}

		return fmt.Errorf("error stat config file %s: %w", filename, err)
	}

	sc, err := getServerContextFromCache()
	if err != nil {
		return err
	}

	serverConfig.ServerContext = *sc

	return nil
}

// getConfigDir get config dir.
func getConfigDir() string {
	userHomeDir := func() string {
		if runtime.GOOS == "windows" {
			home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}

			return home
		}

		return os.Getenv("HOME")
	}

	return path.Join(userHomeDir(), "."+cliName)
}

// getAPIFromCache load api from cache.
func getAPIFromCache() *api.API {
	var (
		api      api.API
		filename = path.Join(getConfigDir(), cacheFileName)
	)

	data, err := os.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, &api)
		if err == nil {
			return &api
		}
	}

	return nil
}

// setAPIToCache set api to cache.
func setAPIToCache(api *api.API) error {
	b, err := json.Marshal(api)
	if err != nil {
		return fmt.Errorf("error marshal API cache: %w", err)
	}

	filename := path.Join(getConfigDir(), cacheFileName)
	if err = os.WriteFile(filename, b, 0o600); err != nil {
		return fmt.Errorf("error write API to cache: %w", err)
	}

	return nil
}

// setServerContextToCache set context to cache.
func setServerContextToCache(s config.ServerContext) error {
	filename := path.Join(getConfigDir(), configFileName)
	content, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return fmt.Errorf("error decode config file %s: %w", filename, err)
	}

	err = os.WriteFile(filename, content, 0o600)
	if err != nil {
		return fmt.Errorf("error save configure: %w", err)
	}

	return nil
}

// getServerContextFromCache load context from cache.
func getServerContextFromCache() (*config.ServerContext, error) {
	filename := path.Join(getConfigDir(), configFileName)
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error read config file %s: %w", filename, err)
	}

	var s config.ServerContext

	err = json.Unmarshal(content, &s)
	if err != nil {
		return nil, fmt.Errorf("error decode config file %s: %w", filename, err)
	}

	return &s, nil
}
