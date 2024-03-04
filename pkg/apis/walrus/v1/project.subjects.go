package v1

import (
	"errors"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ProjectSubjects holds the list of ProjectSubject.
//
// ProjectSubjects is the subresource of Project.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["projsub"]
type ProjectSubjects struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	// +patchStrategy=merge
	// +patchMergeKey=name
	// +listType=map
	// +listMapKey=name
	Items []ProjectSubject `json:"items,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

var _ runtime.Object = (*ProjectSubjects)(nil)

// ProjectSubjectRole describes the role of project subject.
// +enum
type ProjectSubjectRole string

const (
	// ProjectSubjectRoleViewer is the subject role for project viewer.
	ProjectSubjectRoleViewer ProjectSubjectRole = "walrus-project-viewer"
	// ProjectSubjectRoleMember is the subject role for project member.
	ProjectSubjectRoleMember ProjectSubjectRole = "walrus-project-member"
	// ProjectSubjectRoleOwner is the subject role for project owner.
	ProjectSubjectRoleOwner ProjectSubjectRole = "walrus-project-owner"
)

func (in ProjectSubjectRole) String() string {
	return string(in)
}

func (in ProjectSubjectRole) Validate() error {
	switch in {
	case ProjectSubjectRoleViewer, ProjectSubjectRoleMember, ProjectSubjectRoleOwner:
		return nil
	default:
		return errors.New("invalid project subject role")
	}
}

// ProjectSubject is the schema for the project subject API.
type ProjectSubject struct {
	/* NB(thxCode): All attributes must be comparable. */

	// Name is the name of the subject.
	Name string `json:"name"`

	// Kind is the kind of the subject.
	Kind string `json:"kind,omitempty"`

	// Role is the role of the subject.
	//
	// +k8s:validation:enum=["walrus-project-viewer","walrus-project-member","walrus-project-owner"]
	Role ProjectSubjectRole `json:"role"`
}
