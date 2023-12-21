package types

import (
	"encoding/json"
	"time"
)

const (
	// ResourceRelationshipTypeImplicit indicates the resource dependency is auto created by resource reference.
	ResourceRelationshipTypeImplicit = "Implicit"
	// ResourceRelationshipTypeExplicit indicates the resource dependency is manually created by user.
	ResourceRelationshipTypeExplicit = "Explicit"
)

type ResourceDriftDetection struct {
	// Drifted indicates the resource is drifted or not.
	Drifted bool `json:"drifted"`
	// Time indicates the time when the resource is detected.
	Time time.Time `json:"time"`
	// Result indicates the drift result of resource.
	Result *ResourceDrift `json:"drift"`
}

type ResourceDrift struct {
	FormatVersion           string                    `json:"format_version"`
	TerraformVersion        string                    `json:"terraform_version"`
	PlannedValues           *PlannedValues            `json:"-"`
	ResourceComponentDrifts []*ResourceComponentDrift `json:"resource_drift"`
	Configuration           *Configuration            `json:"-"`
	RelevantAttributes      []*RelevantAttributes     `json:"-"`
	Timestamp               time.Time                 `json:"timestamp"`
}

type PlannedValues struct {
	Outputs    json.RawMessage `json:"outputs"`
	RootModule json.RawMessage `json:"root_module"`
}

type Configuration struct {
	ProviderConfig json.RawMessage `json:"provider_config"`
	RootModule     json.RawMessage `json:"root_module"`
}

type RelevantAttributes struct {
	Resource  string `json:"resource"`
	Attribute []any  `json:"attribute"`
}
