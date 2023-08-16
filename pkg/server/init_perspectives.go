package server

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model/perspective"
	"github.com/seal-io/walrus/pkg/dao/types"

	"github.com/seal-io/walrus/pkg/dao/model"
)

// createBuiltinPerspectives creates the built-in Perspective resources.
func (r *Server) createBuiltinPerspectives(ctx context.Context, opts initOptions) error {
	builtin := []*model.Perspective{
		perspectiveAll(),
		perspectiveCluster(),
		perspectiveProject(),
	}

	return opts.ModelClient.Perspectives().CreateBulk().
		Set(builtin...).
		OnConflictColumns(perspective.FieldName).
		UpdateNewValues().
		Exec(ctx)
}

func perspectiveAll() *model.Perspective {
	return &model.Perspective{
		Name:      "All",
		StartTime: "now-7d",
		EndTime:   "now",
		Builtin:   true,
		CostQueries: []types.QueryCondition{
			// Daily cost.
			{
				Filters: types.CostFilters{
					{
						{
							IncludeAll: true,
						},
					},
				},
				GroupBy: types.GroupByFieldDay,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 10,
				},
			},
			// Per project cost.
			{
				Filters: types.CostFilters{
					{
						{
							IncludeAll: true,
						},
					},
				},
				GroupBy: types.GroupByFieldProject,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 10,
				},
			},
			// Per cluster cost.
			{
				Filters: types.CostFilters{
					{
						{
							IncludeAll: true,
						},
					},
				},
				GroupBy: types.GroupByFieldConnectorID,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 10,
				},
			},
		},
	}
}

func perspectiveCluster() *model.Perspective {
	return &model.Perspective{
		Name:      "Cluster",
		StartTime: "now-7d",
		EndTime:   "now",
		Builtin:   true,
		CostQueries: []types.QueryCondition{
			// Daily cost.
			{
				Filters: types.CostFilters{
					{
						{
							FieldName: types.FilterFieldConnectorID,
							Operator:  types.OperatorIn,
							Values:    []string{"${connectorID}"},
						},
					},
				},
				GroupBy: types.GroupByFieldDay,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 10,
				},
			},
			// Per namespace cost.
			{
				Filters: types.CostFilters{
					{
						{
							FieldName: types.FilterFieldConnectorID,
							Operator:  types.OperatorIn,
							Values:    []string{"${connectorID}"},
						},
						{
							FieldName: types.FilterFieldName,
							Operator:  types.OperatorNotIn,
							Values: []string{
								types.ManagementCostItemName,
								types.IdleCostItemName,
							},
						},
					},
				},
				SharedOptions: &types.SharedCostOptions{
					Idle: &types.IdleShareOption{
						SharingStrategy: types.SharingStrategyProportionally,
					},
					Management: &types.ManagementShareOption{
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
				GroupBy: types.GroupByFieldNamespace,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 10,
				},
			},
			// Workload per day cost.
			{
				Filters: types.CostFilters{
					{
						{
							FieldName: types.FilterFieldConnectorID,
							Operator:  types.OperatorIn,
							Values:    []string{"${connectorID}"},
						},
						{
							FieldName: types.FilterFieldName,
							Operator:  types.OperatorNotIn,
							Values: []string{
								types.ManagementCostItemName,
								types.IdleCostItemName,
							},
						},
					},
				},
				SharedOptions: &types.SharedCostOptions{
					Idle: &types.IdleShareOption{
						SharingStrategy: types.SharingStrategyProportionally,
					},
					Management: &types.ManagementShareOption{
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
				GroupBy: types.GroupByFieldWorkload,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 10,
				},
				Step: types.StepDay,
			},
		},
	}
}

func perspectiveProject() *model.Perspective {
	return &model.Perspective{
		Name:      "Project",
		StartTime: "now-7d",
		EndTime:   "now",
		Builtin:   true,
		CostQueries: []types.QueryCondition{
			// Daily cost.
			{
				Filters: types.CostFilters{
					{
						{
							FieldName: types.FilterFieldProject,
							Operator:  types.OperatorIn,
							Values:    []string{"${project}"},
						},
					},
				},
				GroupBy: types.GroupByFieldDay,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 10,
				},
			},
			// Service cost.
			{
				Filters: types.CostFilters{
					{
						{
							FieldName: types.FilterFieldProject,
							Operator:  types.OperatorIn,
							Values:    []string{"${project}"},
						},
					},
				},
				GroupBy: types.GroupByFieldServicePath,
				Step:    types.StepDay,
			},
		},
	}
}
