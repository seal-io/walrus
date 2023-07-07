package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

// ConnectorSeparator is used to separate the connector id and the instance name.
const ConnectorSeparator = "connector--"

type Provider = tfaddr.Provider

// AbsProviderConfig is the absolute address of a provider configuration
// within a particular module instance.
type AbsProviderConfig struct {
	Provider Provider
	Alias    string
}

type Parser struct{}

// ParseServiceRevision parse the service revision output(terraform state) to service resources,
// returns list must not be `nil` unless unexpected input or raising error,
// it can be used to clean stale items safety if got an empty list.
func (p Parser) ParseServiceRevision(revision *model.ServiceRevision) (
	model.ServiceResources, map[string][]string, error,
) {
	return p.ParseState(revision.Output, revision)
}

// ParseState parse the terraform state to service resources,
// returns list must not be `nil` unless unexpected input or raising error,
// it can be used to clean stale items safety if got an empty list.
func (p Parser) ParseState(stateStr string, revision *model.ServiceRevision) (
	applicationResources model.ServiceResources, dependencies map[string][]string, err error,
) {
	logger := log.WithName("deployer").WithName("tf").WithName("parser")
	dependencies = make(map[string][]string)

	var revisionState state
	if err := json.Unmarshal([]byte(stateStr), &revisionState); err != nil {
		return nil, nil, err
	}

	var (
		// ModuleDependencies maps resource unique index to its dependencies resource unique indexes.
		resourceDependencies = make(map[string][]string)
		// ModuleResourceMap maps terraform module key to resource.
		moduleResourceMap = make(map[string]*model.ServiceResource)
		key               = dao.ServiceResourceGetUniqueKey
	)

	for _, rs := range revisionState.Resources {
		switch rs.Mode {
		default:
			logger.Errorf("unknown resource mode: %s", rs.Mode)
			continue
		case types.ServiceResourceModeManaged, types.ServiceResourceModeData:
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

		applicationResource := &model.ServiceResource{
			ServiceID:    revision.ServiceID,
			ConnectorID:  oid.ID(connectorID),
			Mode:         rs.Mode,
			Type:         rs.Type,
			Name:         rs.Name,
			DeployerType: revision.DeployerType,
			Shape:        types.ServiceResourceShapeClass,
		}
		applicationResource.Edges.Instances = make(model.ServiceResources, len(rs.Instances))

		// The module key is used to identify the terraform resource module.
		mk := strs.Join(".", rs.Module, rs.Type, rs.Name)
		if rs.Mode == types.ServiceResourceModeData {
			mk = strs.Join(".", rs.Module, rs.Mode, rs.Type, rs.Name)
		}
		moduleResourceMap[mk] = applicationResource

		for i, is := range rs.Instances {
			instanceID, err := ParseInstanceID(is)
			if err != nil {
				logger.Errorf("parse instance id failed: %v, instance: %v",
					err, is)
				continue
			}

			if instanceID == "" {
				logger.Errorf("instance id is empty, instance: %v", is)
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
			resource := &model.ServiceResource{
				ServiceID:    revision.ServiceID,
				ConnectorID:  oid.ID(connectorID),
				Mode:         rs.Mode,
				Type:         rs.Type,
				Name:         instanceID,
				Shape:        types.ServiceResourceShapeInstance,
				DeployerType: revision.DeployerType,
			}

			// Assume that the first instance's dependencies are the dependencies of the class resource.
			if _, ok := moduleResourceMap[key(applicationResource)]; !ok {
				resourceDependencies[key(applicationResource)] = is.Dependencies
			}

			dependencies[key(resource)] = append(dependencies[key(resource)], key(applicationResource))
			applicationResource.Edges.Instances[i] = resource
			resourceDependencies[key(resource)] = is.Dependencies
		}

		applicationResources = append(applicationResources, applicationResource)
	}

	// Get resource dependencies.
	for k, v := range resourceDependencies {
		for _, d := range v {
			moduleResource, ok := moduleResourceMap[d]
			if !ok {
				logger.Warnf("dependency resource not found, module key: %s", d)
				continue
			}

			dependencies[k] = append(dependencies[k], key(moduleResource))
		}
	}

	return applicationResources, dependencies, nil
}

func ParseStateOutputRawMap(revision *model.ServiceRevision) (map[string]OutputState, error) {
	if len(revision.Output) == 0 {
		return nil, nil
	}

	var revisionState state
	if err := json.Unmarshal([]byte(revision.Output), &revisionState); err != nil {
		return nil, err
	}

	return revisionState.Outputs, nil
}

func ParseStateOutput(revision *model.ServiceRevision) ([]types.OutputValue, error) {
	if len(revision.Output) == 0 {
		return nil, nil
	}

	var revisionState state
	if err := json.Unmarshal([]byte(revision.Output), &revisionState); err != nil {
		return nil, err
	}

	sn := revision.Edges.Service.Name

	var outputs []types.OutputValue

	for n, o := range revisionState.Outputs {
		if strings.Index(n, sn) == 0 {
			val := o.Value
			if o.Sensitive {
				val = []byte(`"<sensitive>"`)
			}

			outputs = append(outputs, types.OutputValue{
				Name:      strings.TrimPrefix(n, sn+"_"), // Name format is serviceName_outputName.
				Value:     val,
				Type:      o.Type,
				Sensitive: o.Sensitive,
			})
		}
	}

	return outputs, nil
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
// The instance id is stored in the "name" attribute of application resource.
func ParseInstanceID(is instanceObjectState) (string, error) {
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

	return "", fmt.Errorf("no id found in instance object state: %v", is)
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

	var revisionState state
	if err := json.Unmarshal([]byte(s), &revisionState); err != nil {
		return nil, err
	}

	for _, resource := range revisionState.Resources {
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
