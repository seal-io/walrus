package operatoraws

import (
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// ec2InstanceStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | pending               | Transitioning         |
// | running               |                       |
// | shutting-down         | Transitioning         |
// | stopping              | Transitioning         |
// | terminated            |                       |
// | stopped               |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
var ec2InstanceStatusPaths = status.NewSummaryWalker(
	[]string{
		"running",
		"terminated",
		"stopped",
	},
	nil,
	[]string{
		"pending",
		"stopping",
		"shutting-down",
	},
)

// ec2ImageStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | pending               | Transitioning         |
// | available             |                       |
// | failed                | Error                 |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeImages.html
var ec2ImageStatusPaths = status.NewSummaryWalker(
	[]string{
		"available",
	},
	[]string{
		"failed",
	},
	[]string{
		"pending",
	},
)

// ec2VolumeStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | in_use                |                       |
// | available             |                       |
// | creating              | Transitioning         |
// | deleting              | Transitioning         |
// | deleted               |                       |
// | error                 |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeVolumes.html
var ec2VolumeStatusPaths = status.NewSummaryWalker(
	[]string{
		"in_use",
		"available",
		"deleted",
	},
	[]string{
		"error",
	},
	[]string{
		"creating",
		"deleting",
	},
)

// ec2SnapshotStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | pending               | Transitioning         |
// | completed             |                       |
// | error                 | Error                 |
// | recoverable           |                       |
// | recovering            | Transitioning         |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Snapshot.html
var ec2SnapshotStatusPaths = status.NewSummaryWalker(
	[]string{
		"completed",
		"recoverable",
	},
	[]string{
		"error",
	},
	[]string{
		"pending",
		"recovering",
	},
)

// ec2NetworkInterfaceStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | available             |                       |
// | attaching             | Transitioning         |
// | associated            |                       |
// | in-use                |                       |
// | detaching             | Transitioning         |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceNetworkInterface.html
var ec2NetworkInterfaceStatusPaths = status.NewSummaryWalker(
	[]string{
		"available",
		"associated",
		"in-use",
	},
	nil,
	[]string{
		"attaching",
		"detaching",
	},
)

// vpcStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | available             |                       |
// | pending               | Transitioning         |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Vpc.html
var vpcStatusPaths = status.NewSummaryWalker(
	[]string{
		"available",
	},
	nil,
	[]string{
		"pending",
	},
)

// subnetStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | available             |                       |
// | pending               | Transitioning         |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Subnet.html
var subnetStatusPaths = status.NewSummaryWalker(
	[]string{
		"available",
	},
	nil,
	[]string{
		"pending",
	},
)

// rdsDBInstanceStatusPaths generate the summary use following table.
//
// | Human Readable Status                           | Human Sensible Status |
// | ----------------------------------------------- | --------------------- |
// | Available                                       |                       |
// | Backing-up                                      | Transitioning         |
// | Configuring-enhanced-monitoring                 | Transitioning         |
// | Configuring-iam-database-auth                   | Transitioning         |
// | Configuring-log-exports                         | Transitioning         |
// | Converting-to-vpc                               | Transitioning         |
// | Creating                                        | Transitioning         |
// | Delete-precheck                                 | Transitioning         |
// | Deleting                                        | Transitioning         |
// | Failed                                          | Error                 |
// | Inaccessible-encryption-credentials             | Error                 |
// | Inaccessible-encryption-credentials-recoverable | Error                 |
// | Incompatible-network                            | Error                 |
// | Incompatible-option-group                       | Error                 |
// | Incompatible-parameters                         | Error                 |
// | Incompatible-restore                            | Error                 |
// | Insufficient-capacity                           | Error                 |
// | Moving-to-vpc                                   | Transitioning         |
// | Rebooting                                       | Transitioning         |
// | Resetting-master-credentials                    | Transitioning         |
// | Renaming                                        | Transitioning         |
// | Restore-error                                   | Error                 |
// | Starting                                        | Transitioning         |
// | Stopped                                         |                       |
// | Stopping                                        | Transitioning         |
// | Storage-full                                    | Error                 |
// | Storage-optimization                            | Transitioning         |
// | Upgrading                                       | Transitioning         |
// ref: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/accessing-monitoring.html#Overview.DBInstance.Status
var rdsDBInstanceStatusPaths = status.NewSummaryWalker(
	[]string{
		"Available",
		"Stopped",
	},

	[]string{
		"Failed",
		"Inaccessible-encryption-credentials",
		"Inaccessible-encryption-credentials-recoverable",
		"Incompatible-network",
		"Incompatible-option-group",
		"Incompatible-parameters",
		"Incompatible-restore",
		"Insufficient-capacity",
		"Restore-error",
		"Storage-full",
	},
	[]string{
		"Backing-up",
		"Configuring-enhanced-monitoring",
		"Configuring-iam-database-auth",
		"Configuring-log-exports",
		"Converting-to-vpc",
		"Creating",
		"Delete-precheck",
		"Deleting",
		"Moving-to-vpc",
		"Rebooting",
		"Resetting-master-credentials",
		"Renaming",
		"Starting",
		"Stopping",
		"Storage-optimization",
		"Upgrading",
	},
)

// rdsDBClusterStatusPaths generate the summary use following table.
// | Human Readable Status                           | Human Sensible Status |
// | ----------------------------------------------- | --------------------- |
// | Available                                       |                       |
// | Backing-up                                      |                       |
// | Backtracking                                    | Transitioning         |
// | Cloning-failed                                  | Error                 |
// | Creating                                        | Transitioning         |
// | Deleting                                        | Transitioning         |
// | Failing-over                                    | Error                 |
// | Inaccessible-encryption-credentials             | Error                 |
// | Inaccessible-encryption-credentials-recoverable | Error                 |
// | Maintenance                                     |                       |
// | Migrating                                       | Transitioning         |
// | Migration-failed                                | Error                 |
// | Modifying                                       | Transitioning         |
// | Promoting                                       | Transitioning         |
// | Renaming                                        | Transitioning         |
// | Resetting-master-credentials                    | Transitioning         |
// | Starting                                        | Transitioning         |
// | Stopped                                         |                       |
// | Stopping                                        | Transitioning         |
// | Storage-optimization                            | Transitioning         |
// | Update-iam-db-auth                              | Transitioning         |
// | Upgrading                                       | Transitioning         |
// https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/accessing-monitoring.html#Aurora.Status
var rdsDBClusterStatusPaths = status.NewSummaryWalker(
	[]string{
		"Available",
		"Maintenance",
		"Stopped",
	},

	[]string{
		"Cloning-failed",
		"Failing-over",
		"Inaccessible-encryption-credentials",
		"Inaccessible-encryption-credentials-recoverable",
		"Migration-failed",
	},
	[]string{
		"Backing-up",
		"Backtracking",
		"Creating",
		"Deleting",
		"Migrating",
		"Modifying",
		"Promoting",
		"Renaming",
		"Resetting-master-credentials",
		"Starting",
		"Stopping",
		"Storage-optimization",
		"Update-iam-db-auth",
		"Upgrading",
	},
)

// cloudFrontStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | InProgress            | Transitioning         |
// | Deployed              |                       |
// ref: https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/distribution-web-values-returned.html
var cloudFrontStatusPaths = status.NewSummaryWalker(
	[]string{
		"Deployed",
	},

	nil,
	[]string{
		"InProgress",
	},
)

// elasticCacheStatusPaths generate the summary use following table.
//
// | Human Readable Status   | Human Sensible Status |
// | ----------------------- | --------------------- |
// | available               |                       |
// | creating                | Transitioning         |
// | deleted                 |                       |
// | deleting                | Transitioning         |
// | incompatible-network    | Error                 |
// | modifying               | Transitioning         |
// | rebooting cluster nodes | Transitioning         |
// | restore-failed          | Error                 |
// | snapshotting            | Transitioning         |
// ref: https://docs.aws.amazon.com/AmazonElastiCache/latest/APIReference/API_CacheCluster.html
var elasticCacheStatusPaths = status.NewSummaryWalker(
	[]string{
		"available",
		"deleted",
	},

	[]string{
		"incompatible-network",
		"restore-failed",
	},
	[]string{
		"creating",
		"deleting",
		"modifying",
		"rebooting cluster nodes",
		"snapshotting",
	},
)

// elbLoadBalancerStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | active                |                       |
// | provisioning          | Transitioning         |
// | active_impaired       | Error                 |
// | failed                | Error                 |
// ref: https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_LoadBalancerState.html
var elbLoadBalancerStatusPaths = status.NewSummaryWalker(
	[]string{
		"active",
	},

	[]string{
		"failed",
		"active_impaired",
	},
	[]string{
		"provisioning",
	},
)

// eksClusterStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | ACTIVE                |                       |
// | PROVISIONING          | Transitioning         |
// | DEPROVISIONING        | Transitioning         |
// | FAILED                | Error                 |
// | INACTIVE              |                       |
// ref: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Cluster.html
var eksClusterStatusPaths = status.NewSummaryWalker(
	[]string{
		"ACTIVE",
		"INACTIVE",
	},

	[]string{
		"FAILED",
	},
	[]string{
		"PROVISIONING",
		"DEPROVISIONING",
	},
)
