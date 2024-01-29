package types

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/zclconf/go-cty/cty"

	"github.com/seal-io/walrus/pkg/dao/types/property"
)

type OutputValue struct {
	Name      string          `json:"name,omitempty"`
	Value     property.Value  `json:"value,omitempty,cli-table-column"`
	Type      cty.Type        `json:"type,omitempty"`
	Sensitive bool            `json:"sensitive,omitempty"`
	Schema    openapi3.Schema `json:"schema,omitempty"`
}

const (
	RunJobTypeApply   = "apply"
	RunJobTypeDestroy = "destroy"
)

type ResourceRunConfigData = []byte
