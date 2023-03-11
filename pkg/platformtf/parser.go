package platformtf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/types"
	tftopic "github.com/seal-io/seal/pkg/topic/platformtf"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

// connectorSeparator is used to separate the connector id and the instance name.
const connectorSeparator = "connector--"

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
		// try to get "singleton" from module
		connector, err := ParseInstanceProviderConnector(rs.Provider)
		if err != nil {
			logger.Errorf("invalid provider format: %s", rs.Provider)
			continue
		}

		for _, is := range rs.Instances {
			instanceID, err := ParseInstanceID(is)
			if err != nil {
				return nil, err
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

// ParseInstanceModuleName get the module name from the module instance string.
func ParseInstanceModuleName(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	traversal, parseDiags := hclsyntax.ParseTraversalAbs([]byte(str), "", hcl.Pos{Line: 1, Column: 1})
	if parseDiags.HasErrors() {
		return "", fmt.Errorf("invalid module format: %s", str)
	}

	var name string
	for _, t := range traversal {
		switch tt := t.(type) {
		case hcl.TraverseAttr:
			name = tt.Name
		default:
			continue
		}
	}

	return name, nil
}

// ParseInstanceProviderConnector get the provider connector from the provider instance string.
func ParseInstanceProviderConnector(providerString string) (string, error) {
	providers := strings.Split(providerString, connectorSeparator)
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
				return value.AsString(), nil
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

func UpdateResource(ctx context.Context, message tftopic.TopicMessage) error {
	var parser Parser
	applicationResources, err := parser.ParseAppRevision(message.ApplicationRevision)
	if err != nil {
		return err
	}

	return message.ModelClient.WithTx(ctx, func(tx *model.Tx) error {
		var (
			existResourceIDs = make([]types.ID, 0)
			newResources     = make(model.ApplicationResources, 0)
		)

		// fetch the old resources of the application
		oldResources, err := message.ModelClient.ApplicationResources().
			Query().
			Where(applicationresource.InstanceID(message.ApplicationRevision.InstanceID)).
			All(ctx)
		if err != nil {
			return err
		}
		oldResourceSet := sets.NewString()
		for _, r := range oldResources {
			uniqueKey := getFingerprint(r)
			oldResourceSet.Insert(uniqueKey)
		}

		for _, ar := range applicationResources {
			// check if the resource is exists.
			key := getFingerprint(ar)
			exists := oldResourceSet.Has(key)
			if exists {
				existResourceIDs = append(existResourceIDs, ar.ID)
			} else {
				newResources = append(newResources, ar)
			}
		}

		// diff application resource of this revision and the latest revision.
		// if the resource is not in the latest revision, delete it.
		_, err = message.ModelClient.ApplicationResources().
			Delete().
			Where(
				applicationresource.InstanceID(message.ApplicationRevision.InstanceID),
				applicationresource.IDNotIn(existResourceIDs...),
			).
			Exec(ctx)
		if err != nil {
			return err
		}

		// create newResource.
		if len(newResources) > 0 {
			resourcesToCreate, err := dao.ApplicationResourceCreates(message.ModelClient, newResources...)
			if err != nil {
				return err
			}
			if _, err = message.ModelClient.ApplicationResources().CreateBulk(resourcesToCreate...).Save(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

// TODO(thxCode): generate by entc.
func getFingerprint(r *model.ApplicationResource) string {
	// align to schema definition.
	return strs.Join("-", string(r.ConnectorID), r.Module, r.Mode, r.Type, r.Name)
}
