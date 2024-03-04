// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

// ScaleStatusApplyConfiguration represents an declarative configuration of the ScaleStatus type for use
// with apply.
type ScaleStatusApplyConfiguration struct {
	Replicas *int32  `json:"replicas,omitempty"`
	Selector *string `json:"selector,omitempty"`
}

// ScaleStatusApplyConfiguration constructs an declarative configuration of the ScaleStatus type for use with
// apply.
func ScaleStatus() *ScaleStatusApplyConfiguration {
	return &ScaleStatusApplyConfiguration{}
}

// WithReplicas sets the Replicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Replicas field is set to the value of the last call.
func (b *ScaleStatusApplyConfiguration) WithReplicas(value int32) *ScaleStatusApplyConfiguration {
	b.Replicas = &value
	return b
}

// WithSelector sets the Selector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Selector field is set to the value of the last call.
func (b *ScaleStatusApplyConfiguration) WithSelector(value string) *ScaleStatusApplyConfiguration {
	b.Selector = &value
	return b
}
