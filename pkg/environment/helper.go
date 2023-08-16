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
