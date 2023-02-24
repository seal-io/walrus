package strs

import "strings"

func Join(sep string, strs ...string) string {
	return strings.Join(strs, sep)
}
