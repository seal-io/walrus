package resourcestatus

import (
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// ec2InstanceStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | running               |                       |
// | terminated            |                       |
// | stopped               |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
var ec2InstanceStatusConverter = status.NewConverter(
	[]string{
		"running",
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
// | deleted               |                       |
// | error                 |                       |
// ref: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeVolumes.html
var ec2VolumeStatusConverter = status.NewConverter(
	[]string{
		"in_use",
		"available",
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
)

// rdsDBInstanceStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status                           | Human Sensible Status |
// | ----------------------------------------------- | --------------------- |
// | available                                       |                       |
// | failed                                          | Error                 |
// | inaccessible-encryption-credentials             | Error                 |
// | inaccessible-encryption-credentials-recoverable | Error                 |
// | incompatible-network                            | Error                 |
// | incompatible-option-group                       | Error                 |
// | incompatible-parameters                         | Error                 |
// | incompatible-restore                            | Error                 |
// | insufficient-capacity                           | Error                 |
// | restore-error                                   | Error                 |
// | stopped                                         |                       |
// | storage-full                                    | Error                 |
// ref: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/accessing-monitoring.html#Overview.DBInstance.Status
// The status of above is capitalized, but it doesn't reflect the actual situation, see terraform ref instead,
//nolint: lll
// terraform ref: https://github.com/hashicorp/terraform-provider-aws/blob/8aea8f1e15c0fdb11036c816c2abc3afea451320/internal/service/rds/consts.go#L67-L98
var rdsDBInstanceStatusConverter = status.NewConverter(
	[]string{
		"available",
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
// | Stopped                                         |                       |
// https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/accessing-monitoring.html#Aurora.Status
var rdsDBClusterStatusConverter = status.NewConverter(
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
)

// cloudFrontStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Deployed              |                       |
// ref: https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/distribution-web-values-returned.html
var cloudFrontStatusConverter = status.NewConverter(
	[]string{
		"Deployed",
	},

	nil,
)

// elasticCacheStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status   | Human Sensible Status |
// | ----------------------- | --------------------- |
// | available               |                       |
// | deleted                 |                       |
// | incompatible-network    | Error                 |
// | restore-failed          | Error                 |
// ref: https://docs.aws.amazon.com/AmazonElastiCache/latest/APIReference/API_CacheCluster.html
var elasticCacheStatusConverter = status.NewConverter(
	[]string{
		"available",
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
// | INACTIVE              |                       |
// ref: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Cluster.html
var eksClusterStatusConverter = status.NewConverter(
	[]string{
		"ACTIVE",
		"INACTIVE",
	},

	[]string{
		"FAILED",
	},
)
