package view

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/seal-io/seal/pkg/apis/cost/validation"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs.

type CreateRequest struct {
	*model.PerspectiveCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	// TODO(michelia): support time range format https://docs.huihoo.com/grafana/2.6/reference/timerange/index.html
	if r.StartTime == "" {
		return errors.New("invalid start time: blank")
	}
	if r.EndTime == "" {
		return errors.New("invalid end time: blank")
	}

	return validation.ValidateAllocationQueries(r.AllocationQueries)
}

type CreateResponse = *model.PerspectiveOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.PerspectiveUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.IsNaive() {
		return errors.New("invalid id: blank")
	}
	if err := validation.ValidateAllocationQueries(r.AllocationQueries); err != nil {
		return err
	}

	var existed, err = modelClient.Perspectives().Query().
		Where(perspective.ID(r.ID)).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get perspective")
	}

	if r.Name == existed.Name &&
		r.StartTime == existed.StartTime &&
		r.EndTime == existed.EndTime &&
		reflect.DeepEqual(r.AllocationQueries, existed.AllocationQueries) {
		return errors.New("invalid input: nothing update")
	}

	return nil
}

type GetRequest struct {
	*model.PerspectiveQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.IsNaive() {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = *model.PerspectiveOutput

// Batch APIs.

type CollectionDeleteRequest []*model.PerspectiveQueryInput

func (r CollectionDeleteRequest) Validate() error {
	if len(r) == 0 {
		return errors.New("invalid input: empty")
	}
	for _, i := range r {
		if !i.ID.Valid(0) {
			return errors.New("invalid id: blank")
		}
	}
	return nil
}

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Perspective, perspective.OrderOption] `query:",inline"`

	Name string `query:"name,omitempty"`
}

type CollectionGetResponse = []*model.PerspectiveOutput

// Extensional APIs.

type FieldType string

const (
	FieldTypeGroupBy = "groupBy"
	FieldTypeFilter  = "filter"
	FieldTypeStep    = "step"
)

type CollectionRouteFieldsRequest struct {
	_ struct{} `route:"GET=/fields"`

	StartTime *time.Time `query:"startTime"`
	EndTime   *time.Time `query:"endTime"`
	FieldType FieldType  `query:"fieldType"`
}

func (r *CollectionRouteFieldsRequest) Validate() error {
	if r.FieldType == "" {
		return errors.New("invalid field type: blank")
	}
	if r.FieldType != FieldTypeFilter && r.FieldType != FieldTypeGroupBy && r.FieldType != FieldTypeStep {
		return errors.New("invalid field type: not support")
	}
	return nil
}

type CollectionRouteFieldsResponse = []PerspectiveField

type CollectionRouteFieldValuesRequest struct {
	_ struct{} `route:"GET=/field-values"`

	StartTime *time.Time        `query:"startTime"`
	EndTime   *time.Time        `query:"endTime"`
	FieldName types.FilterField `query:"fieldName"`
	FieldType FieldType         `query:"fieldType"`
}

func (r *CollectionRouteFieldValuesRequest) Validate() error {
	if r.StartTime != nil && r.EndTime != nil && r.EndTime.Before(*r.StartTime) {
		return errors.New("invalid time range: end time is early than start time")
	}
	if r.FieldName == "" {
		return errors.New("invalid field name: blank")
	}
	if r.FieldType == "" {
		return errors.New("invalid field type: blank")
	}
	if r.FieldType != FieldTypeFilter {
		return errors.New("invalid field type: not support")
	}
	if !strings.HasPrefix(string(r.FieldName), types.LabelPrefix) {
		for _, v := range BuiltInPerspectiveFilterFields {
			if v.FieldName == string(r.FieldName) {
				return nil
			}
		}
		return errors.New("invalid field name: unsupported")
	}
	return nil
}

type CollectionRouteValuesResponse = []PerspectiveValue
