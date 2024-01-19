package resource

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	pkgvariable "github.com/seal-io/walrus/pkg/variable"
	"github.com/seal-io/walrus/utils/json"
)

var (
	// VariableFormat the format of variable.
	variableFormat = `"${var.%s}"`
	// ResourceFormat the format of resource output.
	resourceFormat = `"${res.%s.%s}"`
)

func injectAttributes(
	ctx context.Context,
	client model.ClientSet,
	projectID, environmentID object.ID,
	attrs property.Values,
) (property.Values, error) {
	bs, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
	}

	varNames, resNames, resOuts := getInterpolationTargets(bs)

	if len(varNames) == 0 && len(resNames) == 0 {
		return attrs, nil
	}

	varValues, err := getInterpolationVariableValues(ctx, client, varNames, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	outputValues, err := getInterpolationOutputValues(ctx, client, resNames, projectID, environmentID, resOuts)
	if err != nil {
		return nil, err
	}

	return replaceAttributes(bs, varValues, outputValues)
}

func replaceAttributes(
	attrByte []byte,
	varValues map[string]json.RawMessage,
	outputValues map[string]property.Value,
) (property.Values, error) {
	for n, v := range varValues {
		attrByte = bytes.ReplaceAll(attrByte, []byte(n), v)
	}

	for n, v := range outputValues {
		attrByte = bytes.ReplaceAll(attrByte, []byte(n), v)
	}

	var injectAttrs property.Values
	if err := json.Unmarshal(attrByte, &injectAttrs); err != nil {
		return nil, err
	}

	return injectAttrs, nil
}

// getInterpolationVariableValues returns should inject variable values.
// Output is map of ${var.{variable_name}}:{value}.
func getInterpolationVariableValues(
	ctx context.Context,
	client model.ClientSet,
	varNames []string,
	projectID, environmentID object.ID,
) (map[string]json.RawMessage, error) {
	// Get should inject variables.
	vars, err := pkgvariable.Get(ctx, client, varNames, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	values := make(map[string]json.RawMessage)

	for _, v := range vars {
		n := fmt.Sprintf(variableFormat, v.Name)
		values[n] = []byte(fmt.Sprintf(`"%s"`, string(v.Value)))
	}

	return values, nil
}

// getInterpolationOutputValues returns should inject output values.
// Input resOutputs is map of {resource_name}_{output_name}:{inject_type}.
// Output is map of ${var.{resource_name}.{output_name}}:{value}.
func getInterpolationOutputValues(
	ctx context.Context,
	client model.ClientSet,
	resNames []string,
	projectID, environmentID object.ID,
	resOutputs map[string]string,
) (map[string]property.Value, error) {
	depIDs, err := client.Resources().Query().
		Where(
			resource.NameIn(resNames...),
			resource.ProjectID(projectID),
			resource.EnvironmentID(environmentID)).
		IDs(ctx)
	if err != nil {
		return nil, err
	}

	// Get should inject resource outputs.
	outputs, err := pkgresource.GetDependencyOutputs(ctx, client, depIDs, resOutputs)
	if err != nil {
		return nil, err
	}

	values := make(map[string]property.Value)

	for name := range resOutputs {
		v, ok := outputs[name]
		if !ok {
			continue
		}

		ss := strings.SplitN(name, "_", 2)
		if len(ss) != 2 {
			continue
		}

		n := fmt.Sprintf(resourceFormat, ss[0], ss[1])
		values[n] = v.Value
	}

	return values, nil
}

func getInterpolationTargets(
	attrByte []byte,
) ([]string, []string, map[string]string) {
	variableMatches := pkgresource.VariableReg.FindAllSubmatch(attrByte, -1)
	resourceMatches := pkgresource.ResourceReg.FindAllSubmatch(attrByte, -1)

	var (
		variableMatched = sets.NewString()
		resourceMatched = sets.NewString()
		resourceOutput  = make(map[string]string)
	)

	for _, match := range variableMatches {
		if len(match) > 1 {
			variableMatched.Insert(string(match[1]))
		}
	}

	for _, match := range resourceMatches {
		if len(match) > 1 {
			name, outputName := string(match[1]), string(match[2])

			resourceMatched.Insert(name)

			// This format follow the input of terraform.getResourceDependencyOutputsByID.
			resourceOutput[fmt.Sprintf("%s_%s", name, outputName)] = "res"
		}
	}

	return variableMatched.List(), resourceMatched.List(), resourceOutput
}

func createInputsItemToResource(
	input *model.ResourceCreateInputsItem,
	p *model.ProjectQueryInput,
	e *model.EnvironmentQueryInput,
) *model.Resource {
	return toResource(
		nil, input.Name, input.Type, input.Attributes, input.Labels,
		input.Template,
		input.ResourceDefinition, input.ResourceDefinitionMatchingRule,
		p, e)
}

func updateInputsItemToResource(
	typ string,
	input *model.ResourceUpdateInputsItem,
	p *model.ProjectQueryInput,
	e *model.EnvironmentQueryInput,
) *model.Resource {
	return toResource(
		nil, input.Name, typ, input.Attributes, input.Labels,
		input.Template,
		input.ResourceDefinition, input.ResourceDefinitionMatchingRule,
		p, e)
}

func toResource(
	id *object.ID, name, typ string, attr property.Values, labels map[string]string,
	tmpl *model.TemplateVersionQueryInput,
	rd *model.ResourceDefinitionQueryInput, rdr *model.ResourceDefinitionMatchingRuleQueryInput,
	p *model.ProjectQueryInput, e *model.EnvironmentQueryInput,
) *model.Resource {
	r := &model.Resource{
		Name:       name,
		Labels:     labels,
		Type:       typ,
		Attributes: attr,
	}

	if id != nil {
		r.ID = *id
	}

	if tmpl != nil {
		r.TemplateID = &tmpl.ID
	}

	if rd != nil {
		r.ResourceDefinitionID = &rd.ID
	}

	if rdr != nil {
		r.ResourceDefinitionMatchingRuleID = &rdr.ID
	}

	if p != nil {
		r.ProjectID = p.ID
	}

	if e != nil {
		r.EnvironmentID = e.ID
	}

	return r
}
