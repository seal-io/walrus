package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/types"
)

func getVpcEip(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := vpcClient(cred)
	if err != nil {
		return nil, err
	}

	req := vpc.CreateDescribeEipAddressesRequest()
	req.Scheme = schemeHttps
	req.AllocationId = name

	resp, err := cli.DescribeEipAddresses(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.EipAddresses.EipAddress) == 0 {
		return nil, errNotFound
	}

	st := vpcEipStatusConverter.Convert(resp.EipAddresses.EipAddress[0].Status, "")

	return st, nil
}
