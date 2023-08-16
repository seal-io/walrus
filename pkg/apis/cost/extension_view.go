package cost

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	costvalidation "github.com/seal-io/walrus/pkg/apis/cost/validation"
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/costs/distributor"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/validation"
)

type CollectedTimeRange struct {
	FirstTime time.Time `json:"firstTime,omitempty"`
	LastTime  time.Time `json:"lastTime,omitempty"`
}

type (
	CollectionRouteGetCostReportsRequest struct {
		_ struct{} `route:"POST=/cost-reports"`

		types.QueryCondition `json:",inline"`

		StartTime time.Time `json:"startTime,omitempty"`
		EndTime   time.Time `json:"endTime,omitempty"`

		Context *gin.Context
	}

	CollectionRouteGetCostReportsResponse = *runtime.ResponseCollection
)

func (r *CollectionRouteGetCostReportsRequest) Validate() error {
	switch {
	case slices.Contains([]types.GroupByField{types.GroupByFieldDay, types.GroupByFieldWeek}, r.GroupBy):
		return validation.TimeRangeWithinYear(r.StartTime, r.EndTime)
	case slices.Contains([]types.Step{types.StepDay, types.StepWeek}, r.Step):
		return validation.TimeRangeWithinYear(r.StartTime, r.EndTime)
	case r.GroupBy == types.GroupByFieldMonth:
		return validation.TimeRangeWithinDecade(r.StartTime, r.EndTime)
	case r.Step == types.StepMonth:
		return validation.TimeRangeWithinDecade(r.StartTime, r.EndTime)
	}

	if len(r.Filters) == 0 {
		return errors.New("invalid filter: blank")
	}

	return costvalidation.ValidateCostQuery(r.QueryCondition)
}

func (r *CollectionRouteGetCostReportsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetSummaryCostsRequest struct {
		_ struct{} `route:"POST=/summary-costs"`

		StartTime time.Time `json:"startTime,omitempty"`
		EndTime   time.Time `json:"endTime,omitempty"`

		Context *gin.Context
	}

	CollectionRouteGetSummaryCostsResponse struct {
		Currency           int                 `json:"currency,omitempty"`
		TotalCost          float64             `json:"totalCost,omitempty"`
		AverageDailyCost   float64             `json:"averageDailyCost,omitempty"`
		ProjectCount       int                 `json:"projectCount,omitempty"`
		ClusterCount       int                 `json:"clusterCount,omitempty"`
		CollectedTimeRange *CollectedTimeRange `json:"collectedTimeRange,omitempty"`
	}
)

func (r *CollectionRouteGetSummaryCostsRequest) Validate() error {
	return validation.TimeRangeWithinYear(r.StartTime, r.EndTime)
}

func (r *CollectionRouteGetSummaryCostsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetSummaryClusterCostsRequest struct {
		_ struct{} `route:"POST=/summary-cluster-costs"`

		StartTime   time.Time `json:"startTime,omitempty"`
		EndTime     time.Time `json:"endTime,omitempty"`
		ConnectorID object.ID `json:"connectorID,omitempty"`

		Context *gin.Context
	}

	CollectionRouteGetSummaryClusterCostsResponse struct {
		Currency           int                 `json:"currency,omitempty"`
		TotalCost          float64             `json:"totalCost,omitempty"`
		AverageDailyCost   float64             `json:"averageDailyCost,omitempty"`
		ItemCost           float64             `json:"itemCost,omitempty"`
		ManagementCost     float64             `json:"managementCost,omitempty"`
		IdleCost           float64             `json:"idleCost,omitempty"`
		CollectedTimeRange *CollectedTimeRange `json:"collectedTimeRange,omitempty"`
	}
)

func (r *CollectionRouteGetSummaryClusterCostsRequest) Validate() error {
	if err := validation.TimeRangeWithinYear(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if !r.ConnectorID.Valid() {
		return errors.New("invalid connector id")
	}

	return nil
}

func (r *CollectionRouteGetSummaryClusterCostsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetSummaryProjectCostsRequest struct {
		_ struct{} `route:"POST=/summary-project-costs"`

		StartTime time.Time `json:"startTime,omitempty"`
		EndTime   time.Time `json:"endTime,omitempty"`
		Project   string    `json:"project,omitempty"`

		Context *gin.Context
	}

	CollectionRouteGetSummaryProjectCostsResponse struct {
		Currency           int                 `json:"currency,omitempty"`
		TotalCost          float64             `json:"totalCost,omitempty"`
		AverageDailyCost   float64             `json:"averageDailyCost,omitempty"`
		CollectedTimeRange *CollectedTimeRange `json:"collectedTimeRange,omitempty"`
	}
)

func (r *CollectionRouteGetSummaryProjectCostsRequest) Validate() error {
	if err := validation.TimeRangeWithinYear(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if r.Project == "" {
		return errors.New("invalid project: blank")
	}

	return nil
}

func (r *CollectionRouteGetSummaryProjectCostsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetSummaryQueriedCostsRequest struct {
		_ struct{} `route:"POST=/summary-queried-costs"`

		StartTime     time.Time                `json:"startTime,omitempty"`
		EndTime       time.Time                `json:"endTime,omitempty"`
		Filters       types.CostFilters        `json:"filters,omitempty"`
		SharedOptions *types.SharedCostOptions `json:"sharedOptions,omitempty"`

		Context *gin.Context
	}

	CollectionRouteGetSummaryQueriedCostsResponse struct {
		distributor.Cost `json:",inline"`

		AverageDailyCost   float64             `json:"averageDailyCost,omitempty"`
		ConnectorCount     int                 `json:"connectorCount,omitempty"`
		CollectedTimeRange *CollectedTimeRange `json:"collectedTimeRange,omitempty"`
	}
)

func (r *CollectionRouteGetSummaryQueriedCostsRequest) Validate() error {
	if err := validation.TimeRangeWithinDecade(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if len(r.Filters) == 0 {
		return errors.New("invalid filters: blank")
	}

	if err := costvalidation.ValidateCostFilters(r.Filters); err != nil {
		return err
	}

	if err := costvalidation.ValidateShareCostFilters(r.SharedOptions); err != nil {
		return err
	}

	return nil
}

func (r *CollectionRouteGetSummaryQueriedCostsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}
