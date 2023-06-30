package types

const (
	// ServiceRelationshipTypeImplicit indicates the service dependency is auto created by resource reference.
	ServiceRelationshipTypeImplicit = "Implicit"
	// ServiceRelationshipTypeExplicit indicates the service dependency is manually created by user.
	ServiceRelationshipTypeExplicit = "Explicit"
)
