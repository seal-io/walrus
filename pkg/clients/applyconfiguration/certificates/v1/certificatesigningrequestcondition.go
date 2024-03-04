// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	v1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CertificateSigningRequestConditionApplyConfiguration represents an declarative configuration of the CertificateSigningRequestCondition type for use
// with apply.
type CertificateSigningRequestConditionApplyConfiguration struct {
	Type               *v1.RequestConditionType `json:"type,omitempty"`
	Status             *corev1.ConditionStatus  `json:"status,omitempty"`
	Reason             *string                  `json:"reason,omitempty"`
	Message            *string                  `json:"message,omitempty"`
	LastUpdateTime     *metav1.Time             `json:"lastUpdateTime,omitempty"`
	LastTransitionTime *metav1.Time             `json:"lastTransitionTime,omitempty"`
}

// CertificateSigningRequestConditionApplyConfiguration constructs an declarative configuration of the CertificateSigningRequestCondition type for use with
// apply.
func CertificateSigningRequestCondition() *CertificateSigningRequestConditionApplyConfiguration {
	return &CertificateSigningRequestConditionApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *CertificateSigningRequestConditionApplyConfiguration) WithType(value v1.RequestConditionType) *CertificateSigningRequestConditionApplyConfiguration {
	b.Type = &value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *CertificateSigningRequestConditionApplyConfiguration) WithStatus(value corev1.ConditionStatus) *CertificateSigningRequestConditionApplyConfiguration {
	b.Status = &value
	return b
}

// WithReason sets the Reason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reason field is set to the value of the last call.
func (b *CertificateSigningRequestConditionApplyConfiguration) WithReason(value string) *CertificateSigningRequestConditionApplyConfiguration {
	b.Reason = &value
	return b
}

// WithMessage sets the Message field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Message field is set to the value of the last call.
func (b *CertificateSigningRequestConditionApplyConfiguration) WithMessage(value string) *CertificateSigningRequestConditionApplyConfiguration {
	b.Message = &value
	return b
}

// WithLastUpdateTime sets the LastUpdateTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastUpdateTime field is set to the value of the last call.
func (b *CertificateSigningRequestConditionApplyConfiguration) WithLastUpdateTime(value metav1.Time) *CertificateSigningRequestConditionApplyConfiguration {
	b.LastUpdateTime = &value
	return b
}

// WithLastTransitionTime sets the LastTransitionTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastTransitionTime field is set to the value of the last call.
func (b *CertificateSigningRequestConditionApplyConfiguration) WithLastTransitionTime(value metav1.Time) *CertificateSigningRequestConditionApplyConfiguration {
	b.LastTransitionTime = &value
	return b
}
