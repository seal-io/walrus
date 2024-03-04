package apireg_gen

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DummyOthers is the schema for the dummy sub resources API.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DummyOthers struct {
	meta.TypeMeta `json:",inline"`
}
