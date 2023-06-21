package serviceresources

import (
	"context"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	optypes "github.com/seal-io/seal/pkg/operator/types"
)

// Label applies the labels to the given model.ServiceResource list with the given operator.Operator.
func Label(ctx context.Context, op optypes.Operator, candidates []*model.ServiceResource) (berr error) {
	if op == nil {
		return
	}

	for i := range candidates {
		// Get label values.
		var (
			// Name.
			svcName     string
			projectName string
			envName     string
		)

		if ins := candidates[i].Edges.Service; ins == nil {
			continue
		} else {
			// Service name.
			svcName = ins.Name

			// Project name.
			if proj := ins.Edges.Project; proj != nil {
				projectName = proj.Name
			}

			// Environment name.
			if env := ins.Edges.Environment; env != nil {
				envName = env.Name
			}
		}

		ls := map[string]string{
			// Name.
			types.LabelSealEnvironmentName: envName,
			types.LabelSealProjectName:     projectName,
			types.LabelSealServiceName:     svcName,
		}

		err := op.Label(ctx, candidates[i], ls)
		if multierr.AppendInto(&berr, err) {
			continue
		}
	}

	return berr
}
