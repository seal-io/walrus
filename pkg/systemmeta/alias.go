package systemmeta

import (
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

// MetaObject is the interface for the object with metadata.
type MetaObject = ctrlcli.Object
