package types

const (
	// ServiceDependencyTypeImplicit indicates the service dependency is auto created by resource reference.
	ServiceDependencyTypeImplicit = "Implicit"
	// ServiceDependencyTypeExplicit indicates the service dependency is manually created by user.
	ServiceDependencyTypeExplicit = "Explicit"
)
