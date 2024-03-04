// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

// RollingUpdateDeploymentApplyConfiguration represents an declarative configuration of the RollingUpdateDeployment type for use
// with apply.
type RollingUpdateDeploymentApplyConfiguration struct {
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
	MaxSurge       *intstr.IntOrString `json:"maxSurge,omitempty"`
}

// RollingUpdateDeploymentApplyConfiguration constructs an declarative configuration of the RollingUpdateDeployment type for use with
// apply.
func RollingUpdateDeployment() *RollingUpdateDeploymentApplyConfiguration {
	return &RollingUpdateDeploymentApplyConfiguration{}
}

// WithMaxUnavailable sets the MaxUnavailable field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MaxUnavailable field is set to the value of the last call.
func (b *RollingUpdateDeploymentApplyConfiguration) WithMaxUnavailable(value intstr.IntOrString) *RollingUpdateDeploymentApplyConfiguration {
	b.MaxUnavailable = &value
	return b
}

// WithMaxSurge sets the MaxSurge field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MaxSurge field is set to the value of the last call.
func (b *RollingUpdateDeploymentApplyConfiguration) WithMaxSurge(value intstr.IntOrString) *RollingUpdateDeploymentApplyConfiguration {
	b.MaxSurge = &value
	return b
}
