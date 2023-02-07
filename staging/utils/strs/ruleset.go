package strs

import (
	"strings"

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
}

// Pluralize returns the plural form of a word.
func Pluralize(word string) string {
	return globalRuleset.Pluralize(word)
}

// Singularize returns the singular form of a word.
func Singularize(word string) string {
	return globalRuleset.Singularize(word)
}

// Camelize returns the string in camel case, "dino_party" -> "DinoParty".
func Camelize(word string) string {
	return globalRuleset.Camelize(word)
}

// CamelizeDownFirst is the same as Camelize but with first lowercase char.
func CamelizeDownFirst(word string) string {
	return globalRuleset.CamelizeDownFirst(word)
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
