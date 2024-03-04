// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

// ISCSIPersistentVolumeSourceApplyConfiguration represents an declarative configuration of the ISCSIPersistentVolumeSource type for use
// with apply.
type ISCSIPersistentVolumeSourceApplyConfiguration struct {
	TargetPortal      *string                            `json:"targetPortal,omitempty"`
	IQN               *string                            `json:"iqn,omitempty"`
	Lun               *int32                             `json:"lun,omitempty"`
	ISCSIInterface    *string                            `json:"iscsiInterface,omitempty"`
	FSType            *string                            `json:"fsType,omitempty"`
	ReadOnly          *bool                              `json:"readOnly,omitempty"`
	Portals           []string                           `json:"portals,omitempty"`
	DiscoveryCHAPAuth *bool                              `json:"chapAuthDiscovery,omitempty"`
	SessionCHAPAuth   *bool                              `json:"chapAuthSession,omitempty"`
	SecretRef         *SecretReferenceApplyConfiguration `json:"secretRef,omitempty"`
	InitiatorName     *string                            `json:"initiatorName,omitempty"`
}

// ISCSIPersistentVolumeSourceApplyConfiguration constructs an declarative configuration of the ISCSIPersistentVolumeSource type for use with
// apply.
func ISCSIPersistentVolumeSource() *ISCSIPersistentVolumeSourceApplyConfiguration {
	return &ISCSIPersistentVolumeSourceApplyConfiguration{}
}

// WithTargetPortal sets the TargetPortal field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TargetPortal field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithTargetPortal(value string) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.TargetPortal = &value
	return b
}

// WithIQN sets the IQN field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IQN field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithIQN(value string) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.IQN = &value
	return b
}

// WithLun sets the Lun field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Lun field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithLun(value int32) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.Lun = &value
	return b
}

// WithISCSIInterface sets the ISCSIInterface field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ISCSIInterface field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithISCSIInterface(value string) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.ISCSIInterface = &value
	return b
}

// WithFSType sets the FSType field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FSType field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithFSType(value string) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.FSType = &value
	return b
}

// WithReadOnly sets the ReadOnly field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ReadOnly field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithReadOnly(value bool) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.ReadOnly = &value
	return b
}

// WithPortals adds the given value to the Portals field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Portals field.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithPortals(values ...string) *ISCSIPersistentVolumeSourceApplyConfiguration {
	for i := range values {
		b.Portals = append(b.Portals, values[i])
	}
	return b
}

// WithDiscoveryCHAPAuth sets the DiscoveryCHAPAuth field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DiscoveryCHAPAuth field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithDiscoveryCHAPAuth(value bool) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.DiscoveryCHAPAuth = &value
	return b
}

// WithSessionCHAPAuth sets the SessionCHAPAuth field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SessionCHAPAuth field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithSessionCHAPAuth(value bool) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.SessionCHAPAuth = &value
	return b
}

// WithSecretRef sets the SecretRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SecretRef field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithSecretRef(value *SecretReferenceApplyConfiguration) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.SecretRef = value
	return b
}

// WithInitiatorName sets the InitiatorName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the InitiatorName field is set to the value of the last call.
func (b *ISCSIPersistentVolumeSourceApplyConfiguration) WithInitiatorName(value string) *ISCSIPersistentVolumeSourceApplyConfiguration {
	b.InitiatorName = &value
	return b
}
