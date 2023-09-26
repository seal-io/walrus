package environment

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
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
	return validateEnvironmentCreateInput(r.EnvironmentCreateInput)
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
