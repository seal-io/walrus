package resourcestatus

import "github.com/seal-io/walrus/pkg/dao/types/status"

// virtualMachineStatusConverter generate the summary use following table, other status is treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | running               |                       |
// | stopped               | Inactive              |
// | deallocated           | Inactive              |
// ref: https://learn.microsoft.com/en-us/azure/virtual-machines/states-billing#get-states-using-instance-view
var virtualMachineStatusConverter = status.NewConverter(
	[]string{
		"running",
	},
	[]string{
		"stopped",
		"deallocated",
	},
	nil,
)

// mySQLFlexibleServerStatusConverter generate the summary use following table, other status is treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Ready                 |                       |
// | Stopped               | Inactive              |
// | Disabled              | Inactive              |
// ref: https://learn.microsoft.com/en-us/rest/api/mysql/flexibleserver/servers/get
var mySQLFlexibleServerStatusConverter = status.NewConverter(
	[]string{
		"Ready",
	},
	[]string{
		"Stopped",
		"Disabled",
	},
	nil,
)

// postgreSQLFlexibleServerStatusConverter generate the summary use following table, other status is treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Ready                 |                       |
// | Stopped               | Inactive              |
// | Disabled              | Inactive              |
// ref: https://learn.microsoft.com/en-us/rest/api/postgresql/flexibleserver/servers/get
var postgreSQLFlexibleServerStatusConverter = status.NewConverter(
	[]string{
		"Ready",
	},
	[]string{
		"Stopped",
		"Disabled",
	},
	nil,
)

// redisCacheStatusConverter generate the summary use following table, other status is treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Succeeded             |                       |
// | Failed                | Error                 |
// | Disabled              | Inactive              |
// ref: https://learn.microsoft.com/en-us/rest/api/redis/redis/get
var redisCacheStatusConverter = status.NewConverter(
	[]string{
		"Succeeded",
	},
	[]string{
		"Disabled",
	},
	[]string{
		"Failed",
	},
)

// virtualNetworkStatusConverter generate the summary use following table, other status is treated as transitioning.
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | Succeeded             |                       |
// | Failed                | Error                 |
// ref: https://learn.microsoft.com/en-us/rest/api/virtualnetwork/virtual-networks/get
var virtualNetworkStatusConverter = status.NewConverter(
	[]string{
		"Succeeded",
	},
	nil,
	[]string{
		"Failed",
	},
)
