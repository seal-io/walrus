package role

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/types"
)

type (
	CollectionGetRequest struct {
		model.RoleQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Role, role.OrderOption,
		] `query:",inline"`

		Kind string `query:"kind,omitempty"`
	}

	CollectionGetResponse = []*model.RoleOutput
)

func (r *CollectionGetRequest) Validate() error {
	if err := r.RoleQueryInputs.Validate(); err != nil {
		return err
	}

	if r.Kind != "" && !types.IsRoleKind(r.Kind) {
		return errors.New("invalid kind: unknown")
	}

	return nil
}
