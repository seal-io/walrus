package strs

import (
	"strings"
	"unicode"

	"github.com/akerl/go-indefinite-article/indefinite"
	"github.com/go-openapi/inflect"
)

var globalRuleset = inflect.NewDefaultRuleset()

func init() {
	for _, w := range []string{
		"ACL", "API", "ACME", "ASCII", "AWS", "CPU", "CSS", "DNS", "EOF", "GB", "GUID",
		"HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "KB", "LHS", "MAC", "MB",
		"QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "SSO", "TCP",
		"TLS", "TTL", "UDP", "UI", "UID", "URI", "URL", "UTF8", "UUID", "VM",
		"XML", "XMPP", "XSRF", "XSS",
	} {
		globalRuleset.AddAcronym(w)
	}

	for s, p := range map[string]string{
		"exec": "exec",
	} {
		globalRuleset.AddIrregular(s, p)
	}
}

// Pluralize returns the plural form of a word.
func Pluralize(word string) string {
	return globalRuleset.Pluralize(word)
}

// Singularize returns the singular form of a word.
func Singularize(word string) string {
	return globalRuleset.Singularize(word)
}

// SingularizeWithArticle returns the singular form of a word and prepends the article to it, "dogs" -> "a dog".
func SingularizeWithArticle(word string) string {
	return indefinite.AddArticle(globalRuleset.Singularize(word))
}

// Camelize returns the string in camel case, "dino_party" -> "DinoParty".
func Camelize(word string) string {
	return globalRuleset.Camelize(word)
}

// CamelizeDownFirst is the same as Camelize but with first lowercase char.
func CamelizeDownFirst(word string) string {
	return globalRuleset.CamelizeDownFirst(word)
}

// Decamelize converts camel-case words to space-split words, "DinoParty" -> "Dino Party".
// If lowercase is true, convert the first letter of each word to lowercase, "DinoParty" -> "dino party".
func Decamelize(s string, lowercase bool) string {
	var (
		b          strings.Builder
		splittable = false
	)

	for _, v := range s {
		if splittable && unicode.IsUpper(v) {
			b.WriteByte(' ')
		}

		if lowercase {
			b.WriteString(strings.ToLower(string(v)))
		} else {
			b.WriteRune(v)
		}

		splittable = unicode.IsLower(v) || unicode.IsNumber(v)
	}

	return b.String()
}

// Underscore returns the string in snake case, "BigBen" -> "big_ben".
func Underscore(word string) string {
	return globalRuleset.Underscore(word)
}

// UnderscoreUpper is the same as Underscore but with all uppercase char.
func UnderscoreUpper(word string) string {
	return strings.ToUpper(Underscore(word))
}

// Dasherize returns the string in dasherized case, "SomeText" -> "some-text".
func Dasherize(word string) string {
	return globalRuleset.Dasherize(word)
}

// Capitalize returns the string with first uppercase letter.
func Capitalize(word string) string {
	return globalRuleset.Capitalize(word)
}

// Ordinalize returns the ordination case, "1031" -> "1031st".
func Ordinalize(word string) string {
	return globalRuleset.Ordinalize(word)
}
