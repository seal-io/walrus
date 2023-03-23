// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// RoleQueryInput is the input for the Role query.
type RoleQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the RoleQueryInput to Role.
func (in RoleQueryInput) Model() *Role {
	return &Role{
		ID: in.ID,
	}
}

// RoleCreateInput is the input for the Role creation.
type RoleCreateInput struct {
	// The domain of the role.
	Domain string `json:"domain,omitempty"`
	// The name of the role.
	Name string `json:"name"`
	// The detail of the role.
	Description string `json:"description,omitempty"`
	// The policy list of the role.
	Policies types.RolePolicies `json:"policies,omitempty"`
	// Indicate whether the subject is builtin, decide when creating.
	Builtin bool `json:"builtin,omitempty"`
	// Indicate whether the subject is session level, decide when creating.
	Session bool `json:"session,omitempty"`
}

// Model converts the RoleCreateInput to Role.
func (in RoleCreateInput) Model() *Role {
	var entity = &Role{
		Domain:      in.Domain,
		Name:        in.Name,
		Description: in.Description,
		Policies:    in.Policies,
		Builtin:     in.Builtin,
		Session:     in.Session,
	}
	return entity
}

// RoleUpdateInput is the input for the Role modification.
type RoleUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// The detail of the role.
	Description string `json:"description,omitempty"`
	// The policy list of the role.
	Policies types.RolePolicies `json:"policies,omitempty"`
}

// Model converts the RoleUpdateInput to Role.
func (in RoleUpdateInput) Model() *Role {
	var entity = &Role{
		ID:          in.ID,
		Description: in.Description,
		Policies:    in.Policies,
	}
	return entity
}

// RoleOutput is the output for the Role.
type RoleOutput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The domain of the role.
	Domain string `json:"domain,omitempty"`
	// The name of the role.
	Name string `json:"name,omitempty"`
	// The detail of the role.
	Description string `json:"description,omitempty"`
	// The policy list of the role.
	Policies types.RolePolicies `json:"policies,omitempty"`
	// Indicate whether the subject is builtin, decide when creating.
	Builtin bool `json:"builtin,omitempty"`
	// Indicate whether the subject is session level, decide when creating.
	Session bool `json:"session,omitempty"`
}

// ExposeRole converts the Role to RoleOutput.
func ExposeRole(in *Role) *RoleOutput {
	if in == nil {
		return nil
	}
	var entity = &RoleOutput{
		ID:          in.ID,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
		Domain:      in.Domain,
		Name:        in.Name,
		Description: in.Description,
		Policies:    in.Policies,
		Builtin:     in.Builtin,
		Session:     in.Session,
	}
	return entity
}

// ExposeRoles converts the Role slice to RoleOutput pointer slice.
func ExposeRoles(in []*Role) []*RoleOutput {
	var out = make([]*RoleOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeRole(in[i])
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
