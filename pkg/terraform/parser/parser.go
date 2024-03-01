package parser

import (
	"context"
	"errors"
	"fmt"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcestate"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/templates/translator"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// ConnectorSeparator is used to separate the connector id and the instance name.
const ConnectorSeparator = "connector--"

// The type terraform data implements the standard resource lifecycle,
// but does not directly take any other actions.
// Resource components should skip the data type.
const TerraformTypeData = "terraform_data"

type StateParser struct{}

// GetComponentsAndExtractDependencies returns the components and dependency components after parse the resource state.
//
// GetComponentsAndExtractDependencies returns list must not be `nil` unless unexpected input or raising error,
// it can be used to clean stale items safety if got an empty list.
func (StateParser) GetComponentsAndExtractDependencies(
	ctx context.Context,
	mc model.ClientSet,
	run *model.ResourceRun,
) (
	components model.ResourceComponents,
	dependencies map[string][]string,
	err error,
) {
	logger := log.WithName("deployer").WithName("tf").WithName("parser")
	dependencies = make(map[string][]string)

	s, err := mc.ResourceStates().Query().
		Where(resourcestate.ResourceID(run.ResourceID)).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	var runState state
	if err := json.Unmarshal([]byte(s.Data), &runState); err != nil {
		return nil, nil, err
	}

	var (
		// Maps components unique index to its dependencies components unique indexes.
		componentDependencies = make(map[string][]string)
		// Maps terraform module key to resource.
		moduleComponentMap = make(map[string]*model.ResourceComponent)
		key                = dao.ResourceComponentGetUniqueKey
	)

	for _, rs := range runState.Resources {
		if rs.Type == TerraformTypeData {
			continue
		}

		switch rs.Mode {
		default:
			logger.Errorf("unknown resource mode: %s", rs.Mode)
			continue
		case types.ResourceComponentModeManaged, types.ResourceComponentModeData:
		}

		// Try to get the connectorID id from the provider.
		connectorID, err := ParseInstanceProviderConnector(rs.Provider)
		if err != nil {
			logger.Errorf("invalid provider format: %s", rs.Provider)
			continue
		}

		if connectorID == "" {
			logger.Warnf("connector is empty, provider: %v", rs.Provider)
			continue
		}

		classResourceComponents := &model.ResourceComponent{
			ProjectID:     run.ProjectID,
			EnvironmentID: run.EnvironmentID,
			ResourceID:    run.ResourceID,
			ConnectorID:   object.ID(connectorID),
			Mode:          rs.Mode,
			Type:          rs.Type,
			Name:          rs.Name,
			DeployerType:  run.DeployerType,
			Shape:         types.ResourceComponentShapeClass,
		}
		classResourceComponents.Edges.Instances = make(model.ResourceComponents, len(rs.Instances))

		// The module key is used to identify the terraform resource module.
		mk := strs.Join(".", rs.Module, rs.Type, rs.Name)
		if rs.Mode == types.ResourceComponentModeData {
			mk = strs.Join(".", rs.Module, rs.Mode, rs.Type, rs.Name)
		}
		moduleComponentMap[mk] = classResourceComponents

		for i, is := range rs.Instances {
			instanceID, err := ParseInstanceID(rs, is)
			if err != nil {
				logger.Errorf("parse instance id failed: %v, instance: %v",
					err, is)
				continue
			}

			if instanceID == "" {
				logger.Errorf("instance id is empty, instance: %v", is)
				continue
			}

			// The index key is used to identify the terraform resource instance.
			indexKey, err := ParseIndexKey(rs, is)
			if err != nil {
				logger.Errorf("parse index key failed: %v, instance: %v", err, is)
				continue
			}

			// FIXME(thxCode): as a good solution,
			//  the https://registry.terraform.io/providers/hashicorp/helm should provide a complete ID.
			if rs.Type == "helm_release" && !strings.Contains(instanceID, "/") {
				// NB(thxCode): the ID of helm_release resource doesn't include namespace,
				// so we can't fetch the real Helm Release record that under specified namespace.
				// In order to recognize the real Helm Release record,
				// we should enrich the instanceID with the namespace name.
				md, err := ParseInstanceMetadata(is)
				if err != nil {
					logger.Errorf("parse instance metadata failed: %v, instance attributes: %s",
						err, string(is.Attributes))
					continue
				}

				if nsr := json.Get(md, "namespace"); nsr.String() != "" {
					instanceID = nsr.String() + "/" + instanceID
				} else {
					instanceID = core.NamespaceDefault + "/" + instanceID
				}
			}

			instanceResource := &model.ResourceComponent{
				ProjectID:     run.ProjectID,
				EnvironmentID: run.EnvironmentID,
				ResourceID:    run.ResourceID,
				ConnectorID:   object.ID(connectorID),
				Mode:          rs.Mode,
				Type:          rs.Type,
				Name:          instanceID,
				Shape:         types.ResourceComponentShapeInstance,
				DeployerType:  run.DeployerType,
				IndexKey:      indexKey,
			}

			// Assume that the first instance's dependencies are the dependencies of the class resource.
			if _, ok := moduleComponentMap[key(classResourceComponents)]; !ok {
				componentDependencies[key(classResourceComponents)] = is.Dependencies
			}

			dependencies[key(instanceResource)] = append(
				dependencies[key(instanceResource)],
				key(classResourceComponents),
			)
			classResourceComponents.Edges.Instances[i] = instanceResource
			componentDependencies[key(instanceResource)] = is.Dependencies
		}

		components = append(components, classResourceComponents)
	}

	// Get resource dependencies.
	for k, v := range componentDependencies {
		for _, d := range v {
			moduleResource, ok := moduleComponentMap[d]
			if !ok {
				logger.Warnf("dependency resource not found, module key: %s", d)
				continue
			}

			dependencies[k] = append(dependencies[k], key(moduleResource))
		}
	}

	return components, dependencies, nil
}

// GetOutputMap returns the original outputs after parsed the resource run output(terraform state).
//
// Since we mutate the output names before executing a terraform deployment,
// the output's name(hcl label) is not the same as the original one defined on the terraform template.
//
// This function is used for bridging the referring between multiple (walrus)resources.
// Use GetOriginalOutputsMap if wanna the original outputs.
func (StateParser) GetOutputMap(stateData string) (map[string]types.OutputValue, error) {
	if len(stateData) == 0 {
		return nil, nil
	}

	// Get outputs from state, expected format:
	// {
	//   "outputs": {}
	// }.
	r := json.Get(strs.ToBytes(&stateData), "outputs")
	if !r.Exists() || !r.IsObject() {
		return map[string]types.OutputValue{}, nil
	}

	var osm map[string]types.OutputValue
	if err := json.Unmarshal(strs.ToBytes(&r.Raw), &osm); err != nil {
		return nil, err
	}

	return osm, nil
}

// GetOriginalOutputs returns the original outputs after parsed the resource run output(terraform state).
//
// The given run must carry the resource on the edges, especially the resource's name.
//
// This function returns the original outputs,
// which means the output's name(hcl label) is the same as the original one defined on the terraform template.
func (p StateParser) GetOriginalOutputs(stateData, resourceName string) ([]types.OutputValue, error) {
	osm, err := p.GetOutputMap(stateData)
	if err != nil {
		return nil, err
	}

	var (
		prefix = resourceName + "_"
		oss    = make([]types.OutputValue, 0, len(osm))
		count  int
	)

	for _, mn := range sets.StringKeySet(osm).List() {
		// E.g. `n` is in the form of `{resource name}_{output name}`.
		n := strings.TrimPrefix(mn, prefix)
		if n == mn {
			continue
		}
		o := osm[mn]

		count++

		s := translator.SchemaOfType(
			o.Type,
			translator.Options{
				Name:      n,
				Sensitive: o.Sensitive,
				Order:     count,
			})

		v := o.Value
		if o.Sensitive {
			v = []byte(`"<sensitive>"`)
		}

		oss = append(oss, types.OutputValue{
			Name:   n,
			Value:  v,
			Type:   o.Type,
			Schema: s,
		})
	}

	return oss, nil
}

// GetOriginalOutputsMap is similar to GetOriginalOutputs,
// but returns the original outputs in map form.
func (p StateParser) GetOriginalOutputsMap(stateData, resourceName string) (map[string]types.OutputValue, error) {
	oss, err := p.GetOriginalOutputs(stateData, resourceName)
	if err != nil {
		return nil, err
	}

	osm := make(map[string]types.OutputValue, len(oss))
	for i := range oss {
		osm[oss[i].Name] = oss[i]
	}

	return osm, nil
}

type Provider = tfaddr.Provider

// AbsProviderConfig is the absolute address of a provider configuration
// within a particular module instance.
type AbsProviderConfig struct {
	Provider Provider
	Alias    string
}

// ParseInstanceProviderConnector get the provider connector from the provider instance string.
func ParseInstanceProviderConnector(providerString string) (string, error) {
	providerConfig, err := ParseAbsProviderString(providerString)
	if err != nil {
		return "", err
	}

	if providerConfig.Alias == "" {
		return "", nil
	}

	providers := strings.Split(providerConfig.Alias, ConnectorSeparator)
	if len(providers) != 2 {
		return "", fmt.Errorf("provider name error: %s", providerString)
	}

	return providers[1], nil
}

// ParseInstanceID get the real instance id from the instance object state.
// The instance id is stored in the "name" attribute of resource component.
func ParseInstanceID(rs resourceState, is instanceObjectState) (string, error) {
	if is.Attributes != nil {
		ty, err := ctyjson.ImpliedType(is.Attributes)
		if err != nil {
			return "", err
		}

		val, err := ctyjson.Unmarshal(is.Attributes, ty)
		if err != nil {
			return "", err
		}

		for key, value := range val.AsValueMap() {
			if key == "id" {
				if value.IsNull() {
					return "", nil
				}

				switch value.Type() {
				case cty.String:
					return value.AsString(), nil
				case cty.Number:
					return value.AsBigFloat().String(), nil
				default:
					return "", fmt.Errorf("unsupported type for id: %s, value: %s", value, value.Type().FriendlyName())
				}
			}
		}
	}

	if is.AttributesFlat != nil {
		if id, ok := is.AttributesFlat["id"]; ok {
			return id, nil
		}
	}

	return ParseIndexKey(rs, is)
}

// ParseInstanceMetadata get the metadata from the instance object state.
func ParseInstanceMetadata(is instanceObjectState) ([]byte, error) {
	if is.Attributes == nil {
		return nil, errors.New("no attributes")
	}

	arr := json.Get(is.Attributes, "metadata").Array()

	switch l := len(arr); {
	case l == 0:
		return nil, errors.New("not found metadata")
	case l > 1:
		return nil, errors.New("not singular metadata")
	}

	if !arr[0].IsObject() {
		return nil, errors.New("metadata is not an object")
	}

	return strs.ToBytes(&arr[0].Raw), nil
}

// ParseStateProviders parse terraform state and get providers.
func ParseStateProviders(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}

	providers := sets.NewString()

	var runState state
	if err := json.Unmarshal([]byte(s), &runState); err != nil {
		return nil, err
	}

	for _, resource := range runState.Resources {
		pAddr, err := ParseAbsProviderString(resource.Provider)
		if err != nil {
			return nil, err
		}

		providers.Insert(pAddr.Provider.Type)
	}

	return providers.List(), nil
}

