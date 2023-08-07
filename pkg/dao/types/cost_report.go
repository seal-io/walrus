package types

const (
	// ManagementCostItemName indicates cost item name for cluster management.
	ManagementCostItemName = "__management__"

	// UnmountedCostItemName indicated cost item name for unmounted resource (PV or LB).
	UnmountedCostItemName = "__unmounted__"

	// IdleCostItemName indicated cost item name for idle.
	IdleCostItemName = "__idle__"

	// UnallocatedItemName indicate the cost for the resources unallocated.
	UnallocatedItemName = "__unallocated__"
)

type PVCost struct {
	Cost  float64 `json:"cost"`
	Bytes float64 `json:"bytes"`
}

func IsIdleCost(name string) bool {
	return name == IdleCostItemName
}

func IsManagementCost(name string) bool {
	return name == ManagementCostItemName
}

func IsIdleOrManagementCost(name string) bool {
	return IsIdleCost(name) || IsManagementCost(name)
}
