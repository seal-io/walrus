package view

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type CreateRequest struct {
	*model.Perspective `json:",inline"`
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
	if len(r.AllocationQueries) == 0 {
		return errors.New("invalid allocation queries: blank")
	}

	return nil
}

type IDRequest struct {
	ID types.ID `uri:"id"`
}

func (r *IDRequest) Validate() error {
	if !r.ID.IsNaive() {
		return errors.New("invalid id: blank")
	}

	return nil
}

type UpdateRequest struct {
	*model.Perspective `json:",inline"`

	ID types.ID `uri:"id"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.IsNaive() {
		return errors.New("invalid id: blank")
	}

	var existed, err = modelClient.Perspectives().Query().
		Where(perspective.ID(r.ID)).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return runtime.Error(http.StatusBadRequest, "invalid perspective: not found")
		}
		return runtime.ErrorfP(http.StatusInternalServerError, "failed to get requesting perspective: %w", err)
	}

	if reflect.DeepEqual(r.AllocationQueries, existed.AllocationQueries) &&
		r.StartTime == existed.StartTime &&
		r.EndTime == existed.EndTime {
		return errors.New("invalid input: nothing update")
	}

	return nil
}

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	Name string `query:"name,omitempty"`
}

type CollectionGetResponse = []*model.Perspective

// Extensional APIs

type CollectionRouteFieldsRequest struct {
	_ struct{} `route:"GET=/fields"`

	StartTime *time.Time `query:"startTime"`
	EndTime   *time.Time `query:"endTime"`
}

type CollectionRouteFieldsResponse = []PerspectiveField

type CollectionRouteFieldValuesRequest struct {
	_ struct{} `route:"GET=/field-values"`

	StartTime *time.Time        `query:"startTime"`
	EndTime   *time.Time        `query:"endTime"`
	FieldName types.FilterField `query:"fieldName"`
}

func (r CollectionRouteFieldValuesRequest) Validate() error {
	if r.StartTime != nil && r.EndTime != nil && r.EndTime.Before(*r.StartTime) {
		return errors.New("invalid time range: end time is early than start time")
	}

	if r.FieldName == "" {
		return errors.New("invalid field name: blank")
	}

	if !strings.HasPrefix(string(r.FieldName), types.LabelPrefix) {
		for _, v := range BuiltInPerspectiveFields {
			if v.FieldName == r.FieldName {
				return nil
			}
		}
		return errors.New("invalid field name: unsupported")
	}
	return nil
}

type CollectionRouteValuesResponse = []PerspectiveValue
