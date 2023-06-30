// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// VariableQueryInput is the input for the Variable query.
type VariableQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the VariableQueryInput to Variable.
func (in VariableQueryInput) Model() *Variable {
	return &Variable{
		ID: in.ID,
	}
}

// VariableCreateInput is the input for the Variable creation.
type VariableCreateInput struct {
	// The name of variable.
	Name string `json:"name"`
	// The value of variable, store in string.
	Value crypto.String `json:"value"`
	// The value is sensitive or not.
	Sensitive bool `json:"sensitive,omitempty"`
	// Description of the variable.
	Description string `json:"description,omitempty"`
	// Project to which the variable belongs.
	Project *ProjectQueryInput `json:"project,omitempty"`
	// Environment to which the variable belongs.
	Environment *EnvironmentQueryInput `json:"environment,omitempty"`
}

// Model converts the VariableCreateInput to Variable.
func (in VariableCreateInput) Model() *Variable {
	var entity = &Variable{
		Name:        in.Name,
		Value:       in.Value,
		Sensitive:   in.Sensitive,
		Description: in.Description,
	}
	if in.Project != nil {
		entity.ProjectID = in.Project.ID
	}
	if in.Environment != nil {
		entity.EnvironmentID = in.Environment.ID
	}
	return entity
}

// VariableUpdateInput is the input for the Variable modification.
type VariableUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// The value of variable, store in string.
	Value crypto.String `json:"value,omitempty"`
	// The value is sensitive or not.
	Sensitive bool `json:"sensitive,omitempty"`
	// Description of the variable.
	Description string `json:"description,omitempty"`
}

// Model converts the VariableUpdateInput to Variable.
func (in VariableUpdateInput) Model() *Variable {
	var entity = &Variable{
		ID:          in.ID,
		Value:       in.Value,
		Sensitive:   in.Sensitive,
		Description: in.Description,
	}
	return entity
}

// VariableOutput is the output for the Variable.
type VariableOutput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `json:"id,omitempty"`
	// CreateTime holds the value of the "createTime" field.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// UpdateTime holds the value of the "updateTime" field.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The name of variable.
	Name string `json:"name,omitempty"`
	// The value of variable, store in string.
	Value crypto.String `json:"value,omitempty"`
	// The value is sensitive or not.
	Sensitive bool `json:"sensitive,omitempty"`
	// Description of the variable.
	Description string `json:"description,omitempty"`
	// Project to which the variable belongs.
	Project *ProjectOutput `json:"project,omitempty"`
	// Environment to which the variable belongs.
	Environment *EnvironmentOutput `json:"environment,omitempty"`
}

// ExposeVariable converts the Variable to VariableOutput.
func ExposeVariable(in *Variable) *VariableOutput {
	if in == nil {
		return nil
	}
	var entity = &VariableOutput{
		ID:          in.ID,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
		Name:        in.Name,
		Value:       in.Value,
		Sensitive:   in.Sensitive,
		Description: in.Description,
		Project:     ExposeProject(in.Edges.Project),
		Environment: ExposeEnvironment(in.Edges.Environment),
	}
	if in.ProjectID != "" {
		if entity.Project == nil {
			entity.Project = &ProjectOutput{}
		}
		entity.Project.ID = in.ProjectID
	}
	if in.EnvironmentID != "" {
		if entity.Environment == nil {
			entity.Environment = &EnvironmentOutput{}
		}
		entity.Environment.ID = in.EnvironmentID
	}
	return entity
}

// ExposeVariables converts the Variable slice to VariableOutput pointer slice.
func ExposeVariables(in []*Variable) []*VariableOutput {
	var out = make([]*VariableOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeVariable(in[i])
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
