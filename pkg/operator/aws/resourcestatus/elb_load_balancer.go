package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getElbLoadBalancer(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := elbClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []*string{&name},
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.LoadBalancers) == 0 {
		return nil, errNotFound
	}

	lb := resp.LoadBalancers[0]
	if lb.State == nil || lb.State.Code == nil {
		return &status.Status{}, nil
	}

	var msg string
	if lb.State.Reason != nil {
		msg = *lb.State.Reason
	}

	st := elbLoadBalancerStatusConverter.Convert(*lb.State.Code, msg)

	return st, nil
}
