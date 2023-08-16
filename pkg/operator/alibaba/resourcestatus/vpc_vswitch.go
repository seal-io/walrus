package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getVpcVSwitch(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := vpcClient(cred)
	if err != nil {
		return nil, err
	}

	req := vpc.CreateDescribeVSwitchesRequest()
	req.Scheme = schemeHttps
	req.VSwitchId = name

	resp, err := cli.DescribeVSwitches(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.VSwitches.VSwitch) == 0 {
		return nil, errNotFound
	}

	st := vpcVSwitchStatusConverter.Convert(resp.VSwitches.VSwitch[0].Status, "")

	return st, nil
}
