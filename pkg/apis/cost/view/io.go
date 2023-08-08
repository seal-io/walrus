package view

import (
	"errors"
	"time"

	costvalidation "github.com/seal-io/seal/pkg/apis/cost/validation"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/utils/slice"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

// Batch APIs.

// Extensional APIs.

type (
	Resources []Resource
	Resource  struct {
		Cost

		ItemName  string     `json:"itemName,omitempty"`
		StartTime *time.Time `json:"startTime,omitempty"`
		EndTime   *time.Time `json:"endTime,omitempty"`
	}

	Cost struct {
		Currency         int     `json:"currency,omitempty"`
		TotalCost        float64 `json:"totalCost,omitempty"`
		SharedCost       float64 `json:"sharedCost,omitempty"`
		CpuCost          float64 `json:"cpuCost,omitempty"`
		GpuCost          float64 `json:"gpuCost,omitempty"`
		RamCost          float64 `json:"ramCost,omitempty"`
		PvCost           float64 `json:"pvCost,omitempty"`
		LoadBalancerCost float64 `json:"loadBalancerCost,omitempty"`
	}
)

type (
	CostReportRequest struct {
		_ struct{} `route:"POST=/cost-reports"`

		types.QueryCondition `json:",inline"`

		StartTime time.Time `json:"startTime,omitempty"`
		EndTime   time.Time `json:"endTime,omitempty"`
	}
)

func (r *CostReportRequest) Validate() error {
	if err := r.validateTimeRange(); err != nil {
		return err
	}

	if len(r.Filters) == 0 {
		return errors.New("invalid filter: blank")
	}

	return costvalidation.ValidateCostQuery(r.QueryCondition)
}

func (r *CostReportRequest) validateTimeRange() error {
	switch {
	case slice.ContainsAny([]types.GroupByField{types.GroupByFieldDay, types.GroupByFieldWeek}, r.GroupBy):
		return validation.TimeRangeWithinYear(r.StartTime, r.EndTime)
	case slice.ContainsAny([]types.Step{types.StepDay, types.StepWeek}, r.Step):
		return validation.TimeRangeWithinYear(r.StartTime, r.EndTime)
	case r.GroupBy == types.GroupByFieldMonth:
		return validation.TimeRangeWithinDecade(r.StartTime, r.EndTime)
	case r.Step == types.StepMonth:
		return validation.TimeRangeWithinDecade(r.StartTime, r.EndTime)
	}

	return nil
}

type CollectedTimeRange struct {
	FirstTime time.Time `json:"firstTime,omitempty"`
	LastTime  time.Time `json:"lastTime,omitempty"`
}

type SummaryCostRequest struct {
	_ struct{} `route:"POST=/summary-costs"`

	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
}

func (r *SummaryCostRequest) Validate() error {
	return validation.TimeRangeWithinYear(r.StartTime, r.EndTime)
}

type SummaryCostResponse struct {
	Currency           int                 `json:"currency,omitempty"`
	TotalCost          float64             `json:"totalCost,omitempty"`
	AverageDailyCost   float64             `json:"averageDailyCost,omitempty"`
	ProjectCount       int                 `json:"projectCount,omitempty"`
	ClusterCount       int                 `json:"clusterCount,omitempty"`
	CollectedTimeRange *CollectedTimeRange `json:"collectedTimeRange,omitempty"`
}

type SummaryClusterCostRequest struct {
	_ struct{} `route:"POST=/summary-cluster-costs"`

	StartTime   time.Time `json:"startTime,omitempty"`
	EndTime     time.Time `json:"endTime,omitempty"`
	ConnectorID object.ID `json:"connectorID,omitempty"`
}

func (r *SummaryClusterCostRequest) Validate() error {
	if err := validation.TimeRangeWithinYear(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if !r.ConnectorID.Valid() {
		return errors.New("invalid connector id")
	}

	return nil
}

type SummaryClusterCostResponse struct {
	Currency           int                 `json:"currency,omitempty"`
	TotalCost          float64             `json:"totalCost,omitempty"`
	AverageDailyCost   float64             `json:"averageDailyCost,omitempty"`
	ItemCost           float64             `json:"itemCost,omitempty"`
	ManagementCost     float64             `json:"managementCost,omitempty"`
	IdleCost           float64             `json:"idleCost,omitempty"`
	CollectedTimeRange *CollectedTimeRange `json:"collectedTimeRange,omitempty"`
}

type SummaryProjectCostRequest struct {
	_ struct{} `route:"POST=/summary-project-costs"`

	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
	Project   string    `json:"project,omitempty"`
}

func (r *SummaryProjectCostRequest) Validate() error {
	if err := validation.TimeRangeWithinYear(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if r.Project == "" {
		return errors.New("invalid project: blank")
	}

	return nil
}

type SummaryCostCommonResponse struct {
	Currency           int                 `json:"currency,omitempty"`
	TotalCost          float64             `json:"totalCost,omitempty"`
	AverageDailyCost   float64             `json:"averageDailyCost,omitempty"`
	CollectedTimeRange *CollectedTimeRange `json:"collectedTimeRange,omitempty"`
}

type SummaryQueriedCostRequest struct {
	_ struct{} `route:"POST=/summary-queried-costs"`

	StartTime     time.Time                `json:"startTime,omitempty"`
	EndTime       time.Time                `json:"endTime,omitempty"`
	Filters       types.CostFilters        `json:"filters,omitempty"`
	SharedOptions *types.SharedCostOptions `json:"sharedOptions,omitempty"`
}

func (r *SummaryQueriedCostRequest) Validate() error {
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

type SummaryQueriedCostResponse struct {
	Cost `json:",inline"`

	AverageDailyCost   float64             `json:"averageDailyCost,omitempty"`
	ConnectorCount     int                 `json:"connectorCount,omitempty"`
	CollectedTimeRange *CollectedTimeRange `json:"collectedTimeRange,omitempty"`
}
