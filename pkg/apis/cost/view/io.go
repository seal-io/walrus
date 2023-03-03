package view

import (
	"errors"
	"fmt"
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/slice"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs

// Batch APIs

// Extensional APIs

type (
	Resources []Resource
	Resource  struct {
		Cost

		ItemName  string     `json:"itemName,omitempty"`
		StartTime *time.Time `json:"startTime,omitempty"`
		EndTime   *time.Time `json:"endTime,omitempty"`
	}

	Cost struct {
		Currency   int     `json:"currency,omitempty"`
		TotalCost  float64 `json:"totalCost,omitempty"`
		SharedCost float64 `json:"sharedCost,omitempty"`
		CpuCost    float64 `json:"cpuCost,omitempty"`
		GpuCost    float64 `json:"gpuCost,omitempty"`
		RamCost    float64 `json:"ramCost,omitempty"`
		PvCost     float64 `json:"pvCost,omitempty"`
	}
)

type (
	AllocationCostRequest struct {
		_ struct{} `route:"POST=/allocation-costs"`

		types.QueryCondition `json:",inline"`

		StartTime time.Time `json:"startTime,omitempty"`
		EndTime   time.Time `json:"endTime,omitempty"`
	}
)

func (r *AllocationCostRequest) Validate() error {
	if err := validation.TimeRange(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if len(r.Filters) == 0 {
		return errors.New("invalid filter: blank")
	}

	if r.Step != "" && slice.ContainsAny([]types.GroupByField{
		types.GroupByFieldDay,
		types.GroupByFieldWeek,
		types.GroupByFieldMonth,
		types.GroupByFieldYear,
	}, r.GroupBy) {
		return fmt.Errorf("invalid step: already group by %s", r.GroupBy)
	}
	return nil
}

type SummaryCostRequest struct {
	_ struct{} `route:"POST=/summary-costs"`

	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
}

func (r *SummaryCostRequest) Validate() error {
	return validation.TimeRange(r.StartTime, r.EndTime)
}

type SummaryCostResponse struct {
	Currency         int     `json:"currency,omitempty"`
	TotalCost        float64 `json:"totalCost,omitempty"`
	AverageDailyCost float64 `json:"averageDailyCost,omitempty"`
	ProjectCount     int     `json:"projectCount,omitempty"`
	ClusterCount     int     `json:"clusterCount,omitempty"`
}

type SummaryClusterCostRequest struct {
	_ struct{} `route:"POST=/summary-cluster-costs"`

	StartTime   time.Time `json:"startTime,omitempty"`
	EndTime     time.Time `json:"endTime,omitempty"`
	ConnectorID types.ID  `json:"connectorID,omitempty"`
}

func (r *SummaryClusterCostRequest) Validate() error {
	if err := validation.TimeRange(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if !r.ConnectorID.IsNaive() {
		return errors.New("invalid connector id")
	}
	return nil
}

type SummaryClusterCostResponse struct {
	Currency         int     `json:"currency,omitempty"`
	TotalCost        float64 `json:"totalCost,omitempty"`
	AverageDailyCost float64 `json:"averageDailyCost,omitempty"`
	AllocationCost   float64 `json:"allocationCost,omitempty"`
	ManagementCost   float64 `json:"managementCost,omitempty"`
	IdleCost         float64 `json:"idleCost,omitempty"`
}

type SummaryProjectCostRequest struct {
	_ struct{} `route:"POST=/summary-project-costs"`

	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
	Project   string    `json:"project,omitempty"`
}

func (r *SummaryProjectCostRequest) Validate() error {
	if err := validation.TimeRange(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if r.Project == "" {
		return errors.New("invalid project: blank")
	}
	return nil
}

type SummaryCostCommonResponse struct {
	Currency         int     `json:"currency,omitempty"`
	TotalCost        float64 `json:"totalCost,omitempty"`
	AverageDailyCost float64 `json:"averageDailyCost,omitempty"`
}

type SummaryQueriedCostRequest struct {
	_ struct{} `route:"POST=/summary-queried-costs"`

	StartTime   time.Time                   `json:"startTime,omitempty"`
	EndTime     time.Time                   `json:"endTime,omitempty"`
	Filters     types.AllocationCostFilters `json:"filters,omitempty"`
	SharedCosts types.ShareCosts            `json:"shareCosts,omitempty"`
}

func (r *SummaryQueriedCostRequest) Validate() error {
	if err := validation.TimeRange(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if len(r.Filters) == 0 {
		return errors.New("invalid filters: blank")
	}
	return nil
}

type SummaryQueriedCostResponse struct {
	Cost `json:",inline"`

	AverageDailyCost float64 `json:"averageDailyCost,omitempty"`
	ConnectorCount   int     `json:"connectorCount,omitempty"`
}
