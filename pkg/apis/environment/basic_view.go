package environment

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.EnvironmentCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.EnvironmentOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.EnvironmentCreateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Verify connections.
	connectorIDs := make([]object.ID, len(r.Connectors))
	for i, c := range r.Connectors {
		connectorIDs[i] = c.Connector.ID
	}

	if err := validateConnectors(r.Context, r.Client, connectorIDs); err != nil {
		return err
	}

	// Verify services.
	for i := range r.Services {
		if r.Services[i] == nil {
			return errors.New("empty service")
		}

		if err := validation.IsDNSLabel(r.Services[i].Name); err != nil {
			return fmt.Errorf("invalid service name: %w", err)
		}
	}

	// Get template versions.
	tvIDs := make([]object.ID, len(r.Services))
	for i := range r.Services {
		tvIDs[i] = r.Services[i].Template.ID
	}

	tvs, err := r.Client.TemplateVersions().Query().
		Where(templateversion.IDIn(tvIDs...)).
		Select(
			templateversion.FieldID,
			templateversion.FieldName,
			templateversion.FieldVersion,
			templateversion.FieldSchema).
		All(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get template version: %w", err)
	}

	// Map template version by ID for service validation.
	tvm := make(map[object.ID]*model.TemplateVersion, len(tvs))
	for i := range tvs {
		tvm[tvs[i].ID] = tvs[i]
	}

	// Verify service's variables with variables schema that defined on the template version.
	for _, svc := range r.Services {
		err = svc.Attributes.ValidateWith(tvm[svc.Template.ID].Schema.Variables)
		if err != nil {
			return fmt.Errorf("invalid variables: %w", err)
		}
	}

	return nil
}

type (
	GetRequest struct {
		model.EnvironmentQueryInput `path:",inline"`
	}

	GetResponse = *model.EnvironmentOutput
)

type UpdateRequest struct {
	model.EnvironmentUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.EnvironmentQueryInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	connectorIDs := make([]object.ID, len(r.Connectors))
	for i, c := range r.Connectors {
		connectorIDs[i] = c.Connector.ID
	}

	if err := validateConnectors(r.Context, r.Client, connectorIDs); err != nil {
		return err
	}

	return nil
}

type DeleteRequest = model.EnvironmentDeleteInput

type (
	CollectionGetRequest struct {
		model.EnvironmentQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Environment, environment.OrderOption,
		] `query:",inline"`
	}

	CollectionGetResponse = []*model.EnvironmentOutput
)

type CollectionDeleteRequest = model.EnvironmentDeleteInputs

// validateConnectors checks if given connector IDs are valid within the same project or globally.
func validateConnectors(ctx context.Context, mc model.ClientSet, ids []object.ID) error {
	if len(ids) == 0 {
		return nil
	}

	var typeCount []struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
	}

	err := mc.Connectors().Query().
		Where(connector.IDIn(ids...)).
		GroupBy(connector.FieldType).
		Aggregate(model.Count()).
		Scan(ctx, &typeCount)
	if err != nil {
		return fmt.Errorf("failed to get connector type count: %w", err)
	}

	// Validate connector type is duplicated,
	// only one connector type is allowed in one environment.
	for _, c := range typeCount {
		if c.Count > 1 {
			return fmt.Errorf("invalid connectors: duplicated connector type %s", c.Type)
		}
	}

	return nil
}
