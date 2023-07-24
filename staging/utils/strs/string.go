package strs

import (
	"bytes"
	"strings"
)

func Join[T ~string](sep string, strs ...T) string {
	switch len(strs) {
	case 0:
		return ""
	case 1:
		return string(strs[0])
	}

	n := len(sep) * (len(strs) - 1)
	for i := 0; i < len(strs); i++ {
		n += len(strs[i])
	}

	var b strings.Builder

	b.Grow(n)
	b.WriteString(string(strs[0]))

	for i := range strs[1:] {
		b.WriteString(sep)
		b.WriteString(string(strs[i+1]))
	}

	return b.String()
}

func Indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(v, "\n", "\n"+pad)
}

// NormalizeSpecialChars replaces special characters with their normalized equivalents.
// This is useful to avoid encoding issues in PostgresSQL(Invalid byte sequence for encoding "UTF8").
// Each byte in the input slice is examined:
//   - If the byte is not a special character, it is simply copied to the output slice.
//   - If the byte is a non-breaking space (0xA0), it is replaced with a space (0x20),
//     but only if the last character is non-ASCII character (0xC2)
//   - If the byte is a null character (0x00), it is replaced with the string "NUL" (0x4E 0x55 0x4C).
//
// The resulting byte slice is converted back to a string and returned.
func NormalizeSpecialChars(s string) string {
	var (
		buf      bytes.Buffer
		lastChar byte

		sbs = []byte(s)
	)

	for i := 0; i < len(sbs); i++ {
		switch sbs[i] {
		default:
			buf.WriteByte(sbs[i])
			lastChar = sbs[i]
		case 0xA0:
			if lastChar == 0xC2 {
				buf.WriteByte(0x20)
			} else {
				buf.WriteByte(sbs[i])
			}

			lastChar = 0x20
		case 0x00:
			buf.Write([]byte{0x4e, 0x55, 0x4c})

			lastChar = 0x4c
		}
	}

	return buf.String()
}

// LastContent retrieves the last characters of a string.
func LastContent(content string, length int) string {
	if len(content) < length {
		return content
	}

	return content[len(content)-length:]
}

// FirstContent retrieves the leading characters of a string.
func FirstContent(content string, length int) string {
	if len(content) < length {
		return content
	}

	return content[:length]
}
