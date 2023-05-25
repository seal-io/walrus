package resourcestatus

import (
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// ecsInstanceStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Running               |                       |
// | Stopped               |          			   |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describeinstances
var ecsInstanceStatusConverter = status.NewConverter(
	[]string{
		"Running",
		"Stopped",
	},
	nil,
)

// ecsImageStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Available             |                       |
// | UnAvailable           | Error                 |
// | CreateFailed          | Error                 |
// | Deprecated            | Error                 |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describeimages
var ecsImageStatusConverter = status.NewConverter(
	[]string{
		"Available",
	},

	[]string{
		"UnAvailable",
		"CreateFailed",
		"Deprecated",
	},
)

// ecsDiskStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | In_use                |                       |
// | Available             |                       |
// | All                   |                       |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describedisks
var ecsDiskStatusConverter = status.NewConverter(
	[]string{
		"In_use",
		"Available",
		"All",
	},
	nil,
)

// ecsSnapshotStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | accomplished          |                       |
// | failed                | Error                 |
// | all                   |                       |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describesnapshots
var ecsSnapshotStatusConverter = status.NewConverter(
	[]string{
		"accomplished",
		"all",
	},
	[]string{
		"failed",
	},
)

// ecsNetworkInterfaceStatusConverter generate the summary use following table,
// other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Available             |                       |
// | InUse                 |                       |
// ref: https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describenetworkinterfaces
var ecsNetworkInterfaceStatusConverter = status.NewConverter(
	[]string{
		"Available",
		"InUse",
	},
	nil,
)

// cdnDomainStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | online                |                       |
// | offline               | Error                 |
// | configure_failed      | Error                 |
// | check_failed          | Error                 |
// https://www.alibabacloud.com/help/en/alibaba-cloud-cdn/latest/api-doc-cdn-2018-05-10-api-doc-describecdndomaindetail
var cdnDomainStatusConverter = status.NewConverter(
	[]string{
		"online",
	},
	[]string{
		"offline",
		"configure_failed",
		"check_failed",
	},
)

// rdsDBInstanceStatusConverter generate the summary use following table, other status will be treated as transitioning.
// | Human Readable Status     | Human Sensible Status |
// | ------------------------- | --------------------- |
// | Running                   |                       |
// | Released                  |                       |
// ref: https://help.aliyun.com/document_detail/26315.htm?spm=a2c4g.610394.0.0.910d615eklhZvL
var rdsDBInstanceStatusConverter = status.NewConverter(
	[]string{
		"Running",
		"Released",
	},
	nil,
)

// polarDBClusterStatusConverter generate the summary use following table,
// other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Running               |                       |
// | Deleted               |                       |
// | Stopped               |                       |
// https://www.alibabacloud.com/help/en/polardb/latest/cluster-status
var polarDBClusterStatusConverter = status.NewConverter(
	[]string{
		"Running",
		"Deleted",
		"Stopped",
	},
	nil,
)

// slbLoadBalancerStatusConverter generate the summary use following table,
// other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | inactive              | Error                 |
// | active                |                       |
// | locked                | Error                 |
// ref: https://www.alibabacloud.com/help/en/server-load-balancer/latest/describeloadbalancers
var slbLoadBalancerStatusConverter = status.NewConverter(
	[]string{
		"active",
	},
	[]string{
		"inactive",
		"locked",
	},
)

// vpcStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Available             |                       |
// ref: https://next.api.aliyun.com/api/Vpc/2016-04-28/DescribeVpcs
var vpcStatusConverter = status.NewConverter(
	[]string{
		"Available",
	},
	nil,
)

// vpcVSwitchStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Available             |                       |
// ref: https://www.alibabacloud.com/help/en/ens/latest/describevswitches
var vpcVSwitchStatusConverter = status.NewConverter(
	[]string{
		"Available",
	},
	nil,
)

// vpcEipStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | InUse                 |                       |
// | Available             |                       |
// ref: https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/describeeipaddresses
var vpcEipStatusConverter = status.NewConverter(
	[]string{
		"InUse",
		"Available",
	},
	nil,
)

// csClusterStatusConverter generate the summary use following table, other status will be treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | initial               |                       |
// | failed                | Error                 |
// | running               |                       |
// | updating_failed       | Error                 |
// | disconnected          | Error                 |
// | stopped               |                       |
// | deleted               |                       |
// | delete_failed         | Error                 |
// ref: https://www.alibabacloud.com/help/en/container-service-for-kubernetes/latest/describeclusterdetail
var csClusterStatusConverter = status.NewConverter(
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
)
