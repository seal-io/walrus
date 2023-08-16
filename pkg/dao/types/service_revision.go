package types

import "github.com/seal-io/walrus/pkg/dao/types/property"

type OutputValue struct {
	Name      string         `json:"name,omitempty"`
	Value     property.Value `json:"value,omitempty"`
	Type      property.Type  `json:"type,omitempty"`
	Sensitive bool           `json:"sensitive,omitempty"`
}
