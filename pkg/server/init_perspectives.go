package server

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/types"

	"github.com/seal-io/seal/pkg/dao/model"
)

func (r *Server) initPerspectives(ctx context.Context, opts initOptions) error {
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
		AllocationQueries: []types.QueryCondition{
			// Daily cost.
			{
				Filters: types.AllocationCostFilters{
					{
						{
							IncludeAll: true,
						},
					},
				},
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								IncludeAll: true,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								IncludeAll: true,
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
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
				Filters: types.AllocationCostFilters{
					{
						{
							IncludeAll: true,
						},
					},
				},
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								IncludeAll: true,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								IncludeAll: true,
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
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
				Filters: types.AllocationCostFilters{
					{
						{
							IncludeAll: true,
						},
					},
				},
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								IncludeAll: true,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								IncludeAll: true,
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
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
		AllocationQueries: []types.QueryCondition{
			// Daily cost.
			{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: types.FilterFieldConnectorID,
							Operator:  types.OperatorIn,
							Values:    []string{"${connectorID}"},
						},
					},
				},
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								ConnectorID: "${connectorID}",
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								ConnectorID: "${connectorID}",
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
				GroupBy: types.GroupByFieldDay,
				Step:    types.StepDay,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 10,
				},
			},
			// Per namespace cost.
			{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: types.FilterFieldConnectorID,
							Operator:  types.OperatorIn,
							Values:    []string{"${connectorID}"},
						},
					},
				},
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								ConnectorID: "${connectorID}",
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								ConnectorID: "${connectorID}",
							},
						},
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
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: types.FilterFieldConnectorID,
							Operator:  types.OperatorIn,
							Values:    []string{"${connectorID}"},
						},
					},
				},
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								ConnectorID: "${connectorID}",
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								ConnectorID: "${connectorID}",
							},
						},
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
		AllocationQueries: []types.QueryCondition{
			// Service cost.
			{
				Filters: types.AllocationCostFilters{
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
