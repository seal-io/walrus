package strs

import "strings"

func Join(sep string, strs ...string) string {
	return strings.Join(strs, sep)
}

func Indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(v, "\n", "\n"+pad)
}
