// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

// SessionAffinityConfigApplyConfiguration represents an declarative configuration of the SessionAffinityConfig type for use
// with apply.
type SessionAffinityConfigApplyConfiguration struct {
	ClientIP *ClientIPConfigApplyConfiguration `json:"clientIP,omitempty"`
}

// SessionAffinityConfigApplyConfiguration constructs an declarative configuration of the SessionAffinityConfig type for use with
// apply.
func SessionAffinityConfig() *SessionAffinityConfigApplyConfiguration {
	return &SessionAffinityConfigApplyConfiguration{}
}

// WithClientIP sets the ClientIP field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ClientIP field is set to the value of the last call.
func (b *SessionAffinityConfigApplyConfiguration) WithClientIP(value *ClientIPConfigApplyConfiguration) *SessionAffinityConfigApplyConfiguration {
	b.ClientIP = value
	return b
}
