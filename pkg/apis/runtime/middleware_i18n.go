package runtime

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	_ "github.com/seal-io/seal/pkg/i18n"
	"github.com/seal-io/seal/utils/log"
)

const (
	languageContextKey = "language"
)

var (
	supported = []language.Tag{
		language.English,
		language.Chinese,
	}

	matcher = language.NewMatcher(supported)
)

func I18n() Handle {
	return func(c *gin.Context) {
		setLanguageTag(c)
		c.Next()
	}
}

func setLanguageTag(c *gin.Context) {
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage == "" {
		return
	}

	tags, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil {
		log.Warnf("failed to parse the Accept-Language header %q: %w", acceptLanguage, err)
		return
	}
	tag, _, _ := matcher.Match(tags...)
	c.Set(languageContextKey, tag)
}

func getLanguageTag(c *gin.Context) language.Tag {
	defaultLanguage := language.English
	v, exist := c.Get(languageContextKey)

	if !exist {
		return defaultLanguage
	}

	if tag, ok := v.(language.Tag); ok {
		return tag
	}

	return defaultLanguage
}

func Translate(c *gin.Context, s string) string {
	tag := getLanguageTag(c)
	p := message.NewPrinter(tag)

	return p.Sprintf(s)
}
