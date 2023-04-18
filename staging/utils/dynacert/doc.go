// Package dynacert is used to handle the dynamic TLS generation in a lazy mode.
//   - most of the logic is the same as golang.org/x/crypto/acme/autocert.
//   - dynacert generates the server certificate for each reachable request,
//     do not renew the generated certificate even the CA has changed.
//   - dynacert (re)generates the server certificate if not found or the previous one has expired.
//   - dynacert is not designed for production usage.
package dynacert