func parseAbsProvider(traversal hcl.Traversal) (hcl.Traversal, error) {
	remain := traversal

	for len(remain) > 0 {
		var next string
		switch tt := remain[0].(type) {
		case hcl.TraverseRoot:
			next = tt.Name
		case hcl.TraverseAttr:
			next = tt.Name
		case hcl.TraverseIndex:
			return nil, errors.New("provider address cannot contain module indexes")
		}

		if next != "provider" {
			remain = remain[1:]
			continue
		}

		var retRemain hcl.Traversal
		if len(remain) > 0 {
			retRemain = make(hcl.Traversal, len(remain))
			copy(retRemain, remain)

			if tt, ok := retRemain[0].(hcl.TraverseAttr); ok {
				retRemain[0] = hcl.TraverseRoot{
					Name:     tt.Name,
					SrcRange: tt.SrcRange,
				}
			}

			return retRemain, nil
		}
	}

	return nil, fmt.Errorf("invalid provider configuration address %q", traversal)
}

// ParseAbsProviderConfig parses the given traversal as an absolute provider configuration address.
func ParseAbsProviderConfig(traversal hcl.Traversal) (*AbsProviderConfig, error) {
	remain, err := parseAbsProvider(traversal)
	if err != nil {
		return nil, err
	}

	if len(remain) < 2 || remain.RootName() != "provider" {
		return nil, errors.New("provider address must begin with \"provider.\", followed by a provider type name")
	}

	if len(remain) > 3 {
		return nil, errors.New("extraneous operators after provider configuration alias")
	}

	ret := &AbsProviderConfig{}

	if tt, ok := remain[1].(hcl.TraverseIndex); ok {
		if !tt.Key.Type().Equals(cty.String) {
			return nil, errors.New("the prefix \"provider.\" must be followed by a provider type name")
		}

		p, err := tfaddr.ParseProviderSource(tt.Key.AsString())
		if err != nil {
			return nil, err
		}
		ret.Provider = p
	} else {
		return nil, errors.New("the prefix \"provider.\" must be followed by a provider type name")
	}

	if len(remain) == 3 {
		if tt, ok := remain[2].(hcl.TraverseAttr); ok {
			ret.Alias = tt.Name
		} else {
			return nil, errors.New("provider type name must be followed by a configuration alias name")
		}
	}

	return ret, nil
}

