package types

import "github.com/seal-io/seal/utils/slice"

const (
	SubjectKindUser  = "user"
	SubjectKindGroup = "group"
)

var SubjectKinds = []string{
	SubjectKindUser,
	SubjectKindGroup,
}

func IsSubjectKind(s string) bool {
	return slice.ContainsAny(SubjectKinds, s)
}

const SubjectDomainBuiltin = "builtin"
