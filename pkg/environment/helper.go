package environment

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

func GetManagedNamespaceName(e *model.Environment) string {
	if e == nil || e.Edges.Project == nil {
		return ""
	}

	if e.Annotations[types.AnnotationEnableManagedNamespace] == "false" {
		return ""
	}

	if e.Annotations[types.AnnotationManagedNamespace] != "" {
		return e.Annotations[types.AnnotationManagedNamespace]
	}

	return fmt.Sprintf("%s-%s", e.Edges.Project.Name, e.Name)
}

// ShortenEnvironmentNameIfNeeded shortens environment name if the combined length with project name is greater than 63.
// For maximum length of a namespace name is 63 characters.
func ShortenEnvironmentNameIfNeeded(name, projectName string) string {
	namespace := fmt.Sprintf("%s-%s", projectName, name)
	if len(namespace) <= 63 {
		return name
	}

	return name[:63-len(projectName)-1]
}
