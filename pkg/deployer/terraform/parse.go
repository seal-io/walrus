package terraform

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	pkgvariable "github.com/seal-io/walrus/pkg/variable"
	"github.com/seal-io/walrus/utils/json"
)

type RunOpts struct {
	ResourceRun *model.ResourceRun

	ResourceName string

	ProjectID     object.ID
	EnvironmentID object.ID
}

// ParseModuleAttributes parse module variables and dependencies.
func ParseModuleAttributes(
	ctx context.Context,
	mc model.ClientSet,
	attributes map[string]any,
	onlyValidated bool,
	opts RunOpts,
) (attrs map[string]any, variables model.Variables, outputs map[string]parser.OutputState, err error) {
	var (
		templateVariables         []string
		dependencyResourceOutputs []string
	)

	replaced := !onlyValidated

	attrs, templateVariables, dependencyResourceOutputs, err = parseAttributeReplace(attributes, replaced)
	if err != nil {
		return
	}

	// If resource run has variables that inherit from cloned run, use them directly.
	if opts.ResourceRun != nil && len(opts.ResourceRun.Variables) > 0 {
		for k, v := range opts.ResourceRun.Variables {
			variables = append(variables, &model.Variable{
				Name:  k,
				Value: crypto.String(v),
			})
		}
	} else {
		variables, err = pkgvariable.Get(ctx, mc, templateVariables, opts.ProjectID, opts.EnvironmentID)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if !onlyValidated {
		dependOutputMap := toDependOutputMap(dependencyResourceOutputs)

		outputs, err = getResourceDependencyOutputsByID(
			ctx,
			mc,
			opts.ResourceRun.ResourceID,
			dependOutputMap)
		if err != nil {
			return nil, nil, nil, err
		}

		// Check if all dependency resource outputs are found.
		for outputName := range dependOutputMap {
			if _, ok := outputs[outputName]; !ok {
				return nil, nil, nil, fmt.Errorf("resource %s dependency output %s not found",
					opts.ResourceName, outputName)
			}
		}
	}

	return attrs, variables, outputs, nil
}

// toDependOutputMap splits the dependencyResourceOutputs from {resource}_{resource_name}_{output_name}
// to a map of {resource_name}_{output_name}:{resource}.
func toDependOutputMap(dependencyResourceOutputs []string) map[string]string {
	dependOutputMap := make(map[string]string, 0)

	for _, dependOutput := range dependencyResourceOutputs {
		ss := strings.SplitN(dependOutput, "_", 2)
		if len(ss) != 2 {
			continue
		}
		dependOutputMap[ss[1]] = ss[0]
	}

	return dependOutputMap
}

// parseAttributeReplace parses attribute variable ${var.name} replaces it with ${var._variablePrefix+name},
// resource reference ${res.name.output} replaces it with ${var._resourcePrefix+name}
// and returns variable names and output names.
// Replaced flag indicates whether to replace the module attribute variable with prefix string.
func parseAttributeReplace(
	attributes map[string]any,
	replaced bool,
) (map[string]any, []string, []string, error) {
	bs, err := json.Marshal(attributes)
	if err != nil {
		return nil, nil, nil, err
	}

	variableMatches := pkgresource.VariableReg.FindAllSubmatch(bs, -1)
	resourceMatches := pkgresource.ResourceReg.FindAllSubmatch(bs, -1)

	variableMatched := sets.NewString()

	for _, match := range variableMatches {
		if len(match) > 1 {
			variableMatched.Insert(string(match[1]))
		}
	}

	resourceMatched := sets.NewString()

	for _, match := range resourceMatches {
		if len(match) > 1 {
			// Resource outputs are in the form:
			// - res_{resource_name}_{output_name}.
			resourceMatched.Insert("res_" + string(match[1]) + "_" + string(match[2]))
		}
	}

	variableRepl := "${var." + _variablePrefix + "${1}}"
	bs = pkgresource.VariableReg.ReplaceAll(bs, []byte(variableRepl))

	resourceRepl := "${var." + _resourcePrefix + "${1}_${2}}"
	bs = pkgresource.ResourceReg.ReplaceAll(bs, []byte(resourceRepl))

	// Replace interpolation from ${} to $${} to avoid escape sequences.
	bs = _interpolationReg.ReplaceAllFunc(bs, func(match []byte) []byte {
		m := _interpolationReg.FindSubmatch(match)

		if len(m) != 5 {
			return match
		}

		// If it is variable or resource reference, do not replace.
		if string(m[2]) == "var." {
			return match
		}

		var b bytes.Buffer

		b.WriteString("$")
		b.Write(match)

		return b.Bytes()
	})

	if replaced {
		err = json.Unmarshal(bs, &attributes)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return attributes, variableMatched.List(), resourceMatched.List(), nil
}

// getResourceDependencyOutputsByID gets the dependency outputs of the resource by resource id.
func getResourceDependencyOutputsByID(
	ctx context.Context,
	client model.ClientSet,
	resourceID object.ID,
	dependOutputs map[string]string,
) (map[string]parser.OutputState, error) {
	entity, err := client.Resources().Query().
		Where(resource.ID(resourceID)).
		WithDependencies(func(sq *model.ResourceRelationshipQuery) {
			sq.Where(func(s *sql.Selector) {
				s.Where(sql.ColumnsNEQ(resourcerelationship.FieldResourceID, resourcerelationship.FieldDependencyID))
			})
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dependencyResourceIDs := make([]object.ID, 0, len(entity.Edges.Dependencies))

	for _, d := range entity.Edges.Dependencies {
		if d.Type != types.ResourceRelationshipTypeImplicit {
			continue
		}

		dependencyResourceIDs = append(dependencyResourceIDs, d.DependencyID)
	}

	return pkgresource.GetDependencyOutputs(ctx, client, dependencyResourceIDs, dependOutputs)
}
