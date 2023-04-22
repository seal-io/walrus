package platformtf

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

// connectorSeparator is used to separate the connector id and the instance name.
const connectorSeparator = "connector--"

type Provider = tfaddr.Provider

// AbsProviderConfig is the absolute address of a provider configuration
// within a particular module instance.
type AbsProviderConfig struct {
	Provider Provider
	Alias    string
}

type Parser struct{}

// ParseAppRevision parse the application revision output(terraform state) to application resources.
func (p Parser) ParseAppRevision(revision *model.ApplicationRevision) (model.ApplicationResources, error) {
	return p.ParseState(revision.Output, revision)
}

// ParseState parse the terraform state to application resources.
func (p Parser) ParseState(stateStr string, revision *model.ApplicationRevision) (model.ApplicationResources, error) {
	var logger = log.WithName("platformtf").WithName("parser")
	var revisionState state
	var applicationResources model.ApplicationResources
	if err := json.Unmarshal([]byte(stateStr), &revisionState); err != nil {
		return nil, err
	}

	for _, rs := range revisionState.Resources {
		switch rs.Mode {
		case "managed", "data":
		default:
			logger.Errorf("unknown resource mode: %s", rs.Mode)
			continue
		}

		// "module": "module.singleton[0]" or "module": "module.singleton"
		moduleName, err := ParseInstanceModuleName(rs.Module)
		if err != nil {
			logger.Errorf("invalid module format: %s", rs.Module)
			continue
		}
		// try to get the connector id from the provider.
		connector, err := ParseInstanceProviderConnector(rs.Provider)
		if err != nil {
			logger.Errorf("invalid provider format: %s", rs.Provider)
			continue
		}
		if connector == "" {
			logger.Warnf("connector is empty, provider: %v", rs.Provider)
			continue
		}

		for _, is := range rs.Instances {
			instanceID, err := ParseInstanceID(is)
			if err != nil {
				logger.Errorf("parse instance id failed: %w, instance: %v", err, is)
				continue
			}
			if instanceID == "" {
				logger.Errorf("instance id is empty, instance: %v", is)
				continue
			}

			applicationResource := &model.ApplicationResource{
				InstanceID:   revision.InstanceID,
				ConnectorID:  types.ID(connector),
				Mode:         rs.Mode,
				Module:       moduleName,
				Type:         rs.Type,
				Name:         instanceID,
				DeployerType: revision.DeployerType,
			}
			applicationResources = append(applicationResources, applicationResource)
		}
	}

	return applicationResources, nil
}

func ParseStateOutput(revision *model.ApplicationRevision) ([]types.OutputValue, error) {
	var revisionState state
	if err := json.Unmarshal([]byte(revision.Output), &revisionState); err != nil {
		return nil, err
	}

	if len(revision.Output) == 0 || revision.Status != status.ApplicationRevisionStatusSucceeded {
		return nil, nil
	}

	// sort by the module name length.
	var moduleNames = make([]string, len(revision.Modules))
	for i, v := range revision.Modules {
		moduleNames[i] = v.Name
	}
	sort.SliceStable(moduleNames, func(i, j int) bool {
		return len(moduleNames[i]) > len(moduleNames[j])
	})

	var outputs []types.OutputValue
	for n, o := range revisionState.Outputs {
		for _, mn := range moduleNames {
			if strings.Index(n, mn) == 0 {
				val := o.Value
				if o.Sensitive {
					val = []byte(`"<sensitive>"`)
				}
				outputs = append(outputs, types.OutputValue{
					Name:       strings.TrimPrefix(n, mn+"_"), // name format is moduleName_outputName.
					Value:      val,
					Type:       o.Type,
					Sensitive:  o.Sensitive,
					ModuleName: mn,
				})
				break
			}
		}
	}
	return outputs, nil
}

// ParseInstanceModuleName get the module name from the module instance string.
func ParseInstanceModuleName(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	traversal, parseDiags := hclsyntax.ParseTraversalAbs([]byte(str), "", hcl.Pos{Line: 1, Column: 1})
	if parseDiags.HasErrors() {
		return "", fmt.Errorf("invalid module format: %s", str)
	}

	var names []string
	for len(traversal) > 0 {
		var next string
		switch tt := traversal[0].(type) {
		case hcl.TraverseRoot:
			next = tt.Name
		case hcl.TraverseAttr:
			next = tt.Name
		}
		traversal = traversal[1:]
		if next != "module" {
			continue
		}

		if len(traversal) == 0 {
			return "", errors.New("prefix \"module.\" must be followed by a module name")
		}

		var moduleName string
		switch tt := traversal[0].(type) {
		case hcl.TraverseAttr:
			moduleName = tt.Name
		default:
			return "", errors.New("prefix \"module.\" must be followed by a module name")
		}
		traversal = traversal[1:]
		names = append(names, moduleName)
	}

	return strs.Join("/", names...), nil
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

	providers := strings.Split(providerConfig.Alias, connectorSeparator)
	if len(providers) != 2 {
		return "", fmt.Errorf("provider name error: %s", providerString)
	}

	return providers[1], nil
}

// ParseInstanceID get the real instance id from the instance object state.
// The instance id is stored in the "name" attribute of application resource
func ParseInstanceID(is instanceObjectState) (string, error) {
	if is.AttributesRaw != nil {
		ty, err := ctyjson.ImpliedType(is.AttributesRaw)
		if err != nil {
			return "", err
		}
		val, err := ctyjson.Unmarshal(is.AttributesRaw, ty)
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

// ParseStateProviders parse terraform state and get providers.
func ParseStateProviders(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}

	var providers = sets.NewString()
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

	var ret = &AbsProviderConfig{}
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
