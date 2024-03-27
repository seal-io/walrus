// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

// VolumeDeviceApplyConfiguration represents an declarative configuration of the VolumeDevice type for use
// with apply.
type VolumeDeviceApplyConfiguration struct {
	Name       *string `json:"name,omitempty"`
	DevicePath *string `json:"devicePath,omitempty"`
}

// VolumeDeviceApplyConfiguration constructs an declarative configuration of the VolumeDevice type for use with
// apply.
func VolumeDevice() *VolumeDeviceApplyConfiguration {
	return &VolumeDeviceApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *VolumeDeviceApplyConfiguration) WithName(value string) *VolumeDeviceApplyConfiguration {
	b.Name = &value
	return b
}

// WithDevicePath sets the DevicePath field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DevicePath field is set to the value of the last call.
func (b *VolumeDeviceApplyConfiguration) WithDevicePath(value string) *VolumeDeviceApplyConfiguration {
	b.DevicePath = &value
	return b
}