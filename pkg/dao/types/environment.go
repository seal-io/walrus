package types

import "golang.org/x/exp/slices"

var environmentTypes = []string{
	EnvironmentDevelopment,
	EnvironmentStaging,
	EnvironmentProduction,
}

func EnvironmentTypes() []string {
	return slices.Clone(environmentTypes)
}

func IsEnvironmentType(s string) bool {
	return slices.Contains(environmentTypes, s)
}

const (
	EnvironmentDevelopment = "development"
	EnvironmentStaging     = "staging"
	EnvironmentProduction  = "production"
)
