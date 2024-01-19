package project

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.ProjectCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.ProjectOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.ProjectCreateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsValidName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

type (
	GetRequest struct {
		model.ProjectQueryInput `path:",inline"`
	}

	GetResponse = *model.ProjectOutput
)

type UpdateRequest = model.ProjectUpdateInput

type DeleteRequest = model.ProjectDeleteInput

type (
	CollectionGetRequest struct {
		model.ProjectQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Project, project.OrderOption,
		] `query:",inline"`
	}

	CollectionGetResponse = []*model.ProjectOutput
)

type CollectionDeleteRequest = model.ProjectDeleteInputs
