package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/types"
)

func getVpc(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := vpcClient(cred)
	if err != nil {
		return nil, err
	}

	req := vpc.CreateDescribeVpcsRequest()
	req.Scheme = schemeHttps
	req.VpcId = name

	resp, err := cli.DescribeVpcs(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Vpcs.Vpc) == 0 {
		return nil, errNotFound
	}

	st := vpcStatusConverter.Convert(resp.Vpcs.Vpc[0].Status, "")

	return st, nil
}
