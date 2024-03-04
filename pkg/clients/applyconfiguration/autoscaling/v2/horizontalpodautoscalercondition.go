// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v2

import (
	v2 "k8s.io/api/autoscaling/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HorizontalPodAutoscalerConditionApplyConfiguration represents an declarative configuration of the HorizontalPodAutoscalerCondition type for use
// with apply.
type HorizontalPodAutoscalerConditionApplyConfiguration struct {
	Type               *v2.HorizontalPodAutoscalerConditionType `json:"type,omitempty"`
	Status             *v1.ConditionStatus                      `json:"status,omitempty"`
	LastTransitionTime *metav1.Time                             `json:"lastTransitionTime,omitempty"`
	Reason             *string                                  `json:"reason,omitempty"`
	Message            *string                                  `json:"message,omitempty"`
}

// HorizontalPodAutoscalerConditionApplyConfiguration constructs an declarative configuration of the HorizontalPodAutoscalerCondition type for use with
// apply.
func HorizontalPodAutoscalerCondition() *HorizontalPodAutoscalerConditionApplyConfiguration {
	return &HorizontalPodAutoscalerConditionApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *HorizontalPodAutoscalerConditionApplyConfiguration) WithType(value v2.HorizontalPodAutoscalerConditionType) *HorizontalPodAutoscalerConditionApplyConfiguration {
	b.Type = &value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *HorizontalPodAutoscalerConditionApplyConfiguration) WithStatus(value v1.ConditionStatus) *HorizontalPodAutoscalerConditionApplyConfiguration {
	b.Status = &value
	return b
}

// WithLastTransitionTime sets the LastTransitionTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastTransitionTime field is set to the value of the last call.
func (b *HorizontalPodAutoscalerConditionApplyConfiguration) WithLastTransitionTime(value metav1.Time) *HorizontalPodAutoscalerConditionApplyConfiguration {
	b.LastTransitionTime = &value
	return b
}

// WithReason sets the Reason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reason field is set to the value of the last call.
func (b *HorizontalPodAutoscalerConditionApplyConfiguration) WithReason(value string) *HorizontalPodAutoscalerConditionApplyConfiguration {
	b.Reason = &value
	return b
}

// WithMessage sets the Message field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Message field is set to the value of the last call.
func (b *HorizontalPodAutoscalerConditionApplyConfiguration) WithMessage(value string) *HorizontalPodAutoscalerConditionApplyConfiguration {
	b.Message = &value
	return b
}
