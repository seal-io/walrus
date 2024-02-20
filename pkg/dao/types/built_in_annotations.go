package types

const (
	// AnnotationEnableManagedNamespace specify whether Walrus-managed namespace is enabled.
	// Defaults to true.
	AnnotationEnableManagedNamespace = "walrus.seal.io/enable-managed-namespace"
	// AnnotationManagedNamespace specify custom environment namespace name.
	AnnotationManagedNamespace = "walrus.seal.io/managed-namespace-name"

	// AnnotationSubjectID specify the subject ID of the system resource.
	AnnotationSubjectID = "walrus.seal.io/subject-id"
)
