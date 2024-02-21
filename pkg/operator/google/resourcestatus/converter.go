package resourcestatus

import "github.com/seal-io/walrus/pkg/dao/types/status"

// sqlDatabaseInstanceStatusConverter is a converter for SQL Database Instance status
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | runnable              |                       |
// | failed                | Error                 |
// | suspended             | Inactive              |
// ref: https://github.com/googleapis/googleapis/blob/master/google/cloud/sql/v1/cloud_sql_instances.proto
var sqlDatabaseInstanceStatusConverter = status.NewConverter(
	[]string{
		"runnable",
	},
	[]string{
		"suspended",
	},
	[]string{
		"failed",
	},
)

// redisInstanceStatusConverter is a converter for Redis Instance status
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | ready                 |                       |
// | failing_over          | Error                 |
// ref: https://github.com/googleapis/googleapis/blob/master/google/cloud/redis/v1/cloud_redis.proto
var redisInstanceStatusConverter = status.NewConverter(
	[]string{
		"ready",
	},
	nil,
	[]string{
		"failing_over",
	},
)

// computeInstanceStatusConverter is a converter for Compute Instance status
//
// | Human Readable Status | Human Sensible Status |
// | --------------------- | --------------------- |
// | running               |                       |
// | terminated            | Error                 |
// | stopped               | Inactive              |
// ref: https://cloud.google.com/compute/docs/reference/rest/v1/instances
var computeInstanceStatusConverter = status.NewConverter(
	[]string{
		"running",
	},
	[]string{
		"stopped",
	},
	[]string{
		"terminated",
	},
)
