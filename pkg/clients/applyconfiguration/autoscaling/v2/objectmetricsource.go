// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v2

// ObjectMetricSourceApplyConfiguration represents an declarative configuration of the ObjectMetricSource type for use
// with apply.
type ObjectMetricSourceApplyConfiguration struct {
	DescribedObject *CrossVersionObjectReferenceApplyConfiguration `json:"describedObject,omitempty"`
	Target          *MetricTargetApplyConfiguration                `json:"target,omitempty"`
	Metric          *MetricIdentifierApplyConfiguration            `json:"metric,omitempty"`
}

// ObjectMetricSourceApplyConfiguration constructs an declarative configuration of the ObjectMetricSource type for use with
// apply.
func ObjectMetricSource() *ObjectMetricSourceApplyConfiguration {
	return &ObjectMetricSourceApplyConfiguration{}
}

// WithDescribedObject sets the DescribedObject field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DescribedObject field is set to the value of the last call.
func (b *ObjectMetricSourceApplyConfiguration) WithDescribedObject(value *CrossVersionObjectReferenceApplyConfiguration) *ObjectMetricSourceApplyConfiguration {
	b.DescribedObject = value
	return b
}

// WithTarget sets the Target field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Target field is set to the value of the last call.
func (b *ObjectMetricSourceApplyConfiguration) WithTarget(value *MetricTargetApplyConfiguration) *ObjectMetricSourceApplyConfiguration {
	b.Target = value
	return b
}

// WithMetric sets the Metric field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Metric field is set to the value of the last call.
func (b *ObjectMetricSourceApplyConfiguration) WithMetric(value *MetricIdentifierApplyConfiguration) *ObjectMetricSourceApplyConfiguration {
	b.Metric = value
	return b
}