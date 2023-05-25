package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/types"
)

func getPolarDBCluster(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := polarDBClient(cred)
	if err != nil {
		return nil, err
	}

	req := polardb.CreateDescribeDBClustersRequest()
	req.Scheme = schemeHttps
	req.DBClusterIds = toReqIds(name)

	resp, err := cli.DescribeDBClusters(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Items.DBCluster) == 0 {
		return nil, errNotFound
	}

	st := polarDBClusterStatusConverter.Convert(resp.Items.DBCluster[0].DBClusterStatus, "")

	return st, nil
}
