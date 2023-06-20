package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"

	"github.com/seal-io/seal/pkg/cli/config"
	"github.com/seal-io/seal/pkg/cli/formatter"
	"github.com/seal-io/seal/utils/json"
)

// Operation represents an API action, e.g. list-things or create-user.
type Operation struct {
	Name          string      `json:"name"`
	Group         string      `json:"group,omitempty"`
	Short         string      `json:"short,omitempty"`
	Long          string      `json:"long,omitempty"`
	Method        string      `json:"method,omitempty"`
	URITemplate   string      `json:"uriTemplate"`
	PathParams    []*Param    `json:"pathParams,omitempty"`
	QueryParams   []*Param    `json:"queryParams,omitempty"`
	HeaderParams  []*Param    `json:"headerParams,omitempty"`
	BodyParams    *BodyParams `json:"bodyParams,omitempty"`
	BodyMediaType string      `json:"bodyMediaType,omitempty"`
	Hidden        bool        `json:"hidden,omitempty"`
	Deprecated    string      `json:"deprecated,omitempty"`
}

// Command returns a Cobra command instance for this operation.
func (o Operation) Command(sc *config.Config) *cobra.Command {
	var (
		body  interface{}
		flags = map[string]interface{}{}
	)

	use := o.Name
	for _, p := range o.PathParams {
		use += " " + fmt.Sprintf("<%s>", p.Name)
	}

	argSpec := cobra.ExactArgs(len(o.PathParams))
	if o.BodyMediaType != "" {
		argSpec = cobra.MinimumNArgs(len(o.PathParams))
	}

	sub := &cobra.Command{
		Use:        use,
		Short:      o.Short,
		Long:       o.Long,
		Args:       argSpec,
		Hidden:     o.Hidden,
		Deprecated: o.Deprecated,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := sc.Inject(cmd)
			if err != nil {
				fmt.Fprintln(os.Stderr)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			req, err := o.Request(cmd, args, flags, body)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			resp, err := sc.DoRequest(req)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			b, err := formatter.Format(sc.Format, resp)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			os.Stdout.Write(b)
		},
	}

	for _, p := range o.QueryParams {
		flags[p.Name] = p.AddFlag(sub.Flags())
	}

	for _, p := range o.HeaderParams {
		flags[p.Name] = p.AddFlag(sub.Flags())
	}

	if o.BodyParams != nil {
		switch o.BodyParams.Type {
		case openapi3.TypeArray:
			// Array request body is considered as a single params.
			if len(o.BodyParams.Params) != 0 {
				b := o.BodyParams.Params[0]
				bp := b.AddFlag(sub.Flags())
				body = bp
			}
		case openapi3.TypeObject:
			bps := make(map[string]interface{})
			for _, p := range o.BodyParams.Params {
				bps[p.Name] = p.AddFlag(sub.Flags())
			}
			body = bps
		}
	}

	for _, v := range sc.InjectFields() {
		_ = sub.Flags().MarkHidden(v)
	}

	return sub
}

// Request generate http request base on the operation.
func (o Operation) Request(
	cmd *cobra.Command,
	args []string,
	flags map[string]interface{},
	body interface{},
) (*http.Request, error) {
	// Replaces URL-encoded `{`+name+`}` in the uri.
	uri := o.URITemplate

	for i, param := range o.PathParams {
		uri = strings.Replace(uri, "{"+param.Name+"}", fmt.Sprintf("%v", args[i]), 1)
	}

	// Generate URL queries.
	query := url.Values{}

	for _, param := range o.QueryParams {
		flag := flags[param.Name]
		for _, v := range param.Serialize(flag) {
			if v != "" {
				query.Add(param.Name, v)
			}
		}
	}

	queryEncoded := query.Encode()
	if queryEncoded != "" {
		if strings.Contains(uri, "?") {
			uri += "&"
		} else {
			uri += "?"
		}
		uri += queryEncoded
	}

	// Generate Headers.
	headers := http.Header{}

	for _, param := range o.HeaderParams {
		// Ignore flags not passed from the user.
		if !cmd.Flags().Changed(param.OptionName()) {
			continue
		}

		for _, v := range param.Serialize(flags[param.Name]) {
			headers.Add(param.Name, v)
		}
	}

	// Generate request body.
	var br io.Reader

	if o.BodyMediaType != "" {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("invalid body: %w", err)
		}
		br = bytes.NewReader(b)
	}

	req, err := http.NewRequest(o.Method, uri, br)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	return req, nil
}
