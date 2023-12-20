package types

import (
	"bytes"

	"gopkg.in/yaml.v2"
	k8syaml "sigs.k8s.io/yaml"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/json"
)

const (
	WalrusFileVersion = "v1"
)

// SequenceWalrusFileKeys is the sequence of walrus file keys.
var SequenceWalrusFileKeys = []string{"version", "resources"}

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

	yb, err := k8syaml.JSONToYAML(b)
	if err != nil {
		return nil, err
	}

	m := make(map[string]any)

	err = yaml.Unmarshal(yb, &m)
	if err != nil {
		return nil, err
	}

	// Generate key sorted yaml.
	sorted := make(yaml.MapSlice, len(m))

	for i, v := range SequenceWalrusFileKeys {
		if _, ok := m[v]; !ok {
			continue
		}
		sorted[i] = yaml.MapItem{
			Key:   v,
			Value: m[v],
		}
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)

	err = enc.Encode(sorted)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ResourceSpec include the fields walrus file related.
type ResourceSpec struct {
	// Scope.
	Project     *MetadataName `json:"project,omitempty" yaml:"project,omitempty"`
	Environment *MetadataName `json:"environment,omitempty" yaml:"environment,omitempty"`

	Name        string            `json:"name" yaml:"name,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Attributes  property.Values   `json:"attributes,omitempty" yaml:"attributes,omitempty"`

	// Template or resource definition type.
	Type     string        `json:"type,omitempty" yaml:"type,omitempty"`
	Template *TemplateSpec `json:"template,omitempty" yaml:"template,omitempty"`
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
