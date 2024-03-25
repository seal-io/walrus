package config

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/resources/interpolation"
	"github.com/seal-io/walrus/pkg/resourcestate"
	pkgvariable "github.com/seal-io/walrus/pkg/variable"
	"github.com/seal-io/walrus/utils/json"
)

// _interpolationReg is the regular expression for matching non-reference or non-variable expressions.
// Reference: https://developer.hashicorp.com/terraform/language/expressions/strings#escape-sequences-1
// To handle escape sequences, ${xxx} is converted to $${xxx}.
// If there are more than two consecutive $ symbols, like $${xxx}, they are further converted to $$${xxx}.
// During Terraform processing, $${} is ultimately transformed back to ${};
// this interpolation is used to ensure a WYSIWYG user experience.
var _interpolationReg = regexp.MustCompile(`\$\{((var\.|res\.)?([^.}]+)(?:\.([^.}]+))?)[^\}]*\}`)

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
) (attrs map[string]any, variables model.Variables, outputs map[string]types.OutputValue, err error) {
	var (
		templateVariables         []string
		dependencyResourceOutputs []string
	)

	replaced := !onlyValidated

	attrs, templateVariables, dependencyResourceOutputs, err = parseAttributeReplace(attributes, replaced)
	if err != nil {
		return nil, nil, nil, err
	}

	// If a resource run has variables that inherit from cloned run, use them directly.
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

	vars := make(map[string]*model.Variable)
	for _, v := range variables {
		vars[v.Name] = v
	}

	if !onlyValidated {
		// Replace variables and resource references.
		bs, err := json.Marshal(attributes)
		if err != nil {
			return nil, nil, nil, err
		}

		bs = interpolation.VariableReg.ReplaceAllFunc(bs, func(match []byte) []byte {
			m := interpolation.VariableReg.FindSubmatch(match)

			if len(m) != 2 {
				return match
			}

			if v, ok := vars[string(m[1])]; ok {
				return []byte(v.Value)
			}

			return match
		})

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

		bs = interpolation.ResourceReg.ReplaceAllFunc(bs, func(match []byte) []byte {
			m := interpolation.ResourceReg.FindSubmatch(match)

			if len(m) != 3 {
				return match
			}

			if v, ok := outputs[fmt.Sprintf("%s_%s", m[1], m[2])]; ok {
				var str string
				err := json.Unmarshal(v.Value, &str)
				if err != nil {
					return v.Value
				}
				return []byte(str)
			}

			return match
		})

		attrs = make(map[string]any)
		err = json.Unmarshal(bs, &attrs)
		if err != nil {
			return nil, nil, nil, err
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

	variableMatches := interpolation.VariableReg.FindAllSubmatch(bs, -1)
	resourceMatches := interpolation.ResourceReg.FindAllSubmatch(bs, -1)

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

	// Replace interpolation from ${} to $${} to avoid escape sequences.
	bs = _interpolationReg.ReplaceAllFunc(bs, func(match []byte) []byte {
		m := _interpolationReg.FindSubmatch(match)

		if len(m) != 5 {
			return match
		}

		// If it is a variable or resource reference, do not replace.
		if string(m[2]) == "var." || string(m[2]) == "res." {
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
) (map[string]types.OutputValue, error) {
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

	return resourcestate.GetDependencyOutputs(ctx, client, dependencyResourceIDs, dependOutputs)
}
