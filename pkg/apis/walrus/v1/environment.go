package v1

import (
	"errors"
	"reflect"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Environment is the schema for the environments API.
//
// +genclient
// +genclient:onlyVerbs=create,get,list,watch,apply,update,patch,delete,deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["env"]
type Environment struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnvironmentSpec   `json:"spec,omitempty"`
	Status EnvironmentStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*Environment)(nil)

// EnvironmentType describes the type of environment.
// +enum
type EnvironmentType string

const (
	// EnvironmentTypeDevelopment means the environment is for development.
	EnvironmentTypeDevelopment EnvironmentType = "development"
	// EnvironmentTypeStaging means the environment is for staging.
	EnvironmentTypeStaging EnvironmentType = "staging"
	// EnvironmentTypeProduction means the environment is for production.
	EnvironmentTypeProduction EnvironmentType = "production"
)

func (in EnvironmentType) String() string {
	return string(in)
}

func (in EnvironmentType) Validate() error {
	switch in {
	case EnvironmentTypeDevelopment, EnvironmentTypeStaging, EnvironmentTypeProduction:
		return nil
	default:
		return errors.New("invalid environment type")
	}
}

// EnvironmentSpec defines the desired state of Environment.
type EnvironmentSpec struct {
	// Type is the type of the environment.
	//
	// +k8s:validation:enum=["development","staging","production"]
	Type EnvironmentType `json:"type"`

	// DisplayName is the display name of the environment.
	DisplayName string `json:"displayName,omitempty"`

	// Description is the description of the environment.
	Description string `json:"description,omitempty"`
}

// EnvironmentStatus defines the observed state of Environment.
type EnvironmentStatus struct {
	// Project is the project that the environment belongs to.
	Project string `json:"project"`

	// Phase is the current phase of the environment.
	Phase core.NamespacePhase `json:"phase"`
}

func (in *Environment) Equal(in2 *Environment) bool {
	return reflect.DeepEqual(in.Spec, in2.Spec) &&
		in.Status.Phase == in2.Status.Phase
}

// EnvironmentList holds the list of Environment.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EnvironmentList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []Environment `json:"items"`
}

var _ runtime.Object = (*EnvironmentList)(nil)
