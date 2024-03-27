// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	v1 "github.com/seal-io/walrus/pkg/apis/walrus/v1"
)

// EnvironmentSpecApplyConfiguration represents an declarative configuration of the EnvironmentSpec type for use
// with apply.
type EnvironmentSpecApplyConfiguration struct {
	Type        *v1.EnvironmentType `json:"type,omitempty"`
	DisplayName *string             `json:"displayName,omitempty"`
	Description *string             `json:"description,omitempty"`
}

// EnvironmentSpecApplyConfiguration constructs an declarative configuration of the EnvironmentSpec type for use with
// apply.
func EnvironmentSpec() *EnvironmentSpecApplyConfiguration {
	return &EnvironmentSpecApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *EnvironmentSpecApplyConfiguration) WithType(value v1.EnvironmentType) *EnvironmentSpecApplyConfiguration {
	b.Type = &value
	return b
}

// WithDisplayName sets the DisplayName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DisplayName field is set to the value of the last call.
func (b *EnvironmentSpecApplyConfiguration) WithDisplayName(value string) *EnvironmentSpecApplyConfiguration {
	b.DisplayName = &value
	return b
}

// WithDescription sets the Description field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Description field is set to the value of the last call.
func (b *EnvironmentSpecApplyConfiguration) WithDescription(value string) *EnvironmentSpecApplyConfiguration {
	b.Description = &value
	return b
}