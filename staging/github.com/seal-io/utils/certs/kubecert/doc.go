// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Package kubecert is used to handle the TLS generation in a lazy mode.
//   - most of the logic is the same as golang.org/x/crypto/acme/autocert.
//   - kubecert generates the server certificate at first request arrival,
//     and it relies on the Kubernetes CertificateSigningRequest.
//   - kubecert (re)generates the server certificate if not found or the previous one has expired.
//   - kubecert is designed for Kubernetes APIServer proxy calling,
//     like Mutating/Validating Webhook, Extension APIServer.
package kubecert
