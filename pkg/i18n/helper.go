package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	supported = []language.Tag{
		language.English,
		language.Chinese,
	}
	matcher = language.NewMatcher(supported)
)

// T is a helper function to print the i18n message according to the given key and tags.
func T(key string, tags ...language.Tag) string {
	tag, _, _ := matcher.Match(tags...)
	return message.NewPrinter(tag).Sprintf(key)
}
