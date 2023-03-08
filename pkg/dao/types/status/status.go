package status

const (
	Initializing = "Initializing"
	Ready        = "Ready"
	Error        = "Error"
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
