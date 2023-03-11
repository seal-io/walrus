package status

const (
	ConnectorStatusInitializing = "Initializing"
	ConnectorStatusDeploying    = "Deploying"
	ConnectorStatusReady        = "Ready"
	ConnectorStatusError        = "Error"

	ConnectorFinOpsSyncStatusWaiting = "Waiting"
	ConnectorFinOpsSyncStatusSynced  = "Synced"
	ConnectorFinOpsSyncStatusError   = "Error"
)

const (
	ModuleStatusInitializing = "Initializing"
	ModuleStatusReady        = "Ready"
	ModuleStatusError        = "Error"
)

const (
	ApplicationInstanceStatusDeploying    = "Deploying"
	ApplicationInstanceStatusDeployed     = "Deployed"
	ApplicationInstanceStatusDeployFailed = "DeployFailed"
	ApplicationInstanceStatusDeleting     = "Deleting"
	ApplicationInstanceStatusDeleteFailed = "DeleteFailed"
)

const (
	ApplicationRevisionStatusRunning   = "Running"
	ApplicationRevisionStatusSucceeded = "Succeeded"
	ApplicationRevisionStatusFailed    = "Failed"
)
