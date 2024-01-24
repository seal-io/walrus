package formatter

import (
	"net/http"
)

// Format generate formatted result for response.
func Format(resp *http.Response, opts Options) ([]byte, error) {
	var f Formatter

	switch opts.Format {
	case "json":
		f = &JsonFormatter{}
	case "yaml":
		f = &YamlFormatter{}
	default:
		f = &TableFormatter{
			Columns:   opts.Columns,
			Group:     opts.Group,
			Operation: opts.Operation,
		}
	}

	return f.Format(resp)
}

type Options struct {
	Format    string
	Columns   []string
	Group     string
	Operation string
}

// Formatter generate formatted result for response.
type Formatter interface {
	Format(*http.Response) ([]byte, error)
}
