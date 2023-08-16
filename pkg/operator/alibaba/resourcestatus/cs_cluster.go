package resourcestatus

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getCsKubernetes(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli := csClient(cred)

	resp, err := cli.DescribeCluster(name)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	st := csClusterStatusConverter.Convert(string(resp.State), "")

	return st, nil
}
