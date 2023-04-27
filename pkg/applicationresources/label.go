package applicationresources

import (
	"context"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// Label applies the labels to the given model.ApplicationResource list with the given operator.Operator.
func Label(ctx context.Context, op operator.Operator, candidates []*model.ApplicationResource) (berr error) {
	if op == nil {
		return
	}

	for i := range candidates {
		// get label values.
		var (
			appName     string
			projectName string
			envName     string
		)
		if ins := candidates[i].Edges.Instance; ins == nil {
			continue
		} else {
			// application name
			if app := ins.Edges.Application; app != nil {
				appName = app.Name
				// project name
				if proj := app.Edges.Project; proj != nil {
					projectName = proj.Name
				}
			}
			// environment name
			if env := ins.Edges.Environment; env != nil {
				envName = env.Name
			}
		}

		var ls = map[string]string{
			types.LabelSealEnvironment: envName,
			types.LabelSealProject:     projectName,
			types.LabelSealApplication: appName,
		}
		var err = op.Label(ctx, candidates[i], ls)
		if multierr.AppendInto(&berr, err) {
			continue
		}
	}
	return
}
