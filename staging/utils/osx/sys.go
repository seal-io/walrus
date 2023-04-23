package osx

import "os"

func Hostname(def ...string) string {
	var h, _ = os.Hostname()
	if h == "" && len(def) != 0 {
		return def[0]
	}
	return h
}
