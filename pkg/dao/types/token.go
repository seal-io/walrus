package types

import "golang.org/x/exp/slices"

const (
	TokenKindDeployment = "deployment"
	TokenKindAPI        = "api"
)

var TokenKinds = []string{
	TokenKindDeployment,
	TokenKindAPI,
}

func IsTokenKind(s string) bool {
	return slices.Contains(TokenKinds, s)
}
