package types

import (
	"encoding/json"
	"time"
)

const (
	// ServiceRelationshipTypeImplicit indicates the service dependency is auto created by resource reference.
	ServiceRelationshipTypeImplicit = "Implicit"
	// ServiceRelationshipTypeExplicit indicates the service dependency is manually created by user.
	ServiceRelationshipTypeExplicit = "Explicit"
)

// ServiceDriftResult indicates the drift detection result of the service.
type ServiceDriftResult struct {
	// Drifted indicates whether the service is drifted.
	Drifted bool `json:"drifted"`
	// Time indicates the time when the service is detected.
	Time time.Time `json:"time"`
	// Result indicates the drift result of service.
	Result *ServiceDrift `json:"drift"`
}

type ServiceDrift struct {
	FormatVersion      string                `json:"format_version"`
	TerraformVersion   string                `json:"terraform_version"`
	PlannedValues      *PlannedValues        `json:"-"`
	ResourceDrifts     []*ResourceDrift      `json:"resource_drift"`
	Configuration      *Configuration        `json:"-"`
	RelevantAttributes []*RelevantAttributes `json:"-"`
	Timestamp          time.Time             `json:"timestamp"`
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
