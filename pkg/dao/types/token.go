package types

import "github.com/seal-io/seal/utils/slice"

const (
	TokenKindDeployment = "deployment"
	TokenKindAPI        = "api"
)

var TokenKinds = []string{
	TokenKindDeployment,
	TokenKindAPI,
}

func IsTokenKind(s string) bool {
	return slice.ContainsAny(TokenKinds, s)
}
