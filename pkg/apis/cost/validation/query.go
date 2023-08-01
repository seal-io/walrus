package validation

import (
	"errors"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/seal-io/seal/pkg/dao/types"
)

func ValidateCostQueries(queries []types.QueryCondition) error {
	if len(queries) == 0 {
		return errors.New("invalid cost queries: blank")
	}

	for i := range queries {
		err := ValidateCostQuery(queries[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateCostQuery(query types.QueryCondition) error {
	// Filter.
	if len(query.Filters) == 0 {
		return errors.New("invalid filter: blank")
	}

	err := ValidateCostFilters(query.Filters)
	if err != nil {
		return err
	}

	// Group by.
	if query.GroupBy == "" {
		return errors.New("invalid group by: blank")
	}

	if !slices.Contains([]types.GroupByField{
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
		types.GroupByFieldService,
		types.GroupByFieldEnvironmentPath,
		types.GroupByFieldServicePath,
	}, query.GroupBy) {
		return errors.New("invalid group by: unsupported")
	}

	// Step.
	if query.Step != "" {
		// Check support.
		if !slices.Contains([]types.Step{
			types.StepDay,
			types.StepWeek,
			types.StepMonth,
		}, query.Step) {
			return fmt.Errorf("invalid step: unsupported")
		}

		// Check conflict with group by day bucket.
		if slices.Contains([]types.GroupByField{
			types.GroupByFieldDay,
			types.GroupByFieldWeek,
			types.GroupByFieldMonth,
		}, query.GroupBy) {
			return fmt.Errorf("invalid step: already group by %s", query.GroupBy)
		}
	}

	// Share cost.
	if query.SharedOptions != nil {
		err = ValidateShareCostFilters(query.SharedOptions)
		if err != nil {
			return err
		}
	}

	// Page.
	if query.Paging.Page < 0 {
		return fmt.Errorf("invalid page: negtive value")
	}

	if query.Paging.PerPage < 0 {
		return fmt.Errorf("invalid per page: negtive value")
	}

	return nil
}

func ValidateShareCostFilters(options *types.SharedCostOptions) error {
	isValidStrategy := func(strategy types.SharingStrategy) bool {
		return slices.Contains([]types.SharingStrategy{
			types.SharingStrategyProportionally,
			types.SharingStrategyEqually,
		}, strategy)
	}

	for _, v := range options.Item {
		// Allocation resource.
		if len(v.Filters) != 0 {
			err := ValidateCostFilters(v.Filters)
			if err != nil {
				return err
			}
		}

		// Share strategy.
		if !isValidStrategy(v.SharingStrategy) {
			return fmt.Errorf("invalid share strategy: unsupported")
		}
	}

	// Management.
	if options.Management != nil && !isValidStrategy(options.Management.SharingStrategy) {
		return fmt.Errorf("invalid share strategy: unsupported")
	}

	// Idle cost.
	if options.Idle != nil && !isValidStrategy(options.Idle.SharingStrategy) {
		return fmt.Errorf("invalid share strategy: unsupported")
	}

	return nil
}

func ValidateCostFilters(filters types.CostFilters) error {
	for _, condOr := range filters {
		if len(condOr) == 0 {
			return errors.New("invalid filter: blank")
		}

		for _, condAnd := range condOr {
			// Include all.
			if condAnd.IncludeAll {
				continue
			}

			// Field name.
			if condAnd.FieldName == "" {
				return errors.New("invalid filter: blank field name")
			}

			// Operator.
			if !slices.Contains([]types.Operator{
				types.OperatorIn,
				types.OperatorNotIn,
			}, condAnd.Operator) {
				return errors.New("invalid filter: unsupported operator")
			}

			// Values.
			if len(condAnd.Values) == 0 {
				return errors.New("invalid filter: blank field values")
			}
		}
	}

	return nil
}
