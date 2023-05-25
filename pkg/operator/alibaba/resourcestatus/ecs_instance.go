package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/types"
)

func getEcsInstance(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := ecsClient(cred)
	if err != nil {
		return nil, err
	}

	req := ecs.CreateDescribeInstancesRequest()
	req.Scheme = schemeHttps
	req.InstanceIds = toReqIds(name)

	resp, err := cli.DescribeInstances(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Instances.Instance) == 0 {
		return nil, errNotFound
	}

	st := ecsInstanceStatusConverter.Convert(resp.Instances.Instance[0].Status, "")

	return st, nil
}
