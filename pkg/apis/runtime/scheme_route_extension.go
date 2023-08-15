package runtime

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"

	"github.com/seal-io/seal/pkg/apis/runtime/openapi"
	"github.com/seal-io/seal/utils/strs"
)

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
		"/info",
		"/login",
		"/logout",
		"/tokens",
		"/tokens/:token",
		"/projects/:project/environments/:environment/graph",
		"/projects/:project/environments/:environment/services/:service/revisions/:servicerevision/terraform-states",
		"/projects/:project/environments/:environment/services/_/graph",
		"/projects/:project/environments/:environment/services/:service/graph",
		"/projects/:project/environments/:environment/resources/:serviceresource/keys",
		"/projects/:project/environments/:environment/resources/_/graph",
		"/projects/:project/connectors/:connector/repositories",
		"/projects/:project/connectors/:connector/repository-branches",
		"/connectors/:connector/repositories",
		"/connectors/:connector/repository-branches",
		"/perspectives/_/field-values",
		"/perspectives/_/fields",
	}
	// CliJsonYamlOutputFormatPaths is a list of paths that should be output as JSON/YAML.
	cliJsonYamlOutputFormatPaths = []string{
		"/projects/:project/environments/:environment/services/:service/revisions/:servicerevision/diff-latest",
		"/projects/:project/environments/:environment/services/:service/revisions/:servicerevision/diff-previous",
	}
)

func extendOperationSchema(r *Route, op *openapi3.Operation) {
	if op.Extensions == nil {
		op.Extensions = make(map[string]any)
	}

	switch {
	case len(r.Kinds) != 0 && slices.Contains(cliIgnoreKinds, r.Kinds[len(r.Kinds)-1]):
		op.Extensions[openapi.ExtCliIgnore] = true
	case slices.Contains(cliIgnorePaths, r.Path):
		op.Extensions[openapi.ExtCliIgnore] = true
	case r.Collection:
		if r.Method != http.MethodGet {
			op.Extensions[openapi.ExtCliIgnore] = true
		} else {
			op.Extensions[openapi.ExtCliOperationName] = "list"
		}
	case slices.Contains(
		[]string{
			resourceRouteNameCreate,
			resourceRouteNameGet,
			resourceRouteNameUpdate,
			resourceRouteNameDelete,
		}, r.GoFunc):
		op.Extensions[openapi.ExtCliOperationName] = strs.Dasherize(r.GoFunc)
	case r.Custom:
		op.Extensions[openapi.ExtCliOperationName] = strs.Dasherize(r.CustomName)
	}

	if slices.Contains(cliJsonYamlOutputFormatPaths, r.Path) {
		op.Extensions[openapi.ExtCliOutputFormat] = "json,yaml"
	}
}
