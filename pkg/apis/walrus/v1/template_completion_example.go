package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TemplateCompletionExample is the schema for the template completion examples API.
//
// +genclient
// +genclient:onlyVerbs=get,list
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"]
type TemplateCompletionExample struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Status TemplateCompletionExampleStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*TemplateCompletionExample)(nil)

// TemplateCompletionExampleStatus defines the observed state of TemplateCompletionExample.
type TemplateCompletionExampleStatus struct {
	// Purpose is the purpose of the template completion example.
	Purpose string `json:"purpose,omitempty"`

	// Prompt is the prompt of the template completion example.
	Prompt string `json:"prompt,omitempty"`
}

// TemplateCompletionExampleList holds the list of TemplateCompletionExample.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TemplateCompletionExampleList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []TemplateCompletionExample `json:"items"`
}

var _ runtime.Object = (*TemplateCompletionExampleList)(nil)
