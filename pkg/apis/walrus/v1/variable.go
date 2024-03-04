package v1

import (
	"reflect"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Variable is the schema for the variables API.
//
// +genclient
// +genclient:onlyVerbs=create,get,list,watch,apply,update,patch,delete,deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["var"]
type Variable struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   VariableSpec   `json:"spec,omitempty"`
	Status VariableStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*Variable)(nil)

// VariableSpec defines the desired state of Variable.
type VariableSpec struct {
	// Value contains the configuration data,
	// it is provided as a write-only input field.
	Value *string `json:"value,omitempty"`

	// Sensitive indicates whether the variable is sensitive.
	Sensitive bool `json:"sensitive"`
}

// VariableStatus defines the observed state of Variable.
type VariableStatus struct {
	// Project is the project that the variable belongs to.
	Project string `json:"project"`

	// Environment is the environment that the variable belongs to.
	Environment string `json:"environment"`

	// Value is the current value of the setting,
	// it is provided as a read-only output field.
	//
	// "(sensitive)" returns if the variable is sensitive.
	Value string `json:"value"`
	// Value_ is the shadow of the Value,
	// it is provided for system processing only.
	//
	// DO NOT EXPOSE AND STORE IT.
	Value_ string `json:"-"`
}

func (in *Variable) Equal(in2 *Variable) bool {
	return reflect.DeepEqual(in.Spec, in2.Spec) &&
		in.Status.Value_ == in2.Status.Value_
}

// VariableList holds the list of Variable.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VariableList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []Variable `json:"items"`
}

var _ runtime.Object = (*VariableList)(nil)
