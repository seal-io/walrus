package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"

	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/types"
)

func getSlbLoadBalancer(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := slbClient(cred)
	if err != nil {
		return nil, err
	}

	req := slb.CreateDescribeLoadBalancersRequest()
	req.Scheme = schemeHttps
	req.LoadBalancerId = name

	resp, err := cli.DescribeLoadBalancers(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.LoadBalancers.LoadBalancer) == 0 {
		return nil, errNotFound
	}

	st := slbLoadBalancerStatusConverter.Convert(resp.LoadBalancers.LoadBalancer[0].LoadBalancerStatus, "")

	return st, nil
}
