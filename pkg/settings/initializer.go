package settings

import (
	"fmt"
	"os"

	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/strs"
)

// initializer defines the stereotype for reading initialized value.
type initializer func(string) string

// initializeFromSpecifiedEnv searches the variable of the given environment name,
// returns the default value if not found the specified environment.
//
//nolint:unparam
func initializeFromSpecifiedEnv(envName, defValue string) initializer {
	return func(id string) string {
		envValue := os.Getenv(envName)
		if envValue != "" {
			logger.Debugf("loaded %s initial value from %s environment variable", id, envName)
			return envValue
		}

		return defValue
	}
}

// initializeFromEnv searches the variable of the `SERVER_SETTING_${UpperSnakeCase_of_SettingName}`,
// returns the default value if not found the environment.
func initializeFromEnv(defValue string) initializer {
	return func(id string) string {
		envName := "SERVER_SETTING_" + strs.UnderscoreUpper(id)
		envValue := os.Getenv(envName)

		if envValue != "" {
			logger.Debugf("loaded %s initial value from %s environment variable", id, envName)
			return envValue
		}

		return defValue
	}
}

// initializeFrom initializes with the given val.
func initializeFrom(defValue string) initializer {
	return func(id string) string {
		return defValue
	}
}

// initializeFromJSON initializes with the given val as JSON.
func initializeFromJSON(defValue interface{}) initializer {
	return func(id string) string {
		v, err := json.Marshal(defValue)
		if err != nil {
			panic(fmt.Errorf("error marshaling initialization value: %w", err))
		}

		return string(v)
	}
}
