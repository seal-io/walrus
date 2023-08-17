package setting

import (
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/setting"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type (
	GetRequest struct {
		model.SettingQueryInput `path:",inline"`
	}

	GetResponse = *model.SettingOutput
)

type UpdateRequest struct {
	model.SettingUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.SettingUpdateInput.Validate(); err != nil {
		return err
	}

	// Only allow updating publicly editable setting.
	cnt, err := r.Client.Settings().Query().
		Where(
			setting.ID(r.ID),
			setting.Private(false),
			setting.Editable(true)).
		Count(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get setting: %w", err)
	}

	if cnt != 1 {
		return errors.New("failed to get setting")
	}

	return nil
}

type (
	CollectionGetRequest struct {
		model.SettingQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Setting, setting.OrderOption,
		] `query:",inline"`

		IDs   []object.ID `query:"id,omitempty"`
		Names []string    `query:"name,omitempty"`
	}

	CollectionGetResponse = []*model.SettingOutput
)

func (r *CollectionGetRequest) Validate() error {
	if err := r.SettingQueryInputs.Validate(); err != nil {
		return err
	}

	for i := range r.IDs {
		if !r.IDs[i].Valid() {
			return errors.New("invalid id: blank")
		}
	}

	for i := range r.Names {
		if r.Names[i] == "" {
			return errors.New("invalid name: blank")
		}
	}

	return nil
}

type CollectionUpdateRequest struct {
	model.SettingUpdateInputs `path:",inline" json:",inline"`
}

func (r *CollectionUpdateRequest) Validate() error {
	if err := r.SettingUpdateInputs.Validate(); err != nil {
		return err
	}

	// Only allow updating publicly editable setting.
	cnt, err := r.Client.Settings().Query().
		Where(
			setting.IDIn(r.IDs()...),
			setting.Private(false),
			setting.Editable(true)).
		Select(setting.FieldName).
		Count(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	if cnt != len(r.Items) {
		return errors.New("failed to get settings")
	}

	return nil
}

func exposeSetting(in *model.Setting) *model.SettingOutput {
	if in.Value != "" {
		in.Configured = true
	}

	// Erases the value of sensitive setting.
	if in.Sensitive != nil && *in.Sensitive {
		in.Value = ""
	}

	return model.ExposeSetting(in)
}

func exposeSettings(in []*model.Setting) []*model.SettingOutput {
	out := make([]*model.SettingOutput, 0, len(in))

	for i := 0; i < len(in); i++ {
		o := exposeSetting(in[i])
		if o == nil {
			continue
		}

		out = append(out, o)
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
