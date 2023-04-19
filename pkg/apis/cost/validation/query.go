package validation

import (
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/slice"
)

func ValidateAllocationQueries(queries []types.QueryCondition) error {
	if len(queries) == 0 {
		return errors.New("invalid allocation queries: blank")
	}

	for i := range queries {
		err := ValidateAllocationQuery(queries[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateAllocationQuery(query types.QueryCondition) error {
	// filter
	if len(query.Filters) == 0 {
		return errors.New("invalid filter: blank")
	}
	err := ValidateAllocationCostFilters(query.Filters)
	if err != nil {
		return err
	}

	// group by
	if query.GroupBy == "" {
		return errors.New("invalid group by: blank")
	}
	if !slice.ContainsAny([]types.GroupByField{
		types.GroupByFieldConnectorID,
		types.GroupByFieldNamespace,
		types.GroupByFieldClusterName,
		types.GroupByFieldNode,
		types.GroupByFieldController,
		types.GroupByFieldControllerKind,
		types.GroupByFieldPod,
		types.GroupByFieldContainer,
		types.GroupByFieldWorkload,
		types.GroupByFieldDay,
		types.GroupByFieldWeek,
		types.GroupByFieldMonth,
		types.GroupByFieldProject,
		types.GroupByFieldEnvironment,
		types.GroupByFieldApplication,
	}, query.GroupBy) {
		return errors.New("invalid group by: unsupported")
	}

	// step
	if query.Step != "" {
		// check support
		if !slice.ContainsAny([]types.Step{
			types.StepDay,
			types.StepWeek,
			types.StepMonth,
		}, query.Step) {
			return fmt.Errorf("invalid step: unsupported")
		}

		// check conflict with group by day bucket
		if slice.ContainsAny([]types.GroupByField{
			types.GroupByFieldDay,
			types.GroupByFieldWeek,
			types.GroupByFieldMonth,
		}, query.GroupBy) {
			return fmt.Errorf("invalid step: already group by %s", query.GroupBy)
		}
	}

	// share cost
	if len(query.SharedCosts) != 0 {
		err = ValidateShareCostFilters(query.SharedCosts)
		if err != nil {
			return err
		}
	}

	// page
	if query.Paging.Page < 0 {
		return fmt.Errorf("invalid page: negtive value")
	}
	if query.Paging.PerPage < 0 {
		return fmt.Errorf("invalid per page: negtive value")
	}
	return nil
}

func ValidateShareCostFilters(filters types.ShareCosts) error {
	for _, v := range filters {
		// allocation resource
		if len(v.Filters) != 0 {
			err := ValidateAllocationCostFilters(v.Filters)
			if err != nil {
				return err
			}
		}

		// management
		if len(v.ManagementCostFilters) != 0 {
			for _, mf := range v.ManagementCostFilters {
				if !mf.IncludeAll && !mf.ConnectorID.IsNaive() {
					return errors.New("invalid management share cost: blank connector id")
				}
			}
		}

		// idle cost
		if len(v.IdleCostFilters) != 0 {
			for _, idf := range v.IdleCostFilters {
				if !idf.IncludeAll && !idf.ConnectorID.IsNaive() {
					return errors.New("invalid idle share cost: blank connector id")
				}
			}
		}

		// share strategy
		if !slice.ContainsAny([]types.SharingStrategy{
			types.SharingStrategyProportionally,
			types.SharingStrategyEqually,
		}, v.SharingStrategy) {
			return fmt.Errorf("invalid share strategy: unsupported")
		}
	}
	return nil
}

func ValidateAllocationCostFilters(filters types.AllocationCostFilters) error {
	for _, condOr := range filters {
		if len(condOr) == 0 {
			return errors.New("invalid filter: blank")
		}
		for _, condAnd := range condOr {
			// include all
			if condAnd.IncludeAll {
				continue
			}

			// field name
			if condAnd.FieldName == "" {
				return errors.New("invalid filter: blank field name")
			}

			// operator
			if !slice.ContainsAny([]types.Operator{
				types.OperatorIn,
				types.OperatorNotIn,
			}, condAnd.Operator) {
				return errors.New("invalid filter: unsupported operator")
			}

			// values
			if len(condAnd.Values) == 0 {
				return errors.New("invalid filter: blank field values")
			}
		}
	}
	return nil
}
