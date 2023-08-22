package terraform

import (
	"context"
	"fmt"
	"reflect"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/model/servicerelationship"
	"github.com/seal-io/walrus/pkg/dao/model/servicerevision"
	"github.com/seal-io/walrus/pkg/dao/model/variable"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/terraform/parser"
)

type ServiceOpts struct {
	ServiceRevision *model.ServiceRevision

	ServiceName string

	ProjectID     object.ID
	EnvironmentID object.ID
}

// ParseModuleAttributes parse module variables and dependencies.
func ParseModuleAttributes(
	ctx context.Context,
	mc model.ClientSet,
	attributes map[string]any,
	onlyValidated bool,
	opts ServiceOpts,
) (variables model.Variables, outputs map[string]parser.OutputState, err error) {
	var (
		templateVariables        []string
		dependencyServiceOutputs []string
	)

	replaced := !onlyValidated
	templateVariables, dependencyServiceOutputs = parseAttributeReplace(
		attributes,
		templateVariables,
		dependencyServiceOutputs,
		replaced,
	)

	// If service revision has variables that inherit from cloned revision, use them directly.
	if opts.ServiceRevision != nil && len(opts.ServiceRevision.Variables) > 0 {
		for k, v := range opts.ServiceRevision.Variables {
			variables = append(variables, &model.Variable{
				Name:  k,
				Value: crypto.String(v),
			})
		}
	} else {
		variables, err = getVariables(ctx, mc, templateVariables, opts.ProjectID, opts.EnvironmentID)
		if err != nil {
			return nil, nil, err
		}
	}

	if !onlyValidated {
		outputs, err = getServiceDependencyOutputsByID(
			ctx,
			mc,
			opts.ServiceRevision.ServiceID,
			dependencyServiceOutputs)
		if err != nil {
			return nil, nil, err
		}

		// Check if all dependency service outputs are found.
		for _, o := range dependencyServiceOutputs {
			if _, ok := outputs[o]; !ok {
				return nil, nil, fmt.Errorf("service %s dependency output %s not found", opts.ServiceName, o)
			}
		}
	}

	return variables, outputs, nil
}

// parseAttributeReplace parses attribute variable ${var.name} replaces it with ${var._variablePrefix+name},
// service reference ${service.name.output} replaces it with ${var._servicePrefix+name}
// and returns variable names and service names.
// Replaced flag indicates whether to replace the module attribute variable with prefix string.
func parseAttributeReplace(
	attributes map[string]any,
	variableNames []string,
	serviceOutputs []string,
	replaced bool,
) ([]string, []string) {
	for key, value := range attributes {
		if value == nil {
			continue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			if _, ok := value.(map[string]any); !ok {
				continue
			}

			variableNames, serviceOutputs = parseAttributeReplace(
				value.(map[string]any),
				variableNames,
				serviceOutputs,
				replaced,
			)
		case reflect.String:
			str := value.(string)
			matches := _variableReg.FindAllStringSubmatch(str, -1)
			serviceMatches := _serviceReg.FindAllStringSubmatch(str, -1)

			var matched []string

			for _, match := range matches {
				if len(match) > 1 {
					matched = append(matched, match[1])
				}
			}

			var serviceMatched []string

			for _, match := range serviceMatches {
				if len(match) > 1 {
					serviceMatched = append(serviceMatched, match[1]+"_"+match[2])
				}
			}

			variableNames = append(variableNames, matched...)
			variableRepl := "${var." + _variablePrefix + "${1}}"
			str = _variableReg.ReplaceAllString(str, variableRepl)

			serviceOutputs = append(serviceOutputs, serviceMatched...)
			serviceRepl := "${var." + _servicePrefix + "${1}_${2}}"

			if replaced {
				attributes[key] = _serviceReg.ReplaceAllString(str, serviceRepl)
			}
		case reflect.Slice:
			if _, ok := value.([]any); !ok {
				continue
			}

			for _, v := range value.([]any) {
				if _, ok := v.(map[string]any); !ok {
					continue
				}
				variableNames, serviceOutputs = parseAttributeReplace(
					v.(map[string]any),
					variableNames,
					serviceOutputs,
					replaced,
				)
			}
		}
	}

	return variableNames, serviceOutputs
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
		Difference(WalrusMetadataSet)
	if missingSet.Len() > 0 {
		return nil, fmt.Errorf("missing variables: %s", missingSet.List())
	}

	return variables, nil
}

// getServiceDependencyOutputsByID gets the dependency outputs of the service by service id.
func getServiceDependencyOutputsByID(
	ctx context.Context,
	client model.ClientSet,
	serviceID object.ID,
	dependOutputs []string,
) (map[string]parser.OutputState, error) {
	entity, err := client.Services().Query().
		Where(service.ID(serviceID)).
		WithDependencies(func(sq *model.ServiceRelationshipQuery) {
			sq.Where(func(s *sql.Selector) {
				s.Where(sql.ColumnsNEQ(servicerelationship.FieldServiceID, servicerelationship.FieldDependencyID))
			})
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dependencyServiceIDs := make([]object.ID, 0, len(entity.Edges.Dependencies))

	for _, d := range entity.Edges.Dependencies {
		if d.Type != types.ServiceRelationshipTypeImplicit {
			continue
		}

		dependencyServiceIDs = append(dependencyServiceIDs, d.DependencyID)
	}

	return getServiceDependencyOutputs(ctx, client, dependencyServiceIDs, dependOutputs)
}

// getServiceDependencyOutputs gets the dependency outputs of the service.
func getServiceDependencyOutputs(
	ctx context.Context,
	client model.ClientSet,
	dependencyServiceIDs []object.ID,
	dependOutputs []string,
) (map[string]parser.OutputState, error) {
	dependencyRevisions, err := client.ServiceRevisions().Query().
		Select(
			servicerevision.FieldID,
			servicerevision.FieldAttributes,
			servicerevision.FieldOutput,
			servicerevision.FieldServiceID,
			servicerevision.FieldProjectID,
		).
		Where(func(s *sql.Selector) {
			sq := s.Clone().
				AppendSelectExprAs(
					sql.RowNumber().
						PartitionBy(servicerevision.FieldServiceID).
						OrderBy(sql.Desc(servicerevision.FieldCreateTime)),
					"row_number",
				).
				Where(s.P()).
				From(s.Table()).
				As(servicerevision.Table)

			// Query the latest revision of the service.
			s.Where(sql.EQ(s.C("row_number"), 1)).
				From(sq)
		}).Where(servicerevision.ServiceIDIn(dependencyServiceIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make(map[string]parser.OutputState, 0)
	dependSets := sets.NewString(dependOutputs...)

	for _, r := range dependencyRevisions {
		revisionOutput, err := parser.ParseStateOutputRawMap(r)
		if err != nil {
			return nil, err
		}

		for n, o := range revisionOutput {
			if dependSets.Has(n) {
				outputs[n] = o
			}
		}
	}

	return outputs, nil
}
