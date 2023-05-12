package osx

import "os"

func Getenv(key string, def ...string) string {
	e := os.Getenv(key)
	if e == "" && len(def) != 0 {
		return def[0]
	}

	return e
}
