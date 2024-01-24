package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/formatter"
	"github.com/seal-io/walrus/utils/json"
)

// Operation represents an API action, e.g. list-things or create-user.
type Operation struct {
	Name          string      `json:"name"`
	Group         string      `json:"group,omitempty"`
	Short         string      `json:"short,omitempty"`
	Long          string      `json:"long,omitempty"`
	Method        string      `json:"method,omitempty"`
	URITemplate   string      `json:"uriTemplate"`
	URIParams     []string    `json:"uriParams"`
	PathParams    []*Param    `json:"pathParams,omitempty"`
	QueryParams   []*Param    `json:"queryParams,omitempty"`
	HeaderParams  []*Param    `json:"headerParams,omitempty"`
	BodyParams    *BodyParams `json:"bodyParams,omitempty"`
	BodyMediaType string      `json:"bodyMediaType,omitempty"`
	Hidden        bool        `json:"hidden,omitempty"`
	Deprecated    string      `json:"deprecated,omitempty"`
	Formats       []string    `json:"formats,omitempty"`
	TableColumns  []string    `json:"tableColumns,omitempty"`

	// CmdIgnore is used to ignore the operation when generating CLI commands.
	CmdIgnore bool `json:"cmdIgnore,omitempty"`
}

// Command returns a Cobra command instance for this operation.
func (o Operation) Command(sc *config.Config) *cobra.Command {
	var (
		body  any
		flags = map[string]any{}
	)

	var (
		use      = o.Name
		argCount = 0
	)

	for _, p := range o.PathParams {
		switch {
		case p.DataFrom == DataFromContextAndArg && sc.ContextExisted(p.Name):
			continue
		case p.DataFrom == DataFromFlag:
			continue
		default:
		}

		use += " " + fmt.Sprintf("<%s>", p.Name)
		argCount += 1
	}

	argSpec := cobra.ExactArgs(argCount)
	if o.BodyMediaType != "" {
		argSpec = cobra.MinimumNArgs(argCount)
	}

	sub := &cobra.Command{
		Use:        use,
		Short:      o.Short,
		Long:       o.Long,
		Args:       argSpec,
		Hidden:     o.Hidden,
		Deprecated: o.Deprecated,
		PreRun: func(cmd *cobra.Command, args []string) {
			err := sc.Inject(cmd)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			req, err := o.Request(cmd, args, flags, body, sc.ServerContext)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			resp, err := sc.DoRequest(req)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			format, ok := flags["output"].(*string)
			if !ok {
				fmt.Fprintln(os.Stderr, fmt.Errorf("invalid output format"))
				os.Exit(1)
			}

			b, err := formatter.Format(resp, formatter.Options{
				Format:    o.format(*format),
				Columns:   o.TableColumns,
				Group:     o.Group,
				Operation: o.Name,
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			os.Stdout.Write(b)
			os.Stdout.Write([]byte{'\n'})
		},
	}

	for _, p := range o.QueryParams {
		flags[p.Name] = p.AddFlag(sub.Flags())
	}

	for _, p := range o.HeaderParams {
		flags[p.Name] = p.AddFlag(sub.Flags())
	}

	for _, p := range o.PathParams {
		if p.DataFrom == DataFromFlag {
			flags[p.Name] = p.AddFlag(sub.Flags())
		}
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
			bps := make(map[string]any)
			for _, p := range o.BodyParams.Params {
				bps[p.Name] = p.AddFlag(sub.Flags())
			}
			body = bps
		}
	}

	flags["output"] = sub.Flags().StringP("output", "o", "table", "Output format [table, json, yaml]")

	return sub
}

// Request generate http request base on the operation.
func (o Operation) Request(
	cmd *cobra.Command,
	args []string,
	flags map[string]any,
	body any,
	sc config.ServerContext,
) (*http.Request, error) {
	// Generate URI template.
	uriTemplate := o.URITemplate

	if len(o.URIParams) != 0 {
		data := make(map[string]string)

		for _, k := range o.URIParams {
			if val := flags[k]; val != nil {
				data[k] = *val.(*string)
			}
		}

		var buf bytes.Buffer
		tmpl := template.Must(
			template.New("uri").Parse(uriTemplate),
		)

		err := tmpl.Execute(&buf, data)
		if err != nil {
			return nil, err
		}

		uriTemplate = buf.String()
	}

	// Replaces URL-encoded `{`+name+`}` in the uri.
	var argCount int

	for _, param := range o.PathParams {
		paramPlaceholder := "{" + param.Name + "}"

		switch {
		case param.DataFrom == DataFromContextAndArg:
			// Inject from context.
			uriTemplate = sc.InjectURI(uriTemplate, param.Name)

			// Inject from arg.
			if argCount < len(args) && strings.Contains(uriTemplate, paramPlaceholder) {
				uriTemplate = strings.Replace(uriTemplate, paramPlaceholder, fmt.Sprintf("%v", args[argCount]), 1)
				argCount += 1
			}
		case param.DataFrom == DataFromFlag:
			flag := flags[param.Name]

			se := param.Serialize(flag)
			if len(se) == 0 {
				continue
			}
			uriTemplate = strings.Replace(uriTemplate, paramPlaceholder, fmt.Sprintf("%v", se[0]), 1)
		case param.DataFrom == DataFromArg:
			uriTemplate = strings.Replace(uriTemplate, paramPlaceholder, fmt.Sprintf("%v", args[argCount]), 1)
			argCount += 1
		}
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
		if strings.Contains(uriTemplate, "?") {
			uriTemplate += "&"
		} else {
			uriTemplate += "?"
		}
		uriTemplate += queryEncoded
	}

	// Generate Headers.
	headers := http.Header{}

	for _, param := range o.HeaderParams {
		// Ignore flags not passed from the user.
		if cmd != nil && !cmd.Flags().Changed(param.OptionName()) {
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

	req, err := http.NewRequest(o.Method, uriTemplate, br)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	return req, nil
}

func (o Operation) format(flagFormat string) string {
	if len(o.Formats) != 0 {
		if slices.Contains(o.Formats, flagFormat) {
			return flagFormat
		}

		return o.Formats[0]
	}

	return flagFormat
}
