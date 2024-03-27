// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

// EndpointHintsApplyConfiguration represents an declarative configuration of the EndpointHints type for use
// with apply.
type EndpointHintsApplyConfiguration struct {
	ForZones []ForZoneApplyConfiguration `json:"forZones,omitempty"`
}

// EndpointHintsApplyConfiguration constructs an declarative configuration of the EndpointHints type for use with
// apply.
func EndpointHints() *EndpointHintsApplyConfiguration {
	return &EndpointHintsApplyConfiguration{}
}

// WithForZones adds the given value to the ForZones field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the ForZones field.
func (b *EndpointHintsApplyConfiguration) WithForZones(values ...*ForZoneApplyConfiguration) *EndpointHintsApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithForZones")
		}
		b.ForZones = append(b.ForZones, *values[i])
	}
	return b
}