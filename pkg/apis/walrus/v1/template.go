package v1

import (
	"k8s.io/apimachinery/pkg/runtime"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
)

// Template is the schema for the templates API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["tpl"],subResources=["status"]
type Template walruscore.Template

var _ runtime.Object = (*Template)(nil)

// TemplateList holds the list of Template.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TemplateList walruscore.TemplateList

var _ runtime.Object = (*TemplateList)(nil)
