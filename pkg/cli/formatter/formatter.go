package formatter

import (
	"net/http"
)

// Format generate formatted result for response.
func Format(format string, resp *http.Response) ([]byte, error) {
	var f Formatter

	switch format {
	case "json":
		f = &JsonFormatter{}
	case "yaml":
		f = &YamlFormatter{}
	default:
		f = &TableFormatter{}
	}

	return f.Format(resp)
}

// Formatter generate formatted result for response.
type Formatter interface {
	Format(*http.Response) ([]byte, error)
}
