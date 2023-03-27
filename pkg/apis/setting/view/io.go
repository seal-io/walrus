package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type UpdateRequest struct {
	ID    types.ID `uri:"id"`
	Value *string  `json:"value"`

	Name string `json:"-"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(1) {
		return errors.New("invalid id: blank")
	}
	if r.Value == nil {
		return errors.New("invalid input: nil value")
	}

	var confirmSetting = []predicate.Setting{
		setting.Private(false),
		setting.Editable(true),
	}
	switch {
	case r.ID.IsNaive():
		confirmSetting = append(confirmSetting, setting.ID(r.ID))
	default:
		var keys = r.ID.Split()
		confirmSetting = append(confirmSetting, setting.Name(keys[0]))
	}
	var settingEntity, err = modelClient.Settings().Query().
		Where(confirmSetting...).
		Select(setting.FieldName, setting.FieldValue).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get setting")
	}
	r.Name = settingEntity.Name

	return nil
}

type GetRequest struct {
	*model.SettingQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(1) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = *model.SettingOutput

// Batch APIs

type CollectionUpdateRequest []*UpdateRequest

func (r CollectionUpdateRequest) ValidateWith(ctx context.Context, input any) error {
	if len(r) == 0 {
		return errors.New("invalid input: empty list")
	}
	for _, i := range r {
		var err = i.ValidateWith(ctx, input)
		if err != nil {
			return err
		}
	}
	return nil
}

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`

	IDs []types.ID `query:"id,omitempty"`
}

func (r *CollectionGetRequest) Validate() error {
	for i := range r.IDs {
		if !r.IDs[i].Valid(1) {
			return errors.New("invalid id: blank")
		}
	}
	return nil
}

type CollectionGetResponse = []*model.SettingOutput

// Extensional APIs
