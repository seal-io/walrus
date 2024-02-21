package resourcestatus

import (
	"github.com/seal-io/walrus/pkg/dao/types/status"
)

// ec2InstanceStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | running               |                       |
// | terminated            | Inactive              |
// | stopped               | Inactive              |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
var ec2InstanceStatusConverter = status.NewConverter(
	[]string{
		"running",
	},
	[]string{
		"terminated",
		"stopped",
	},
	nil,
)

// ec2ImageStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | available             |                       |
// | failed                | Error                 |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeImages.html
var ec2ImageStatusConverter = status.NewConverter(
	[]string{
		"available",
	},
	nil,
	[]string{
		"failed",
	},
)

// ec2VolumeStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | in_use                |                       |
// | available             |                       |
// | deleted               | Inactive              |
// | error                 |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeVolumes.html
var ec2VolumeStatusConverter = status.NewConverter(
	[]string{
		"in_use",
		"available",
	},
	[]string{
		"deleted",
	},
	[]string{
		"error",
	},
)

// ec2SnapshotStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | completed             |                       |
// | error                 | Error                 |
// | recoverable           |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Snapshot.html
var ec2SnapshotStatusConverter = status.NewConverter(
	[]string{
		"completed",
		"recoverable",
	},
	nil,
	[]string{
		"error",
	},
)

// ec2NetworkInterfaceStatusConverter generate the summary use following table,
// other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | available             |                       |
// | associated            |                       |
// | in-use                |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceNetworkInterface.html
var ec2NetworkInterfaceStatusConverter = status.NewConverter(
	[]string{
		"available",
		"associated",
		"in-use",
	},
	nil,
	nil,
)

// vpcStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | available             |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Vpc.html
var vpcStatusConverter = status.NewConverter(
	[]string{
		"available",
	},
	nil,
	nil,
)

// ec2SubnetStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | available             |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Subnet.html
var ec2SubnetStatusConverter = status.NewConverter(
	[]string{
		"available",
	},
	nil,
	nil,
)

// rdsDBInstanceStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status                           | Human Sensible Status |
// | ----------------------------------------------- | --------------------- |
// | Available                                       |                       |
// | Failed                                          | Error                 |
// | Inaccessible-encryption-credentials             | Error                 |
// | Inaccessible-encryption-credentials-recoverable | Error                 |
// | Incompatible-network                            | Error                 |
// | Incompatible-option-group                       | Error                 |
// | Incompatible-parameters                         | Error                 |
// | Incompatible-restore                            | Error                 |
// | Insufficient-capacity                           | Error                 |
// | Restore-error                                   | Error                 |
// | Stopped                                         | Inactive              |
// | Storage-full                                    | Error                 |
// ref: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/accessing-monitoring.html#Overview.DBInstance.Status
var rdsDBInstanceStatusConverter = status.NewConverter(
	[]string{
		"available",
	},
	[]string{
		"stopped",
	},
	[]string{
		"failed",
		"inaccessible-encryption-credentials",
		"inaccessible-encryption-credentials-recoverable",
		"incompatible-network",
		"incompatible-option-group",
		"incompatible-parameters",
		"incompatible-restore",
		"insufficient-capacity",
		"restore-error",
		"storage-full",
	},
)

// rdsDBClusterStatusConverter generate the summary use following table, other status will be treated as transitioning.
// | Human Readable Status                           | Human Sensible Status |
// | ----------------------------------------------- | --------------------- |
// | Available                                       |                       |
// | Backing-up                                      |                       |
// | Cloning-failed                                  | Error                 |
// | Failing-over                                    | Error                 |
// | Inaccessible-encryption-credentials             | Error                 |
// | Inaccessible-encryption-credentials-recoverable | Error                 |
// | Maintenance                                     |                       |
// | Migration-failed                                | Error                 |
// | Stopped                                         | Inactive              |
// https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/accessing-monitoring.html#Aurora.Status
var rdsDBClusterStatusConverter = status.NewConverter(
	[]string{
		"available",
		"maintenance",
	},
	[]string{
		"stopped",
	},
	[]string{
		"cloning-failed",
		"failing-over",
		"inaccessible-encryption-credentials",
		"inaccessible-encryption-credentials-recoverable",
		"migration-failed",
	},
)

// cloudFrontStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Deployed              |                       |
// ref: https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/distribution-web-values-returned.html
var cloudFrontStatusConverter = status.NewConverter(
	[]string{
		"deployed",
	},
	nil,
	nil,
)

// elasticCacheStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status   | Human Sensible Status |
// | ----------------------- | --------------------- |
// | available               |                       |
// | deleted                 | Inactive              |
// | incompatible-network    | Error                 |
// | restore-failed          | Error                 |
// ref: https://docs.aws.amazon.com/AmazonElastiCache/latest/APIReference/API_CacheCluster.html
var elasticCacheStatusConverter = status.NewConverter(
	[]string{
		"available",
	},
	[]string{
		"deleted",
	},
	[]string{
		"incompatible-network",
		"restore-failed",
	},
)

// elbLoadBalancerStatusConverter generate the summary use following table,
// other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | active                |                       |
// | active_impaired       | Error                 |
// | failed                | Error                 |
// ref: https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_LoadBalancerState.html
var elbLoadBalancerStatusConverter = status.NewConverter(
	[]string{
		"active",
	},
	nil,
	[]string{
		"failed",
		"active_impaired",
	},
)

// eksClusterStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | ACTIVE                |                       |
// | FAILED                | Error                 |
// | INACTIVE              | Inactive              |
// ref: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Cluster.html
var eksClusterStatusConverter = status.NewConverter(
	[]string{
		"ACTIVE",
	},
	[]string{
		"INACTIVE",
	},
	[]string{
		"FAILED",
	},
)
