package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Setting is the schema for the settings API.
//
// +genclient
// +genclient:onlyVerbs=get,list,watch,apply,update,patch
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["set"]
type Setting struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   SettingSpec   `json:"spec,omitempty"`
	Status SettingStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*Setting)(nil)

// SettingSpec defines the desired state of Setting.
type SettingSpec struct {
	// Value contains the configuration data,
	// it is provided as a write-only input field.
	Value *string `json:"value,omitempty"`
}

// SettingStatus defines the observed state of Setting.
type SettingStatus struct {
	// Description is the description of the settings,
	// it is readonly.
	Description string `json:"description,omitempty"`

	// Hidden indicates whether the setting is hidden on UI,
	// it is readonly.
	Hidden bool `json:"hidden"`

	// Editable indicates whether the setting is editable on UI,
	// it is readonly.
	Editable bool `json:"editable"`

	// Sensitive indicates whether the setting is sensitive,
	// it is readonly.
	Sensitive bool `json:"sensitive"`

	// Value is the current value of the setting,
	// it is provided as a read-only output field.
	//
	// "(sensitive)" returns if the setting is sensitive.
	Value string `json:"value"`
	// Value_ is the shadow of the Value,
	// it is provided for system processing only.
	//
	// DO NOT EXPOSE AND STORE IT.
	Value_ string `json:"-"`
}

func (in *Setting) Equal(in2 *Setting) bool {
	return in.Status.Value_ == in2.Status.Value_
}

// SettingList holds the list of Setting.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SettingList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []Setting `json:"items"`
}

var _ runtime.Object = (*SettingList)(nil)
