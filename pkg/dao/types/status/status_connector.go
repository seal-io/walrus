package status

const (
	ConnectorStatusConnected         ConditionType = "Connected"
	ConnectorStatusCostToolsDeployed ConditionType = "CostToolDeployed"
	ConnectorStatusCostSynced        ConditionType = "CostSynced"
	ConnectorStatusReady             ConditionType = "Ready"
)

// connectorStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Connected        | Unknown                 | Connecting            | Transitioning         |
//	| Connected        | False                   | ConnectFailed         | Error                 |
//	| Connected        | True                    | Connected             |                       |
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
			ConnectorStatusConnected,
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
