package perspective

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/seal-io/seal/pkg/apis/cost/validation"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	utilvalidation "github.com/seal-io/seal/utils/validation"
)

type (
	CreateRequest struct {
		model.PerspectiveCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.PerspectiveOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.PerspectiveCreateInput.Validate(); err != nil {
		return err
	}

	if err := utilvalidation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// TODO(michelia): support time range format https://docs.huihoo.com/grafana/2.6/reference/timerange/index.html
	if r.StartTime == "" {
		return errors.New("invalid start time: blank")
	}

	if r.EndTime == "" {
		return errors.New("invalid end time: blank")
	}

	return validation.ValidateCostQueries(r.CostQueries)
}

type (
	GetRequest struct {
		model.PerspectiveQueryInput `path:",inline"`
	}

	GetResponse = *model.PerspectiveOutput
)

type UpdateRequest struct {
	model.PerspectiveUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.PerspectiveUpdateInput.Validate(); err != nil {
		return err
	}

	if err := utilvalidation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := validation.ValidateCostQueries(r.CostQueries); err != nil {
		return err
	}

	entity, err := r.Client.Perspectives().Query().
		Where(perspective.ID(r.ID)).
		Only(r.Context)
	if err != nil {
		return runtime.Errorw(err, "failed to get perspective")
	}

	if r.Name == entity.Name &&
		r.StartTime == entity.StartTime &&
		r.EndTime == entity.EndTime &&
		reflect.DeepEqual(r.CostQueries, entity.CostQueries) {
		return errors.New("invalid input: nothing update")
	}

	return nil
}

type DeleteRequest = model.PerspectiveDeleteInput

type (
	CollectionGetRequest struct {
		model.PerspectiveQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Perspective, perspective.OrderOption,
		] `query:",inline"`

		Name string `query:"name,omitempty"`
	}

	CollectionGetResponse = []*model.PerspectiveOutput
)

type CollectionDeleteRequest = model.PerspectiveDeleteInputs
