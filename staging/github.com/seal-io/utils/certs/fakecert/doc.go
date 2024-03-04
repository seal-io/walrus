// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Package fakecert is used to handle the TLS generation in a lazy mode.
//   - most of the logic is the same as golang.org/x/crypto/acme/autocert.
//   - fakecert generates the host certificate for each reachable request,
//     do not renew the generated certificate even the CA has changed.
//   - fakecert (re)generates the host certificate if not found or the previous one has expired.
//   - fakecert is not designed for production usage.
package fakecert
