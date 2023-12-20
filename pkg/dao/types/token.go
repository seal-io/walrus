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

const (
	WalrusOperationTokenName = "walrus-operation-token"
)

func IsTokenKind(s string) bool {
	return slices.Contains(TokenKinds, s)
}
