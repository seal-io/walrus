package types

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/seal-io/walrus/pkg/dao/types/property"
)

type OutputValue struct {
	Name   string          `json:"name,omitempty"`
	Value  property.Value  `json:"value,omitempty"`
	Schema openapi3.Schema `json:"schema,omitempty"`
}

const (
	ResourceRevisionTypeApply   = "apply"
	ResourceRevisionTypeDestory = "destroy"
	// ResourceRevisionTypeSync try to sync the resource from remote and update the revision state to match.
	ResourceRevisionTypeSync = "sync"
	// ResourceRevisionTypeDetect try to detect drift of the resource from remote.
	ResourceRevisionTypeDetect = "detect"
)