func ParseAbsProviderString(str string) (*AbsProviderConfig, error) {
	traversal, diags := hclsyntax.ParseTraversalAbs([]byte(str), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("invalid provider configuration address %s", str)
	}

	ret, err := ParseAbsProviderConfig(traversal)
	if err != nil {
		return nil, fmt.Errorf("invalid provider configuration address %q: %w", str, err)
	}

	return ret, nil
}

// ParseIndexKey parse the index key from the instance object state.
// The index key is used to identify the terraform resource instance, e.g. `helm_release.foo[0]`.
func ParseIndexKey(rs resourceState, is instanceObjectState) (string, error) {
	if rs.Type == "" || rs.Name == "" {
		return "", errors.New("resource type or name is empty")
	}

	if is.IndexKey == nil {
		return strs.Join(".", rs.Module, rs.Type, rs.Name), nil
	}

	switch is.IndexKey.(type) {
	case string:
		return strs.Join(".", rs.Module, rs.Type, rs.Name) + fmt.Sprintf("[\"%v\"]", is.IndexKey), nil
	case int:
		return strs.Join(".", rs.Module, rs.Type, rs.Name) + fmt.Sprintf("[%v]", is.IndexKey), nil
	}

	return "", fmt.Errorf("invalid index key: %v", is.IndexKey)
}
