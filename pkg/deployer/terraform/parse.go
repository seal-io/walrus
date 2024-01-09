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
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model/variable"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/json"
)

type RevisionOpts struct {
	ResourceRevision *model.ResourceRevision

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
	opts RevisionOpts,
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

	// If resource revision has variables that inherit from cloned revision, use them directly.
	if opts.ResourceRevision != nil && len(opts.ResourceRevision.Variables) > 0 {
		for k, v := range opts.ResourceRevision.Variables {
			variables = append(variables, &model.Variable{
				Name:  k,
				Value: crypto.String(v),
			})
		}
	} else {
		variables, err = getVariables(ctx, mc, templateVariables, opts.ProjectID, opts.EnvironmentID)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if !onlyValidated {
		dependOutputMap := toDependOutputMap(dependencyResourceOutputs)

		outputs, err = getResourceDependencyOutputsByID(
			ctx,
			mc,
			opts.ResourceRevision.ResourceID,
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

	variableMatches := _variableReg.FindAllSubmatch(bs, -1)
	resourceMatches := _resourceReg.FindAllSubmatch(bs, -1)

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
	bs = _variableReg.ReplaceAll(bs, []byte(variableRepl))

	resourceRepl := "${var." + _resourcePrefix + "${1}_${2}}"
	bs = _resourceReg.ReplaceAll(bs, []byte(resourceRepl))

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

func getVariables(
	ctx context.Context,
	client model.ClientSet,
	variableNames []string,
	projectID, environmentID object.ID,
) (model.Variables, error) {
	var variables model.Variables

	if len(variableNames) == 0 {
		return variables, nil
	}

	nameIn := make([]any, len(variableNames))
	for i, name := range variableNames {
		nameIn[i] = name
	}

	type scanVariable struct {
		Name      string        `json:"name"`
		Value     crypto.String `json:"value"`
		Sensitive bool          `json:"sensitive"`
		Scope     int           `json:"scope"`
	}

	var vars []scanVariable

	err := client.Variables().Query().
		Modify(func(s *sql.Selector) {
			var (
				envPs = sql.And(
					sql.EQ(variable.FieldProjectID, projectID),
					sql.EQ(variable.FieldEnvironmentID, environmentID),
				)

				projPs = sql.And(
					sql.EQ(variable.FieldProjectID, projectID),
					sql.IsNull(variable.FieldEnvironmentID),
				)

				globalPs = sql.IsNull(variable.FieldProjectID)
			)

			s.Where(
				sql.And(
					sql.In(variable.FieldName, nameIn...),
					sql.Or(
						envPs,
						projPs,
						globalPs,
					),
				),
			).SelectExpr(
				sql.Expr("CASE "+
					"WHEN project_id IS NOT NULL AND environment_id IS NOT NULL THEN 3 "+
					"WHEN project_id IS NOT NULL AND environment_id IS NULL THEN 2 "+
					"ELSE 1 "+
					"END AS scope"),
			).AppendSelect(
				variable.FieldName,
				variable.FieldValue,
				variable.FieldSensitive,
			)
		}).
		Scan(ctx, &vars)
	if err != nil {
		return nil, err
	}

	found := make(map[string]scanVariable)
	for _, v := range vars {
		ev, ok := found[v.Name]
		if !ok {
			found[v.Name] = v
			continue
		}

		if v.Scope > ev.Scope {
			found[v.Name] = v
		}
	}

	// Validate module variable are all exist.
	foundSet := sets.NewString()
	for n, e := range found {
		foundSet.Insert(n)
		variables = append(variables, &model.Variable{
			Name:      n,
			Value:     e.Value,
			Sensitive: e.Sensitive,
		})
	}
	requiredSet := sets.NewString(variableNames...)

	missingSet := requiredSet.
		Difference(foundSet).
		Difference(sets.NewString(WalrusContextVariableName))
	if missingSet.Len() > 0 {
		return nil, fmt.Errorf("missing variables: %s", missingSet.List())
	}

	return variables, nil
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

	return getServiceDependencyOutputs(ctx, client, dependencyResourceIDs, dependOutputs)
}

// getServiceDependencyOutputs gets the dependency outputs of the resource.
func getServiceDependencyOutputs(
	ctx context.Context,
	client model.ClientSet,
	dependencyResourceIDs []object.ID,
	dependOutputs map[string]string,
) (map[string]parser.OutputState, error) {
	dependencyRevisions, err := client.ResourceRevisions().Query().
		Select(
			resourcerevision.FieldID,
			resourcerevision.FieldAttributes,
			resourcerevision.FieldOutput,
			resourcerevision.FieldResourceID,
			resourcerevision.FieldProjectID,
		).
		Where(func(s *sql.Selector) {
			sq := s.Clone().
				AppendSelectExprAs(
					sql.RowNumber().
						PartitionBy(resourcerevision.FieldResourceID).
						OrderBy(sql.Desc(resourcerevision.FieldCreateTime)),
					"row_number",
				).
				Where(s.P()).
				From(s.Table()).
				As(resourcerevision.Table)

			// Query the latest revision of the resource.
			s.Where(sql.EQ(s.C("row_number"), 1)).
				From(sq)
		}).
		Where(resourcerevision.ResourceIDIn(dependencyResourceIDs...)).
		WithResource(func(rq *model.ResourceQuery) {
			rq.Select(
				resource.FieldTemplateID,
				resource.FieldResourceDefinitionID,
			)
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make(map[string]parser.OutputState, 0)

	for _, r := range dependencyRevisions {
		revisionOutput, err := parser.ParseStateOutputRawMap(r)
		if err != nil {
			return nil, err
		}

		for n, o := range revisionOutput {
			if _, ok := dependOutputs[n]; !ok {
				continue
			}

			outputs[n] = o
		}
	}

	return outputs, nil
}
