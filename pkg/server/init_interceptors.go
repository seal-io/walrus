package server

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/types"
)

func (r *Server) initInterceptors(_ context.Context, opts initOptions) error {
	opts.ModelClient.Intercept(
		model.InterceptFunc(func(next model.Querier) model.Querier {
			return model.QuerierFunc(func(ctx context.Context, query model.Query) (model.Value, error) {
				// add default filter for application resource.
				if q, ok := query.(*model.ApplicationResourceQuery); ok {
					q.Where(applicationresource.Mode(types.ApplicationResourceModeManaged))
				}

				return next.Query(ctx, query)
			})
		}),
	)

	return nil
}
