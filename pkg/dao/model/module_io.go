// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
)

// ModuleQueryInput is the input for the Module query.
type ModuleQueryInput struct {
	// It is also the name of the module.
	ID string `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the ModuleQueryInput to Module.
func (in ModuleQueryInput) Model() *Module {
	return &Module{
		ID: in.ID,
	}
}

// ModuleCreateInput is the input for the Module creation.
type ModuleCreateInput struct {
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Description of the module.
	Description string `json:"description,omitempty"`
	// Labels of the module.
	Labels map[string]string `json:"labels,omitempty"`
	// Source of the module.
	Source string `json:"source"`
	// Version of the module.
	Version string `json:"version,omitempty"`
	// Schema of the module.
	Schema *types.ModuleSchema `json:"schema,omitempty"`
}

// Model converts the ModuleCreateInput to Module.
func (in ModuleCreateInput) Model() *Module {
	var entity = &Module{
		Status:        in.Status,
		StatusMessage: in.StatusMessage,
		Description:   in.Description,
		Labels:        in.Labels,
		Source:        in.Source,
		Version:       in.Version,
		Schema:        in.Schema,
	}
	return entity
}

// ModuleUpdateInput is the input for the Module modification.
type ModuleUpdateInput struct {
	// It is also the name of the module.
	ID string `uri:"id" json:"-"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Description of the module.
	Description string `json:"description,omitempty"`
	// Labels of the module.
	Labels map[string]string `json:"labels,omitempty"`
	// Source of the module.
	Source string `json:"source,omitempty"`
	// Version of the module.
	Version string `json:"version,omitempty"`
	// Schema of the module.
	Schema *types.ModuleSchema `json:"schema,omitempty"`
}

// Model converts the ModuleUpdateInput to Module.
func (in ModuleUpdateInput) Model() *Module {
	var entity = &Module{
		ID:            in.ID,
		Status:        in.Status,
		StatusMessage: in.StatusMessage,
		Description:   in.Description,
		Labels:        in.Labels,
		Source:        in.Source,
		Version:       in.Version,
		Schema:        in.Schema,
	}
	return entity
}

// ModuleOutput is the output for the Module.
type ModuleOutput struct {
	// It is also the name of the module.
	ID string `json:"id,omitempty"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Description of the module.
	Description string `json:"description,omitempty"`
	// Labels of the module.
	Labels map[string]string `json:"labels,omitempty"`
	// Source of the module.
	Source string `json:"source,omitempty"`
	// Version of the module.
	Version string `json:"version,omitempty"`
	// Schema of the module.
	Schema *types.ModuleSchema `json:"schema,omitempty"`
	// Applications holds the value of the applications edge.
	Applications []*ApplicationModuleRelationshipOutput `json:"applications,omitempty"`
}

// ExposeModule converts the Module to ModuleOutput.
func ExposeModule(in *Module) *ModuleOutput {
	if in == nil {
		return nil
	}
	var entity = &ModuleOutput{
		ID:            in.ID,
		Status:        in.Status,
		StatusMessage: in.StatusMessage,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
		Description:   in.Description,
		Labels:        in.Labels,
		Source:        in.Source,
		Version:       in.Version,
		Schema:        in.Schema,
		Applications:  ExposeApplicationModuleRelationships(in.Edges.Applications),
	}
	return entity
}

// ExposeModules converts the Module slice to ModuleOutput pointer slice.
func ExposeModules(in []*Module) []*ModuleOutput {
	var out = make([]*ModuleOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeModule(in[i])
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
