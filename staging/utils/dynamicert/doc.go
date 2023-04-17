// Package dynamicert is used to handle the dynamic TLS generation in a lazy mode.
//   - most of the logic is the same as golang.org/x/crypto/acme/autocert.
//   - dynamicert generates the server certificate for each reachable request,
//     do not renew the generated certificate even the CA has changed.
//   - dynamicert (re)generates the server certificate if not found or the previous one has expired.
//   - dynamicert is not designed for production usage.
package dynamicert
