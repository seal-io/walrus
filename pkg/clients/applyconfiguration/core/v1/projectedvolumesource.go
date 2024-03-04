// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

// ProjectedVolumeSourceApplyConfiguration represents an declarative configuration of the ProjectedVolumeSource type for use
// with apply.
type ProjectedVolumeSourceApplyConfiguration struct {
	Sources     []VolumeProjectionApplyConfiguration `json:"sources,omitempty"`
	DefaultMode *int32                               `json:"defaultMode,omitempty"`
}

// ProjectedVolumeSourceApplyConfiguration constructs an declarative configuration of the ProjectedVolumeSource type for use with
// apply.
func ProjectedVolumeSource() *ProjectedVolumeSourceApplyConfiguration {
	return &ProjectedVolumeSourceApplyConfiguration{}
}

// WithSources adds the given value to the Sources field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Sources field.
func (b *ProjectedVolumeSourceApplyConfiguration) WithSources(values ...*VolumeProjectionApplyConfiguration) *ProjectedVolumeSourceApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithSources")
		}
		b.Sources = append(b.Sources, *values[i])
	}
	return b
}

// WithDefaultMode sets the DefaultMode field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DefaultMode field is set to the value of the last call.
func (b *ProjectedVolumeSourceApplyConfiguration) WithDefaultMode(value int32) *ProjectedVolumeSourceApplyConfiguration {
	b.DefaultMode = &value
	return b
}
