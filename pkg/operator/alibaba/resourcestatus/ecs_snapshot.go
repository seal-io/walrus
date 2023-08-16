package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getEcsSnapshot(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := ecsClient(cred)
	if err != nil {
		return nil, err
	}

	req := ecs.CreateDescribeSnapshotsRequest()
	req.Scheme = schemeHttps
	req.SnapshotIds = toReqIds(name)

	resp, err := cli.DescribeSnapshots(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Snapshots.Snapshot) == 0 {
		return nil, errNotFound
	}

	st := ecsSnapshotStatusConverter.Convert(resp.Snapshots.Snapshot[0].Status, "")

	return st, nil
}
