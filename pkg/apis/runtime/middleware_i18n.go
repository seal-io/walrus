package runtime

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/seal-io/walrus/utils/log"

	_ "github.com/seal-io/walrus/pkg/i18n"
)

const (
	languageContextKey = "request_language"
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
		log.WithName("api").
			Warnf("failed to parse the Accept-Language header %q: %w", acceptLanguage, err)
		return
	}
	tag, _, _ := matcher.Match(tags...)
	c.Set(languageContextKey, tag)
}

func getLanguageTag(ctx context.Context) language.Tag {
	def := language.English

	c, ok := ctx.(*gin.Context)
	if !ok {
		c, ok = ctx.Value(gin.ContextKey).(*gin.Context)
		if !ok {
			return def
		}
	}

	v, exist := c.Get(languageContextKey)
	if !exist {
		return def
	}

	if tag, ok := v.(language.Tag); ok {
		return tag
	}

	return def
}

func Translate(ctx context.Context, s string) string {
	tag := getLanguageTag(ctx)
	p := message.NewPrinter(tag)

	return p.Sprintf(s)
}
