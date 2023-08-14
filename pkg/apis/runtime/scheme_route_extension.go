package runtime

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"

	cliapi "github.com/seal-io/seal/pkg/cli/api"
)

//nolint:lll
var (
	// CliIgnoreKinds is a list of kinds to ignore when generating CLI commands.
	cliIgnoreKinds = []string{
		"Subject",
		"Token",
		"SubjectRole",
		"Dashboard",
		"TemplateCompletion",
		"Cost",
	}
	// CliIgnorePaths is a list of paths to ignore when generating CLI commands.
	cliIgnorePaths = []string{
		"/projects/:project/environments/:environment/services/:service/service-revisions/:servicerevision/terraform-states",
		"/projects/:project/environments/:environment/services/_/graph",
		"/projects/:project/environments/:environment/service-resources/:serviceresource/keys",
		"/projects/:project/environments/:environment/service-resources/_/graph",
		"/projects/:project/connectors/:connector/repositories",
		"/projects/:project/connectors/:connector/repository-branches",
		"/connectors/:connector/repositories",
		"/connectors/:connector/repository-branches",
		"/perspectives/_/field-values",
		"/perspectives/_/fields",
	}
	// CliJsonYamlOutputFormatPaths is a list of paths that should be output as JSON/YAML.
	cliJsonYamlOutputFormatPaths = []string{
		"/projects/:project/environments/:environment/services/:service/service-revisions/:servicerevision/diff-latest",
		"/projects/:project/environments/:environment/services/:service/service-revisions/:servicerevision/diff-previous",
	}
)

func extendOperationSchema(r *Route, op *openapi3.Operation) {
	if op.Extensions == nil {
		op.Extensions = make(map[string]any)
	}

	switch {
	case len(r.Kinds) != 0 && slices.Contains(cliIgnoreKinds, r.Kinds[len(r.Kinds)-1]):
		op.Extensions[cliapi.ExtCliIgnore] = true
	case slices.Contains(cliIgnorePaths, r.Path):
		op.Extensions[cliapi.ExtCliIgnore] = true
	case r.Collection:
		if r.Method != http.MethodGet {
			op.Extensions[cliapi.ExtCliIgnore] = true
		} else {
			op.Extensions[cliapi.ExtCliOperationName] = "list"
		}
	}

	if slices.Contains(cliJsonYamlOutputFormatPaths, r.Path) {
		op.Extensions[cliapi.ExtCliOutputFormat] = "json,yaml"
	}
}
