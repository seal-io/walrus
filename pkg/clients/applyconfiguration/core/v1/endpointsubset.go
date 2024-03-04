// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

// EndpointSubsetApplyConfiguration represents an declarative configuration of the EndpointSubset type for use
// with apply.
type EndpointSubsetApplyConfiguration struct {
	Addresses         []EndpointAddressApplyConfiguration `json:"addresses,omitempty"`
	NotReadyAddresses []EndpointAddressApplyConfiguration `json:"notReadyAddresses,omitempty"`
	Ports             []EndpointPortApplyConfiguration    `json:"ports,omitempty"`
}

// EndpointSubsetApplyConfiguration constructs an declarative configuration of the EndpointSubset type for use with
// apply.
func EndpointSubset() *EndpointSubsetApplyConfiguration {
	return &EndpointSubsetApplyConfiguration{}
}

// WithAddresses adds the given value to the Addresses field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Addresses field.
func (b *EndpointSubsetApplyConfiguration) WithAddresses(values ...*EndpointAddressApplyConfiguration) *EndpointSubsetApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithAddresses")
		}
		b.Addresses = append(b.Addresses, *values[i])
	}
	return b
}

// WithNotReadyAddresses adds the given value to the NotReadyAddresses field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NotReadyAddresses field.
func (b *EndpointSubsetApplyConfiguration) WithNotReadyAddresses(values ...*EndpointAddressApplyConfiguration) *EndpointSubsetApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNotReadyAddresses")
		}
		b.NotReadyAddresses = append(b.NotReadyAddresses, *values[i])
	}
	return b
}

// WithPorts adds the given value to the Ports field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Ports field.
func (b *EndpointSubsetApplyConfiguration) WithPorts(values ...*EndpointPortApplyConfiguration) *EndpointSubsetApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithPorts")
		}
		b.Ports = append(b.Ports, *values[i])
	}
	return b
}
