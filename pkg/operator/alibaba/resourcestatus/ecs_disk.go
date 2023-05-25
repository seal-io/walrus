package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/types"
)

func getEcsDisk(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := ecsClient(cred)
	if err != nil {
		return nil, err
	}

	req := ecs.CreateDescribeDisksRequest()
	req.Scheme = schemeHttps
	req.DiskIds = toReqIds(name)

	resp, err := cli.DescribeDisks(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Disks.Disk) == 0 {
		return nil, errNotFound
	}

	st := ecsDiskStatusConverter.Convert(resp.Disks.Disk[0].Status, "")

	return st, nil
}
