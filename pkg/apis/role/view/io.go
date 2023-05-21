package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs.

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Role, role.OrderOption] `query:",inline"`

	Kind string `query:"kind,omitempty"`
}

func (r *CollectionGetRequest) Validate() error {
	if r.Kind != "" && !types.IsRoleKind(r.Kind) {
		return errors.New("invalid kind: unknown")
	}

	return nil
}

type CollectionGetResponse = []*model.RoleOutput

// Extensional APIs.
