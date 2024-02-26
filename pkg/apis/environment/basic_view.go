package environment

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.EnvironmentCreateInput `path:",inline" json:",inline"`

		Preview bool `json:"preview,default=false"`
	}

	CreateResponse = *model.EnvironmentOutput
)

func (r *CreateRequest) Validate() error {
	return validateEnvironmentCreateInput(r.EnvironmentCreateInput)
}

type (
	GetRequest struct {
		model.EnvironmentQueryInput `path:",inline"`

		IncludeSummary bool `query:"includeSummary,omitempty"`
	}

	GetResponse = *environmentOutput
)

type UpdateRequest struct {
	model.EnvironmentUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.EnvironmentQueryInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsValidName(r.Name); err != nil {
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

		IncludeSummary bool `query:"includeSummary,omitempty"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*environmentOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest = model.EnvironmentDeleteInputs

type environmentOutput struct {
	model.EnvironmentOutput `json:",inline"`
	StatusSummary           status.Count `json:"statusSummary"`
}

func exposeEnvironment(entity *model.Environment) *environmentOutput {
	output := &environmentOutput{
		EnvironmentOutput: *model.ExposeEnvironment(entity),
	}

	if len(entity.Edges.Resources) > 0 {
		for _, v := range entity.Edges.Resources {
			switch {
			case v.Status.Error:
				output.StatusSummary.Error++
			case v.Status.Transitioning:
				output.StatusSummary.Transitioning++
			case v.Status.Inactive:
				output.StatusSummary.Inactive++
			default:
				output.StatusSummary.Ready++
			}
		}
	}

	return output
}

func exposeEnvironments(entities []*model.Environment) []*environmentOutput {
	output := make([]*environmentOutput, len(entities))
	for i, v := range entities {
		output[i] = exposeEnvironment(v)
	}

	return output
}
