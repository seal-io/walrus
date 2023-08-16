package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getEcsNetworkInterface(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := ecsClient(cred)
	if err != nil {
		return nil, err
	}

	req := ecs.CreateDescribeNetworkInterfacesRequest()
	req.Scheme = schemeHttps
	req.NetworkInterfaceId = &[]string{name}

	resp, err := cli.DescribeNetworkInterfaces(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.NetworkInterfaceSets.NetworkInterfaceSet) == 0 {
		return nil, errNotFound
	}

	st := ecsNetworkInterfaceStatusConverter.Convert(resp.NetworkInterfaceSets.NetworkInterfaceSet[0].Status, "")

	return st, nil
}
