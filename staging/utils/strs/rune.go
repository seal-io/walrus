package strs

import "github.com/mattn/go-runewidth"

func init() {
	runewidth.DefaultCondition.EastAsianWidth = true
	runewidth.DefaultCondition.CreateLUT()
}

// RuneWidth returns the number of cells in r.
// See http://www.unicode.org/reports/tr11/
func RuneWidth(r rune) int {
	return runewidth.RuneWidth(r)
}

// IsAmbiguousWidth returns whether is ambiguous width or not.
func IsAmbiguousWidth(r rune) bool {
	return runewidth.IsAmbiguousWidth(r)
}

// IsNeutralWidth returns whether is neutral width or not.
func IsNeutralWidth(r rune) bool {
	return runewidth.IsNeutralWidth(r)
}

// StringWidth return width as you can see.
func StringWidth(s string) int {
	return runewidth.StringWidth(s)
}

// Truncate return string truncated with w cells.
func Truncate(s string, w int, tail string) string {
	return runewidth.Truncate(s, w, tail)
}

// TruncateLeft cuts w cells from the beginning of the `s`.
func TruncateLeft(s string, w int, prefix string) string {
	return runewidth.TruncateLeft(s, w, prefix)
}

// Wrap return string wrapped with w cells.
func Wrap(s string, w int) string {
	return runewidth.Wrap(s, w)
}

// FillLeft return string filled in left by spaces in w cells.
func FillLeft(s string, w int) string {
	return runewidth.FillLeft(s, w)
}

// FillRight return string filled in left by spaces in w cells.
func FillRight(s string, w int) string {
	return runewidth.FillRight(s, w)
}

// HasSuffix returns true when one of the suffixes matched.
func HasSuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
			return true
		}
	}

	return false
}
