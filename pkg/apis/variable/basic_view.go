package variable

import (
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/variable"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.VariableCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.VariableOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.VariableCreateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if r.Value == "" {
		return errors.New("invalid value: blank")
	}

	return nil
}

type UpdateRequest = model.VariableUpdateInput

type DeleteRequest = model.VariableDeleteInput

type (
	CollectionGetRequest struct {
		model.VariableQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Variable, variable.OrderOption,
		] `query:",inline"`

		IncludeInherited bool `query:"includeInherited,omitempty"`
	}

	CollectionGetResponse = []*model.VariableOutput
)

type CollectionDeleteRequest = model.VariableDeleteInputs

func exposeVariable(in *model.Variable) *model.VariableOutput {
	if in.Sensitive {
		in.Value = ""
	}

	return model.ExposeVariable(in)
}

func exposeVariables(in []*model.Variable) []*model.VariableOutput {
	out := make([]*model.VariableOutput, 0, len(in))

	for i := 0; i < len(in); i++ {
		o := exposeVariable(in[i])
		if o == nil {
			continue
		}

		out = append(out, o)
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
