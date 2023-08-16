package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getRdsDBInstance(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := rdsClient(cred)
	if err != nil {
		return nil, err
	}

	req := rds.CreateDescribeDBInstanceAttributeRequest()
	req.Scheme = schemeHttps
	req.DBInstanceId = name

	resp, err := cli.DescribeDBInstanceAttribute(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Items.DBInstanceAttribute) == 0 {
		return nil, errNotFound
	}

	st := rdsDBInstanceStatusConverter.Convert(resp.Items.DBInstanceAttribute[0].DBInstanceStatus, "")

	return st, nil
}
