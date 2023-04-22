package platformtf

import (
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/utils/json"
)

type state struct {
	Version          int                    `json:"version"`
	TerraformVersion string                 `json:"terraform_version"`
	Serial           int                    `json:"serial"`
	Lineage          string                 `json:"lineage"`
	Outputs          map[string]outputState `json:"outputs"`
	Resources        []resourceState        `json:"resources"`
	CheckResults     checkResults           `json:"check_results"`
}

type outputState struct {
	Value     property.Value `json:"value"`
	Type      property.Type  `json:"type"`
	Sensitive bool           `json:"sensitive,omitempty"`
}

type resourceState struct {
	Module    string                `json:"module,omitempty"`
	Mode      string                `json:"mode"`
	Type      string                `json:"type"`
	Name      string                `json:"name"`
	EachMode  string                `json:"each,omitempty"`
	Provider  string                `json:"provider"`
	Instances []instanceObjectState `json:"instances"`
}

type instanceObjectState struct {
	IndexKey interface{} `json:"index_key,omitempty"`
	Status   string      `json:"status,omitempty"`
	Deposed  string      `json:"deposed,omitempty"`

	SchemaVersion           uint64            `json:"schema_version"`
	AttributesRaw           json.RawMessage   `json:"attributes,omitempty"`
	AttributesFlat          map[string]string `json:"attributes_flat,omitempty"`
	AttributeSensitivePaths json.RawMessage   `json:"sensitive_attributes,omitempty"`

	PrivateRaw []byte `json:"private,omitempty"`

	Dependencies []string `json:"dependencies,omitempty"`

	CreateBeforeDestroy bool `json:"create_before_destroy,omitempty"`
}

type checkResults struct {
	ObjectKind string               `json:"object_kind"`
	ConfigAddr string               `json:"config_addr"`
	Status     string               `json:"status"`
	Objects    []checkResultsObject `json:"objects"`
}

func (c checkResults) MarshalJSON() ([]byte, error) {
	// if c is empty return empty json
	if c.ObjectKind == "" && c.ConfigAddr == "" && c.Status == "" && c.Objects == nil {
		return []byte("null"), nil
	}

	return json.Marshal(c)
}

type checkResultsObject struct {
	ObjectAddr      string   `json:"object_addr"`
	Status          string   `json:"status"`
	FailureMessages []string `json:"failure_messages,omitempty"`
}
