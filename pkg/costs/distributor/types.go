package distributor

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
)

type SharedCost struct {
	StartTime      time.Time        `json:"startTime"`
	TotalCost      float64          `json:"totalCost"`
	IdleCost       float64          `json:"idleCost"`
	ManagementCost float64          `json:"managementCost"`
	AllocationCost float64          `json:"allocationCost"`
	Condition      types.SharedCost `json:"condition"`
}
