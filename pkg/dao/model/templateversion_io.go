// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// TemplateVersionQueryInput is the input for the TemplateVersion query.
type TemplateVersionQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the TemplateVersionQueryInput to TemplateVersion.
func (in TemplateVersionQueryInput) Model() *TemplateVersion {
	return &TemplateVersion{
		ID: in.ID,
	}
}

// TemplateVersionCreateInput is the input for the TemplateVersion creation.
type TemplateVersionCreateInput struct {
	// Template version.
	Version string `json:"version"`
	// Template version source.
	Source string `json:"source"`
	// Schema of the template.
	Schema *types.TemplateSchema `json:"schema,omitempty"`
	// Template holds the value of the template edge.
	Template TemplateQueryInput `json:"template"`
}

// Model converts the TemplateVersionCreateInput to TemplateVersion.
func (in TemplateVersionCreateInput) Model() *TemplateVersion {
	var entity = &TemplateVersion{
		Version: in.Version,
		Source:  in.Source,
		Schema:  in.Schema,
	}
	entity.TemplateID = in.Template.ID
	return entity
}

// TemplateVersionUpdateInput is the input for the TemplateVersion modification.
type TemplateVersionUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// Schema of the template.
	Schema *types.TemplateSchema `json:"schema,omitempty"`
}

// Model converts the TemplateVersionUpdateInput to TemplateVersion.
func (in TemplateVersionUpdateInput) Model() *TemplateVersion {
	var entity = &TemplateVersion{
		ID:     in.ID,
		Schema: in.Schema,
	}
	return entity
}

// TemplateVersionOutput is the output for the TemplateVersion.
type TemplateVersionOutput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Template version.
	Version string `json:"version,omitempty"`
	// Template version source.
	Source string `json:"source,omitempty"`
	// Schema of the template.
	Schema *types.TemplateSchema `json:"schema,omitempty"`
	// Template holds the value of the template edge.
	Template *TemplateOutput `json:"template,omitempty"`
}

// ExposeTemplateVersion converts the TemplateVersion to TemplateVersionOutput.
func ExposeTemplateVersion(in *TemplateVersion) *TemplateVersionOutput {
	if in == nil {
		return nil
	}
	var entity = &TemplateVersionOutput{
		ID:         in.ID,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
		Version:    in.Version,
		Source:     in.Source,
		Schema:     in.Schema,
		Template:   ExposeTemplate(in.Edges.Template),
	}
	if in.TemplateID != "" {
		if entity.Template == nil {
			entity.Template = &TemplateOutput{}
		}
		entity.Template.ID = in.TemplateID
	}
	return entity
}

// ExposeTemplateVersions converts the TemplateVersion slice to TemplateVersionOutput pointer slice.
func ExposeTemplateVersions(in []*TemplateVersion) []*TemplateVersionOutput {
	var out = make([]*TemplateVersionOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeTemplateVersion(in[i])
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
