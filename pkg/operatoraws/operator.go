package operatoraws

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/cloudprovider"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform/operator"
)

const OperatorType = "AWS"

// New returns operator.Operator with the given options.
func New(ctx context.Context, opts operator.CreateOptions) (operator.Operator, error) {
	cred, err := cloudprovider.CredentialFromConnector(&opts.Connector)
	if err != nil {
		return nil, err
	}

	cli, err := cloudprovider.NewAWS(opts.Connector.ID.String(), cred)
	if err != nil {
		return nil, err
	}

	return Operator{
		cli: cli,
	}, nil
}

type Operator struct {
	cli *cloudprovider.AWS
}

func (o Operator) Type() operator.Type {
	return OperatorType
}

func (o Operator) IsConnected(_ context.Context) error {
	return o.cli.IsConnected()
}

func (o Operator) GetStatus(ctx context.Context, resource *model.ApplicationResource) (*status.Status, error) {
	var (
		st = &status.Status{}
		sm = &status.Summary{}
	)
	// EC2.
	switch resource.Type {
	case "aws_instance":
		res, err := o.cli.DescribeEc2Instance(ctx, resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.State.Name)
		sm = ec2InstanceStatusPaths.Walk(st)
	case "aws_ami":
		res, err := o.cli.DescribeEc2Image(ctx, resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.State)
		sm = ec2ImageStatusPaths.Walk(st)
	case "aws_ebs_volume":
		res, err := o.cli.DescribeEc2Volume(ctx, resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.State)
		sm = ec2VolumeStatusPaths.Walk(st)
	case "aws_ebs_snapshot":
		res, err := o.cli.DescribeEc2Snapshot(ctx, resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.State)
		sm = ec2SnapshotStatusPaths.Walk(st)
	case "aws_network_interface":
		res, err := o.cli.DescribeEc2NetworkInterface(ctx, resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.Status)
		sm = ec2NetworkInterfaceStatusPaths.Walk(st)
	case "aws_vpc":
		res, err := o.cli.DescribeVpc(ctx, resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.State)
		sm = vpcStatusPaths.Walk(st)
	case "aws_subnet":
		res, err := o.cli.DescribeSubnet(ctx, resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.State)
		sm = subnetStatusPaths.Walk(st)
	// RDS.
	case "aws_db_instance":
		res, err := o.cli.DescribeRDSDBInstance(ctx, resource.Name)
		if err != nil {
			return st, err
		}

		if res.DBInstanceStatus != nil {
			st.SummaryStatus = *res.DBInstanceStatus
			sm = rdsDBInstanceStatusPaths.Walk(st)
		}
	case "alicloud_db_instance":
		res, err := o.cli.DescribeRDSDBCluster(ctx, resource.Name)
		if err != nil {
			return st, err
		}

		if res.Status != nil {
			st.SummaryStatus = *res.Status
			sm = rdsDBClusterStatusPaths.Walk(st)
		}
	// Cloud Front.
	case "aws_cloudfront_distribution":
		res, err := o.cli.DescribeCloudFrontDistribute(ctx, resource.Name)
		if err != nil {
			return st, err
		}

		if res.Status != nil {
			st.SummaryStatus = *res.Status
			sm = cloudFrontStatusPaths.Walk(st)
		}
	// Cache.
	case "aws_elasticache_cluster":
		res, err := o.cli.DescribeElasticCache(ctx, resource.Name)
		if err != nil {
			return st, err
		}

		if res.CacheClusterStatus != nil {
			st.SummaryStatus = *res.CacheClusterStatus
			sm = elasticCacheStatusPaths.Walk(st)
		}
	// LB.
	case "aws_lb":
		res, err := o.cli.DescribeLoadBalancer(ctx, resource.Name)
		if err != nil {
			return st, err
		}

		if res.State != nil {
			st.SummaryStatus = res.State.String()
			sm = elbLoadBalancerStatusPaths.Walk(st)
		}
	// Kubernetes.
	case "aws_eks_cluster":
		res, err := o.cli.DescribeEKSCluster(ctx, resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.Status)
		sm = eksClusterStatusPaths.Walk(st)
	}

	st.Summary = *sm

	return st, nil
}

func (o Operator) GetKeys(ctx context.Context, resource *model.ApplicationResource) (*operator.Keys, error) {
	return nil, nil
}

func (o Operator) GetEndpoints(
	ctx context.Context,
	resource *model.ApplicationResource,
) ([]types.ApplicationResourceEndpoint, error) {
	return nil, nil
}

func (o Operator) GetComponents(
	ctx context.Context,
	resource *model.ApplicationResource,
) ([]*model.ApplicationResource, error) {
	return nil, nil
}

func (o Operator) Log(ctx context.Context, s string, options operator.LogOptions) error {
	return errors.New("cannot log")
}

func (o Operator) Exec(ctx context.Context, s string, options operator.ExecOptions) error {
	return errors.New("cannot execute")
}

func (o Operator) Label(ctx context.Context, resource *model.ApplicationResource, m map[string]string) error {
	return nil
}
