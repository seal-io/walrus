package v1

import (
	"errors"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TemplateCompletion is the schema for the template completions API.
//
// +genclient
// +genclient:onlyVerbs=create
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"]
type TemplateCompletion struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplateCompletionSpec   `json:"spec,omitempty"`
	Status TemplateCompletionStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*TemplateCompletion)(nil)

// TemplateCompletionPurpose describes the purpose of template completion.
// +enum
type TemplateCompletionPurpose string

const (
	// TemplateCompletionPurposeGenerate means the template completion is for generating.
	TemplateCompletionPurposeGenerate TemplateCompletionPurpose = "generate"
	// TemplateCompletionPurposeCorrect means the template completion is for correcting.
	TemplateCompletionPurposeCorrect TemplateCompletionPurpose = "correct"
	// TemplateCompletionPurposeExplain means the template completion is for explaining.
	TemplateCompletionPurposeExplain TemplateCompletionPurpose = "explain"
)

func (in TemplateCompletionPurpose) String() string {
	return string(in)
}

func (in TemplateCompletionPurpose) Validate() error {
	switch in {
	case TemplateCompletionPurposeGenerate, TemplateCompletionPurposeCorrect, TemplateCompletionPurposeExplain:
		return nil
	default:
		return errors.New("invalid template completion purpose")
	}
}

// TemplateCompletionSpec defines the desired state of TemplateCompletion.
type TemplateCompletionSpec struct {
	// Purpose is the purpose of template completion.
	//
	// +k8s:validation:enum=["generate","correct","explain"]
	Purpose TemplateCompletionPurpose `json:"purpose"`

	// Content contains the content for completion,
	// it is provided as a write-only input field.
	Content string `json:"content"`
}

// TemplateCompletionStatus defines the observed state of TemplateCompletion.
type TemplateCompletionStatus struct {
	// Content contains the content for completion,
	// it is provided as a read-only output field.
	Content string `json:"content"`
}

// TemplateCompletionList holds the list of TemplateCompletion.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TemplateCompletionList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []FileExample `json:"items"`
}

var _ runtime.Object = (*TemplateCompletionList)(nil)
