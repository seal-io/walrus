// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package model

import "time"

// TemplateQueryInput is the input for the Template query.
type TemplateQueryInput struct {
	// It is also the name of the template.
	ID string `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the TemplateQueryInput to Template.
func (in TemplateQueryInput) Model() *Template {
	return &Template{
		ID: in.ID,
	}
}

// TemplateCreateInput is the input for the Template creation.
type TemplateCreateInput struct {
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Description of the template.
	Description string `json:"description,omitempty"`
	// A URL to an SVG or PNG image to be used as an icon.
	Icon string `json:"icon,omitempty"`
	// Labels of the template.
	Labels map[string]string `json:"labels,omitempty"`
	// Source of the template.
	Source string `json:"source"`
}

// Model converts the TemplateCreateInput to Template.
func (in TemplateCreateInput) Model() *Template {
	var entity = &Template{
		StatusMessage: in.StatusMessage,
		Description:   in.Description,
		Icon:          in.Icon,
		Labels:        in.Labels,
		Source:        in.Source,
	}
	return entity
}

// TemplateUpdateInput is the input for the Template modification.
type TemplateUpdateInput struct {
	// It is also the name of the template.
	ID string `uri:"id" json:"-"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Description of the template.
	Description string `json:"description,omitempty"`
	// A URL to an SVG or PNG image to be used as an icon.
	Icon string `json:"icon,omitempty"`
	// Labels of the template.
	Labels map[string]string `json:"labels,omitempty"`
	// Source of the template.
	Source string `json:"source,omitempty"`
}

// Model converts the TemplateUpdateInput to Template.
func (in TemplateUpdateInput) Model() *Template {
	var entity = &Template{
		ID:            in.ID,
		StatusMessage: in.StatusMessage,
		Description:   in.Description,
		Icon:          in.Icon,
		Labels:        in.Labels,
		Source:        in.Source,
	}
	return entity
}

// TemplateOutput is the output for the Template.
type TemplateOutput struct {
	// It is also the name of the template.
	ID string `json:"id,omitempty"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Description of the template.
	Description string `json:"description,omitempty"`
	// A URL to an SVG or PNG image to be used as an icon.
	Icon string `json:"icon,omitempty"`
	// Labels of the template.
	Labels map[string]string `json:"labels,omitempty"`
	// Source of the template.
	Source string `json:"source,omitempty"`
	// versions of the template.
	Versions []*TemplateVersionOutput `json:"versions,omitempty"`
}

// ExposeTemplate converts the Template to TemplateOutput.
func ExposeTemplate(in *Template) *TemplateOutput {
	if in == nil {
		return nil
	}
	var entity = &TemplateOutput{
		ID:            in.ID,
		Status:        in.Status,
		StatusMessage: in.StatusMessage,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
		Description:   in.Description,
		Icon:          in.Icon,
		Labels:        in.Labels,
		Source:        in.Source,
		Versions:      ExposeTemplateVersions(in.Edges.Versions),
	}
	return entity
}

// ExposeTemplates converts the Template slice to TemplateOutput pointer slice.
func ExposeTemplates(in []*Template) []*TemplateOutput {
	var out = make([]*TemplateOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeTemplate(in[i])
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
