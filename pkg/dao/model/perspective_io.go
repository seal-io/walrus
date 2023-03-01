// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
)

// PerspectiveQueryInput is the input for the Perspective query.
type PerspectiveQueryInput struct {
	// ID holds the value of the "id" field.
	ID types.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the PerspectiveQueryInput to Perspective.
func (in PerspectiveQueryInput) Model() *Perspective {
	return &Perspective{
		ID: in.ID,
	}
}

// PerspectiveCreateInput is the input for the Perspective creation.
type PerspectiveCreateInput struct {
	// Name for current perspective
	Name string `json:"name"`
	// Start time for current perspective
	StartTime string `json:"startTime"`
	// End time for current perspective
	EndTime string `json:"endTime"`
	// Is builtin Perspective
	Builtin bool `json:"builtin,omitempty"`
	// Indicated the perspective included allocation queries, record the used query condition
	AllocationQueries []types.QueryCondition `json:"allocationQueries,omitempty"`
}

// Model converts the PerspectiveCreateInput to Perspective.
func (in PerspectiveCreateInput) Model() *Perspective {
	var entity = &Perspective{
		Name:              in.Name,
		StartTime:         in.StartTime,
		EndTime:           in.EndTime,
		Builtin:           in.Builtin,
		AllocationQueries: in.AllocationQueries,
	}
	return entity
}

// PerspectiveUpdateInput is the input for the Perspective modification.
type PerspectiveUpdateInput struct {
	// ID holds the value of the "id" field.
	ID types.ID `uri:"id" json:"-"`
	// Start time for current perspective
	StartTime string `json:"startTime,omitempty"`
	// End time for current perspective
	EndTime string `json:"endTime,omitempty"`
	// Is builtin Perspective
	Builtin bool `json:"builtin,omitempty"`
	// Indicated the perspective included allocation queries, record the used query condition
	AllocationQueries []types.QueryCondition `json:"allocationQueries,omitempty"`
}

// Model converts the PerspectiveUpdateInput to Perspective.
func (in PerspectiveUpdateInput) Model() *Perspective {
	var entity = &Perspective{
		ID:                in.ID,
		StartTime:         in.StartTime,
		EndTime:           in.EndTime,
		Builtin:           in.Builtin,
		AllocationQueries: in.AllocationQueries,
	}
	return entity
}

// PerspectiveOutput is the output for the Perspective.
type PerspectiveOutput struct {
	// ID holds the value of the "id" field.
	ID types.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Name for current perspective
	Name string `json:"name,omitempty"`
	// Start time for current perspective
	StartTime string `json:"startTime,omitempty"`
	// End time for current perspective
	EndTime string `json:"endTime,omitempty"`
	// Is builtin Perspective
	Builtin bool `json:"builtin,omitempty"`
	// Indicated the perspective included allocation queries, record the used query condition
	AllocationQueries []types.QueryCondition `json:"allocationQueries,omitempty"`
}

// ExposePerspective converts the Perspective to PerspectiveOutput.
func ExposePerspective(in *Perspective) *PerspectiveOutput {
	if in == nil {
		return nil
	}
	var entity = &PerspectiveOutput{
		ID:                in.ID,
		CreateTime:        in.CreateTime,
		UpdateTime:        in.UpdateTime,
		Name:              in.Name,
		StartTime:         in.StartTime,
		EndTime:           in.EndTime,
		Builtin:           in.Builtin,
		AllocationQueries: in.AllocationQueries,
	}
	return entity
}

// ExposePerspectives converts the Perspective slice to PerspectiveOutput pointer slice.
func ExposePerspectives(in []*Perspective) []*PerspectiveOutput {
	var out = make([]*PerspectiveOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposePerspective(in[i])
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
