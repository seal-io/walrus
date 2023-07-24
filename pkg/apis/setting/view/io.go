package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// Basic APIs.

type UpdateRequest struct {
	ID    object.ID     `uri:"id,omitempty" json:"id,omitempty"`
	Value crypto.String `json:"value"`

	Name string `json:"-"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if r.ID == "" {
		return errors.New("invalid id: blank")
	}

	var p predicate.Setting
	if r.ID.Valid() {
		p = setting.ID(r.ID)
	} else {
		p = setting.Name(string(r.ID))
	}

	// Only allow updating publicly editable setting.
	entity, err := modelClient.Settings().Query().
		Where(
			p,
			setting.Private(false),
			setting.Editable(true)).
		Select(setting.FieldName).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get setting")
	}

	// Get setting name by id.
	r.Name = entity.Name

	return nil
}

type GetRequest struct {
	model.SettingQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if r.ID == "" {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse struct {
	*model.SettingOutput `json:",inline"`

	// Configured indicates the setting whether to be configured.
	Configured bool `json:"configured"`
}

// Batch APIs.

type CollectionUpdateRequest []*UpdateRequest

func (r CollectionUpdateRequest) ValidateWith(ctx context.Context, input any) error {
	if len(r) == 0 {
		return errors.New("invalid input: empty list")
	}

	for _, i := range r {
		if i == nil {
			return errors.New("invalid input: empty item")
		}

		err := i.ValidateWith(ctx, input)
		if err != nil {
			return err
		}
	}

	return nil
}

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`

	IDs []object.ID `query:"id,omitempty"`
}

func (r *CollectionGetRequest) Validate() error {
	for i := range r.IDs {
		if r.IDs[i] == "" {
			return errors.New("invalid id: blank")
		}
	}

	return nil
}

type CollectionGetResponse = []*GetResponse

// Extensional APIs.
