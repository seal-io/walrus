package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

func ConnectorCreates(mc model.ClientSet, input ...*model.Connector) ([]*model.ConnectorCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ConnectorCreate, len(input))

	for i := range input {
		r := input[i]
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Connectors().Create().
			SetName(r.Name).
			SetType(r.Type).
			SetConfigVersion(r.ConfigVersion).
			SetConfigData(r.ConfigData).
			SetEnableFinOps(r.EnableFinOps).
			SetCategory(r.Category)

		// Optional.
		c.SetDescription(r.Description)

		if r.ProjectID.IsNaive() {
			c.SetProjectID(r.ProjectID)
		}

		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}

		if !r.FinOpsCustomPricing.IsZero() {
			c.SetFinOpsCustomPricing(r.FinOpsCustomPricing)
		} else if r.Type == types.ConnectorTypeK8s {
			// Set default pricing for Kubernetes connector.
			c.SetFinOpsCustomPricing(types.DefaultFinOpsCustomPricing())
		}

		if r.Type == types.ConnectorTypeK8s {
			if r.EnableFinOps {
				status.ConnectorStatusCostToolsDeployed.Unknown(r, "Deploying cost tools")
				status.ConnectorStatusCostSynced.Unknown(r,
					"It takes about an hour to generate hour-level cost data")
			}
		}

		r.Status.SetSummary(status.WalkConnector(&r.Status))
		c.SetStatus(r.Status)
		rrs[i] = c
	}

	return rrs, nil
}

func ConnectorUpdate(mc model.ClientSet, input *model.Connector) (*model.ConnectorUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	// Predicated.
	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	if input.Category != types.ConnectorCategoryKubernetes && input.EnableFinOps {
		return nil, errors.New("invalid connector: finOps not supported")
	}

	// Conditional.
	c := mc.Connectors().UpdateOne(input).
		SetDescription(input.Description).
		SetEnableFinOps(input.EnableFinOps)

	if input.Status.Changed() {
		c.SetStatus(input.Status)
	}

	if !input.FinOpsCustomPricing.IsZero() {
		c.SetFinOpsCustomPricing(input.FinOpsCustomPricing)
	}

	if input.Name != "" {
		c.SetName(input.Name)
	}

	if input.Labels != nil {
		c.SetLabels(input.Labels)
	}

	if input.ConfigVersion != "" {
		c.SetConfigVersion(input.ConfigVersion)
	}

	if input.ConfigData != nil {
		c.SetConfigData(input.ConfigData)
	}

	return c, nil
}
