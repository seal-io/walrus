package config

import (
	"fmt"
	"os"
	"path"

	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
)

const (
	configFileName = "config.json"
	cliName        = "walrus"
)

// InitConfig init config dir and load server context from cache.
func InitConfig() (*ServerContext, error) {
	// Log Level.
	log.SetLevel(log.InfoLevel)

	// Config dir.
	configDir := GetConfigDir()

	err := os.MkdirAll(configDir, 0o700)
	if err != nil {
		return nil, err
	}

	// Config file.
	filename := path.Join(configDir, configFileName)

	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &ServerContext{}, os.WriteFile(filename, []byte("{}"), 0o600)
		}

		return nil, fmt.Errorf("error stat config file %s: %w", filename, err)
	}

	return GetServerContextFromCache()
}

// GetConfigDir get config dir.
func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home dir: %w", err))
	}

	return path.Join(home, "."+cliName)
}

// GetServerContextFromCache load context from cache.
func GetServerContextFromCache() (*ServerContext, error) {
	filename := path.Join(GetConfigDir(), configFileName)
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error read config file %s: %w", filename, err)
	}

	var s ServerContext

	err = json.Unmarshal(content, &s)
	if err != nil {
		return nil, fmt.Errorf("error decode config file %s: %w", filename, err)
	}

	return &s, nil
}

// SetServerContextToCache set context to cache.
func SetServerContextToCache(s ServerContext) error {
	filename := path.Join(GetConfigDir(), configFileName)
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
