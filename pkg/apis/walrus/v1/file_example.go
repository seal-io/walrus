package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// FileExample is the schema for the file example API.
//
// +genclient
// +genclient:onlyVerbs=get,list
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"]
type FileExample struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Status FileExampleStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*FileExample)(nil)

// FileExampleStatus defines the observed state of FileExample.
type FileExampleStatus struct {
	// Icon is the icon of the file example.
	Icon string `json:"icon,omitempty"`
	// Readme is the readme of the file example.
	Readme string `json:"readme,omitempty"`
	// Content is the content of the file example.
	Content string `json:"content,omitempty"`
}

// FileExampleList holds the list of FileExample.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FileExampleList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []FileExample `json:"items"`
}

var _ runtime.Object = (*FileExampleList)(nil)
