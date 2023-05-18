package operatoralibaba

import (
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// ecsInstanceStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Pending               | Transitioning         |
// | Running               |                       |
// | Starting              | Transitioning         |
// | Stopping              | Transitioning         |
// | Stopped               |          			   |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describeinstances
var ecsInstanceStatusPaths = status.NewSummaryWalker(
	[]string{
		"Running",
		"Stopped",
	},
	nil,
	[]string{
		"Pending",
		"Starting",
		"Stopping",
	})

// ecsImageStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Creating              | Transitioning         |
// | Waiting               | Transitioning         |
// | Available             |                       |
// | UnAvailable           | Error                 |
// | CreateFailed          | Error                 |
// | Deprecated            | Error                 |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describeimages
var ecsImageStatusPaths = status.NewSummaryWalker(
	[]string{
		"Available",
	},

	[]string{
		"UnAvailable",
		"CreateFailed",
		"Deprecated",
	},
	[]string{
		"Creating",
		"Waiting",
	},
)

// ecsDiskStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | In_use                |                       |
// | Available             |                       |
// | Attaching             | Transitioning         |
// | Detaching             | Transitioning         |
// | ReIniting             | Transitioning         |
// | All                   |                       |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describedisks
var ecsDiskStatusPaths = status.NewSummaryWalker(
	[]string{
		"In_use",
		"Available",
		"All",
	},
	nil,
	[]string{
		"ReIniting",
		"Attaching",
		"Detaching",
	},
)

// ecsSnapshotStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | progressing           | Transitioning         |
// | accomplished          |                       |
// | failed                | Error                 |
// | all                   |                       |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describesnapshots
var ecsSnapshotStatusPaths = status.NewSummaryWalker(
	[]string{
		"accomplished",
		"all",
	},
	[]string{
		"failed",
	},
	[]string{
		"progressing",
	},
)

// ecsNetworkInterfaceStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Available             |                       |
// | Attaching             | Transitioning         |
// | InUse                 |                       |
// | Detaching             | Transitioning         |
// | Deleting              | Transitioning         |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describenetworkinterfaces
var ecsNetworkInterfaceStatusPaths = status.NewSummaryWalker(
	[]string{
		"Available",
		"InUse",
	},
	nil,
	[]string{
		"Attaching",
		"Detaching",
		"Deleting",
	},
)

// cdnDomainStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | online                |                       |
// | offline               | Error                 |
// | configuring           | Transitioning         |
// | configure_failed      | Error                 |
// | checking              | Transitioning         |
// | check_failed          | Error                 |
// | stopping              | Transitioning         |
// | deleting              | Transitioning         |
// https://www.alibabacloud.com/help/en/alibaba-cloud-cdn/latest/api-doc-cdn-2018-05-10-api-doc-describecdndomaindetail
var cdnDomainStatusPaths = status.NewSummaryWalker(
	[]string{
		"online",
	},
	[]string{
		"offline",
		"configure_failed",
		"check_failed",
	},
	[]string{
		"configuring",
		"checking",
		"stopping",
		"deleting",
	},
)

// rdsDBInstanceStatusPaths generate the summary use following table.
// | Human Readable Status     | Human Sensible Status |
// | ------------------------- | --------------------- |
// | Creating                  | Transitioning         |
// | Running                   |                       |
// | Deleting                  | Transitioning         |
// | Rebooting                 | Transitioning         |
// | DBInstanceClassChanging   | Transitioning         |
// | TRANSING                  | Transitioning         |
// | EngineVersionUpgrading    | Transitioning         |
// | TransingToOthers          | Transitioning         |
// | GuardDBInstanceCreating   | Transitioning         |
// | Restoring                 | Transitioning         |
// | Importing                 | Transitioning         |
// | ImportingFromOthers       | Transitioning         |
// | DBInstanceNetTypeChanging | Transitioning         |
// | GuardSwitching            | Transitioning         |
// | INS_CLONING               | Transitioning         |
// | Released                  |                       |
// ref: https://help.aliyun.com/document_detail/26315.htm?spm=a2c4g.610394.0.0.910d615eklhZvL
var rdsDBInstanceStatusPaths = status.NewSummaryWalker(
	[]string{
		"Running",
		"Released",
	},
	nil,
	[]string{
		"Creating",
		"Deleting",
		"Rebooting",
		"DBInstanceClassChanging",
		"TRANSING",
		"EngineVersionUpgrading",
		"TransingToOthers",
		"GuardDBInstanceCreating",
		"Restoring",
		"Importing",
		"ImportingFromOthers",
		"DBInstanceNetTypeChanging",
		"GuardSwitching",
		"INS_CLONING",
	},
)

// polarDBClusterStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Creating              | Transitioning         |
// | Running               |                       |
// | Deleting              | Transitioning         |
// | Deleted               |                       |
// | Rebooting             | Transitioning         |
// | Starting              | Transitioning         |
// | Stopped               |                       |
// | INS_MAINTAINING       | Transitioning         |
// | Switching             | Transitioning         |
// | DBNodeCreating        | Transitioning         |
// | DBNodeDeleting        | Transitioning         |
// | ClassChanging         | Transitioning         |
// | NetAddressCreating    | Transitioning         |
// | NetAddressDeleting    | Transitioning         |
// | NetAddressModifying   | Transitioning         |
// | MinorVersionUpgrading | Transitioning         |
// | STORAGE_EXPANDIN      | Transitioning         |
// | UPGRADE_FORBIDDEN     | Transitioning         |
// | TRANSING              | Transitioning         |
// | SSL_MODIFYING         | Transitioning         |
// | TDEModifying          | Transitioning         |
// | CONFIG_SWITCHING      | Transitioning         |
// | ROLE_SWITCHING        | Transitioning         |
// | ClassChanged          | Transitioning         |
// | MajorVersionUpgrading | Transitioning         |
// https://www.alibabacloud.com/help/en/polardb/latest/cluster-status
var polarDBClusterStatusPaths = status.NewSummaryWalker(
	[]string{
		"Running",
		"Deleted",
		"Stopped",
	},
	nil,
	[]string{
		"Creating",
		"Deleting",
		"Rebooting",
		"Starting",
		"INS_MAINTAINING",
		"Switching",
		"DBNodeCreating",
		"DBNodeDeleting",
		"ClassChanging",
		"NetAddressCreating",
		"NetAddressDeleting",
		"NetAddressModifying",
		"MinorVersionUpgrading",
		"STORAGE_EXPANDIN",
		"UPGRADE_FORBIDDEN",
		"TRANSING",
		"SSL_MODIFYING",
		"TDEModifying",
		"CONFIG_SWITCHING",
		"ROLE_SWITCHING",
		"ClassChanged",
		"MajorVersionUpgrading",
	},
)

// slbLoadBalancerStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | inactive              | Error                 |
// | active                |                       |
// | locked                | Error                 |
// ref: https://www.alibabacloud.com/help/en/server-load-balancer/latest/describeloadbalancers
var slbLoadBalancerStatusPaths = status.NewSummaryWalker(
	[]string{
		"active",
	},
	[]string{
		"inactive",
		"locked",
	},
	nil,
)

// vpcStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Available             |                       |
// | Pending               | Transitioning         |
// ref: https://next.api.aliyun.com/api/Vpc/2016-04-28/DescribeVpcs
var vpcStatusPaths = status.NewSummaryWalker(
	[]string{
		"Available",
	},
	nil,
	[]string{
		"Pending",
	},
)

// vSwitchStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Available             |                       |
// | Pending               | Transitioning         |
// ref: https://www.alibabacloud.com/help/en/ens/latest/describevswitches
var vSwitchStatusPaths = status.NewSummaryWalker(
	[]string{
		"Available",
	},
	nil,
	[]string{
		"Pending",
	},
)

// eipStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Associating           | Transitioning         |
// | Unassociating         | Transitioning         |
// | InUse                 |                       |
// | Available             |                       |
// | Releasing             | Transitioning         |
// ref: https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/describeeipaddresses
var eipStatusPaths = status.NewSummaryWalker(
	[]string{
		"InUse",
		"Available",
	},
	nil,
	[]string{
		"Associating",
		"Unassociating",
		"Releasing",
	},
)

// csClusterStatusPaths generate the summary use following table.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | initial               |                       |
// | failed                | Error                 |
// | running               |                       |
// | updating              | Transitioning         |
// | updating_failed       | Error                 |
// | scaling               | Transitioning         |
// | waiting               | Transitioning         |
// | disconnected          | Error                 |
// | stopped               |                       |
// | deleting              | Transitioning         |
// | deleted               |                       |
// | delete_failed         | Error                 |
// ref: https://www.alibabacloud.com/help/en/container-service-for-kubernetes/latest/describeclusterdetail
var csClusterStatusPaths = status.NewSummaryWalker(
	[]string{
		"initial",
		"running",
		"stopped",
		"deleted",
	},
	[]string{
		"failed",
		"updating_failed",
		"disconnected",
		"delete_failed",
	},
	[]string{
		"initial",
		"updating",
		"scaling",
		"waiting",
		"deleting",
	},
)
