package perspective

import (
	"errors"
	"strings"
	"time"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

type FieldType string

const (
	FieldTypeGroupBy = "groupBy"
	FieldTypeFilter  = "filter"
	FieldTypeStep    = "step"
)

type (
	CollectionRouteGetFieldsRequest struct {
		_ struct{} `route:"GET=/fields"`

		model.PerspectiveQueryInputs `path:",inline" query:",inline"`

		StartTime *time.Time `query:"startTime"`
		EndTime   *time.Time `query:"endTime"`
		FieldType FieldType  `query:"fieldType"`
	}

	CollectionRouteGetFieldsResponse = []Field
)

func (r *CollectionRouteGetFieldsRequest) Validate() error {
	if err := r.PerspectiveQueryInputs.Validate(); err != nil {
		return err
	}

	if r.FieldType == "" {
		return errors.New("invalid field type: blank")
	}

	if r.FieldType != FieldTypeFilter && r.FieldType != FieldTypeGroupBy && r.FieldType != FieldTypeStep {
		return errors.New("invalid field type: not support")
	}

	return nil
}

type (
	CollectionRouteGetFieldValuesRequest struct {
		_ struct{} `route:"GET=/field-values"`

		model.PerspectiveQueryInputs `path:",inline" query:",inline"`

		StartTime *time.Time        `query:"startTime"`
		EndTime   *time.Time        `query:"endTime"`
		FieldName types.FilterField `query:"fieldName"`
		FieldType FieldType         `query:"fieldType"`
	}

	CollectionRouteGetFieldValuesResponse = []Value
)

func (r *CollectionRouteGetFieldValuesRequest) Validate() error {
	if err := r.PerspectiveQueryInputs.Validate(); err != nil {
		return err
	}

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
		for _, v := range builtinFilterFields {
			if v.FieldName == string(r.FieldName) {
				return nil
			}
		}

		return errors.New("invalid field name: unsupported")
	}

	return nil
}
