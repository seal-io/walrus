// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// SettingQueryInput is the input for the Setting query.
type SettingQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the SettingQueryInput to Setting.
func (in SettingQueryInput) Model() *Setting {
	return &Setting{
		ID: in.ID,
	}
}

// SettingCreateInput is the input for the Setting creation.
type SettingCreateInput struct {
	// The name of system setting.
	Name string `json:"name"`
	// The value of system setting, store in string.
	Value string `json:"value,omitempty"`
	// Indicate the system setting should be hidden or not, default is visible.
	Hidden *bool `json:"hidden,omitempty"`
	// Indicate the system setting should be edited or not, default is readonly.
	Editable *bool `json:"editable,omitempty"`
	// Indicate the system setting should be exposed or not, default is exposed.
	Private *bool `json:"private,omitempty"`
}

// Model converts the SettingCreateInput to Setting.
func (in SettingCreateInput) Model() *Setting {
	var entity = &Setting{
		Name:     in.Name,
		Value:    in.Value,
		Hidden:   in.Hidden,
		Editable: in.Editable,
		Private:  in.Private,
	}
	return entity
}

// SettingUpdateInput is the input for the Setting modification.
type SettingUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// The name of system setting.
	Name string `json:"name,omitempty"`
	// The value of system setting, store in string.
	Value string `json:"value,omitempty"`
	// Indicate the system setting should be hidden or not, default is visible.
	Hidden *bool `json:"hidden,omitempty"`
	// Indicate the system setting should be edited or not, default is readonly.
	Editable *bool `json:"editable,omitempty"`
	// Indicate the system setting should be exposed or not, default is exposed.
	Private *bool `json:"private,omitempty"`
}

// Model converts the SettingUpdateInput to Setting.
func (in SettingUpdateInput) Model() *Setting {
	var entity = &Setting{
		ID:       in.ID,
		Name:     in.Name,
		Value:    in.Value,
		Hidden:   in.Hidden,
		Editable: in.Editable,
		Private:  in.Private,
	}
	return entity
}

// SettingOutput is the output for the Setting.
type SettingOutput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The name of system setting.
	Name string `json:"name,omitempty"`
	// The value of system setting, store in string.
	Value string `json:"value,omitempty"`
	// Indicate the system setting should be hidden or not, default is visible.
	Hidden *bool `json:"hidden,omitempty"`
	// Indicate the system setting should be edited or not, default is readonly.
	Editable *bool `json:"editable,omitempty"`
	// Indicate the system setting should be exposed or not, default is exposed.
	Private *bool `json:"private,omitempty"`
}

// ExposeSetting converts the Setting to SettingOutput.
func ExposeSetting(in *Setting) *SettingOutput {
	if in == nil {
		return nil
	}
	var entity = &SettingOutput{
		ID:         in.ID,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
		Name:       in.Name,
		Value:      in.Value,
		Hidden:     in.Hidden,
		Editable:   in.Editable,
		Private:    in.Private,
	}
	return entity
}

// ExposeSettings converts the Setting slice to SettingOutput pointer slice.
func ExposeSettings(in []*Setting) []*SettingOutput {
	var out = make([]*SettingOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeSetting(in[i])
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
