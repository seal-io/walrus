package status

const (
	ConnectorStatusProvisioned       ConditionType = "Provisioned"
	ConnectorStatusCostToolsDeployed ConditionType = "CostToolDeployed"
	ConnectorStatusCostSynced        ConditionType = "CostSynced"
	ConnectorStatusReady             ConditionType = "Ready"
)

// connectorStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Provisioned      | Unknown                 | Provisioning          | Transitioning         |
//	| Provisioned      | False                   | ProvisionFailed       | Error                 |
//	| Provisioned      | True                    | Provisioned           |                       |
//	| CostToolDeployed | Unknown                 | CostToolDeploying     | Transitioning         |
//	| CostToolDeployed | False                   | CostToolDeployFailed  | Error                 |
//	| CostToolDeployed | True                    | CostToolDeployed      |                       |
//	| CostSynced       | Unknown                 | CostSyncing           | Transitioning         |
//	| CostSynced       | False                   | CostSyncFailed        | Error                 |
//	| CostSynced       | True                    | CostSynced            |                       |
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | Unready               | Error                 |
//	| Ready            | True                    | Ready                 |                       |
var connectorStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ConnectorStatusProvisioned,
			ConnectorStatusCostToolsDeployed,
			ConnectorStatusCostSynced,
			ConnectorStatusReady,
		},
	},
)

// WalkConnector walks the given status by connector flow.
func WalkConnector(st *Status) *Summary {
	return connectorStatusPaths.Walk(st)
}
