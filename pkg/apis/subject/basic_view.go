package subject

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs.

type (
	CreateRequest struct {
		model.SubjectCreateInput `path:",inline" json:",inline"`

		Password string `json:"password"`
	}

	CreateResponse = *model.SubjectOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.SubjectCreateInput.Validate(); err != nil {
		return err
	}

	if !types.IsSubjectKind(r.Kind) {
		return errors.New("invalid kind: unknown")
	}

	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	if r.Password == "" {
		return errors.New("invalid password: blank")
	}

	r.Builtin = false

	return nil
}

type UpdateRequest struct {
	model.SubjectUpdateInput `path:",inline" json:",inline"`

	Password string `json:"password,omitempty"` // Allow to reset password.
}

type DeleteRequest = model.SubjectDeleteInput

type (
	CollectionGetRequest struct {
		model.SubjectQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Subject, subject.OrderOption,
		] `query:",inline"`

		Kind string `query:"kind,omitempty"`
	}

	CollectionGetResponse = []*model.SubjectOutput
)

func (r *CollectionGetRequest) Validate() error {
	if err := r.SubjectQueryInputs.Validate(); err != nil {
		return err
	}

	if r.Kind != "" && !types.IsSubjectKind(r.Kind) {
		return errors.New("invalid kind: unknown")
	}

	return nil
}
