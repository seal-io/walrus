package resourcedefinition

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.ResourceDefinitionCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.ResourceDefinitionOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.ResourceDefinitionCreateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return err
	}

	if r.UiSchema != nil {
		if err := r.UiSchema.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type (
	GetRequest = model.ResourceDefinitionQueryInput

	GetResponse = *model.ResourceDefinitionOutput
)

type UpdateRequest struct {
	model.ResourceDefinitionUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.ResourceDefinitionUpdateInput.Validate(); err != nil {
		return err
	}

	if r.UiSchema != nil {
		if err := r.UiSchema.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type DeleteRequest = model.ResourceDefinitionDeleteInput

type (
	CollectionGetRequest struct {
		model.ResourceDefinitionQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.ResourceDefinition, resourcedefinition.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ResourceDefinitionOutput
)

func (r *CollectionGetRequest) Validate() error {
	if err := r.ResourceDefinitionQueryInputs.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest = model.ResourceDefinitionDeleteInputs
