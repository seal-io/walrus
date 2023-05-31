package applicationresources

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/operator"
	"github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
)

type ApplicationResourceDetail struct {
	Resource *model.ApplicationResourceOutput `json:"resource"`
	Keys     *types.Keys                      `json:"keys,omitempty"`
}

// MarshalJSON implements the json.Marshaler to avoid the impact from model.ApplicationResourceOutput's marshaller.
func (in ApplicationResourceDetail) MarshalJSON() ([]byte, error) {
	type (
		AliasResource model.ApplicationResourceOutput
		Alias         struct {
			*AliasResource `json:",inline"`
			Keys           *types.Keys `json:"keys"`
		}
	)

	return json.Marshal(&Alias{
		AliasResource: (*AliasResource)(in.Resource.Normalize()),
		Keys:          in.Keys,
	})
}

func GetResourcesDetail(
	ctx context.Context,
	resources model.ApplicationResources,
	withoutKeys bool,
) []ApplicationResourceDetail {
	logger := log.WithName("application-resources")

	details := make([]ApplicationResourceDetail, len(resources))
	for i := 0; i < len(resources); i++ {
		details[i].Resource = model.ExposeApplicationResource(resources[i])
	}

	if !withoutKeys {
		// NB(thxCode): we can safety index the connector with its pointer here,
		// as the ent can keep the connector pointer is the same between those resources related by the same connector.
		m := make(map[*model.Connector][]int)
		for i := 0; i < len(resources); i++ {
			m[resources[i].Edges.Connector] = append(m[resources[i].Edges.Connector], i)
		}

		for c, idxs := range m {
			// Get operator by connector.
			op, err := operator.Get(ctx, types.CreateOptions{Connector: *c})
			if err != nil {
				logger.Warnf("cannot get operator of connector: %v", err)
				continue
			}

			if err = op.IsConnected(ctx); err != nil {
				logger.Warnf("unreachable connector: %v", err)
				continue
			}
			// Fetch keys for the resources that related to same connector.
			for _, i := range idxs {
				details[i].Keys, err = op.GetKeys(ctx, resources[i])
				if err != nil {
					logger.Errorf("error getting keys: %v", err)
				}
			}
		}
	}

	return details
}
