package types

const (
	// AnnotationEnableManagedNamespace specify whether Seal-managed namespace is enabled.
	// Defaults to true.
	AnnotationEnableManagedNamespace = "seal.io/enable-managed-namespace"
	// AnnotationManagedNamespace specify custom environment namespace name.
	AnnotationManagedNamespace = "seal.io/managed-namespace-name"
)
