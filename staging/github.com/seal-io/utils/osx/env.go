package osx

import (
	"os"
)

// LookupEnv retrieves the value of the environment variable named
// by the key. If the variable is present in the environment the
// value (which may be empty) is returned and the boolean is true.
// Otherwise, the returned value will be empty and the boolean will
// be false.
func LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Getenv retrieves the value of the environment variable named by the key.
// It returns the default, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use LookupEnv.
func Getenv(key string, def ...string) string {
	e, ok := LookupEnv(key)
	if !ok && len(def) != 0 {
		return def[0]
	}

	return e
}

// ExpandEnv is similar to Getenv,
// but replaces ${var} or $var in the result.
func ExpandEnv(key string, def ...string) string {
	return os.ExpandEnv(Getenv(key, def...))
}
