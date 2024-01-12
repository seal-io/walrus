package types

import (
	"gopkg.in/yaml.v2"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/json"
)

const (
	WalrusFileVersion = "v1"
)

// WalrusFile is the spec for the WalrusFile.
type WalrusFile struct {
	Version   string         `json:"version,omitempty" yaml:"version,omitempty"`
	Resources []ResourceSpec `json:"resources,omitempty" yaml:"resources,omitempty"`
}

// Yaml returns the yaml bytes of the walrus file.
func (w *WalrusFile) Yaml() ([]byte, error) {
	if w == nil {
		return nil, nil
	}
	// Marshal to json since some field use json.RawMessage.
	b, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}

	// Convert to yaml mapSlice to keep the order of keys.
	var m yaml.MapSlice

	err = yaml.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	out, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// ResourceSpec include the fields walrus file related.
// The walrus file yaml preserves the order of the fields in this struct.
type ResourceSpec struct {
	Name        string `json:"name" yaml:"name,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Scope.
	Project     *MetadataName `json:"project,omitempty" yaml:"project,omitempty"`
	Environment *MetadataName `json:"environment,omitempty" yaml:"environment,omitempty"`

	// Template or resource definition type.
	Type     string        `json:"type,omitempty" yaml:"type,omitempty"`
	Template *TemplateSpec `json:"template,omitempty" yaml:"template,omitempty"`

	Labels     map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Attributes property.Values   `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// MetadataName is the name of the resource.
type MetadataName struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

// TemplateSpec include the field walrus file related.
type TemplateSpec struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}
