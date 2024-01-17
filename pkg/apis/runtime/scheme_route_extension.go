package runtime

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/pkg/apis/runtime/openapi"
	"github.com/seal-io/walrus/utils/strs"
)

var (
	// CliIgnoreKinds is a list of kinds to ignore when generating CLI commands.
	cliIgnoreKinds = []string{
		"Subject",
		"ProjectSubject",
		"Token",
		"SubjectRole",
		"Dashboard",
		"TemplateCompletion",
		"Cost",
		"WorkflowStepExecution",
		"Perspective",
	}
	// CliIgnorePaths is a list of paths to ignore when generating CLI commands.
	cliIgnorePaths = []string{
		"/info",
		"/login",
		"/logout",
		"/tokens",
		"/tokens/:token",
		"/projects/:project/environments/:environment/graph",
		"/projects/:project/environments/:environment/apply",
		"/projects/:project/environments/:environment/export",
		"/projects/:project/environments/:environment/resources/:resource/revisions/:resourcerevision/log",
		"/projects/:project/environments/:environment/resources/:resource/revisions/:resourcerevision/terraform-states",
		"/projects/:project/environments/:environment/resources/_/graph",
		"/projects/:project/environments/:environment/resources/:resource/graph",
		"/projects/:project/environments/:environment/resources/:resourcecomponent/keys",
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
		"/projects/:project/environments/:environment/resources/:resource/revisions/:resourcerevision/diff-latest",
		"/projects/:project/environments/:environment/resources/:resource/revisions/:resourcerevision/diff-previous",
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
		switch {
		case r.Method == http.MethodGet:
			op.Extensions[openapi.ExtCliOperationName] = "list"
		case r.Method == http.MethodPost:
			op.Extensions[openapi.ExtCliCmdIgnore] = true
			op.Extensions[openapi.ExtCliOperationName] = strs.Dasherize(r.GoFunc)
		case r.Method == http.MethodDelete:
			op.Extensions[openapi.ExtCliCmdIgnore] = true
			op.Extensions[openapi.ExtCliOperationName] = strs.Dasherize(r.GoFunc)
		default:
			op.Extensions[openapi.ExtCliIgnore] = true
		}
	case r.GoFunc == resourceRouteNamePatch:
		op.Extensions[openapi.ExtCliCmdIgnore] = true
		op.Extensions[openapi.ExtCliOperationName] = strs.Dasherize(r.GoFunc)
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
