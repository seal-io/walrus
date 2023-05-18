package operatoralibaba

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/cloudprovider"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform/operator"
)

const OperatorType = "Alibaba"

// New returns operator.Operator with the given options.
func New(ctx context.Context, opts operator.CreateOptions) (operator.Operator, error) {
	cred, err := cloudprovider.CredentialFromConnector(&opts.Connector)
	if err != nil {
		return nil, err
	}

	cli, err := cloudprovider.NewAlibaba(opts.Connector.ID.String(), cred)
	if err != nil {
		return nil, err
	}

	return Operator{
		cli: cli,
	}, nil
}

type Operator struct {
	cli *cloudprovider.Alibaba
}

func (o Operator) Type() operator.Type {
	return OperatorType
}

func (o Operator) IsConnected(_ context.Context) error {
	return o.cli.IsConnected()
}

func (o Operator) GetStatus(_ context.Context, resource *model.ApplicationResource) (*status.Status, error) {
	var (
		st = &status.Status{}
		sm = &status.Summary{}
	)

	switch resource.Type {
	// ECS.
	case "alicloud_instance":
		res, err := o.cli.DescribeEcsInstance(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.Status
		sm = ecsInstanceStatusPaths.Walk(st)
	case "alicloud_image":
		res, err := o.cli.DescribeEcsImage(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.Status
		sm = ecsImageStatusPaths.Walk(st)
	case "alicloud_disk":
		res, err := o.cli.DescribeEcsDisk(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.Status
		sm = ecsDiskStatusPaths.Walk(st)
	case "alicloud_snapshot", "alicloud_ecs_snapshot":
		res, err := o.cli.DescribeEcsSnapshot(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.Status
		sm = ecsSnapshotStatusPaths.Walk(st)
	case "alicloud_network_interface", "alicloud_ecs_network_interface":
		res, err := o.cli.DescribeEcsNetworkInterface(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.Status
		sm = ecsNetworkInterfaceStatusPaths.Walk(st)
	// CDN.
	case "alicloud_cdn_domain":
		res, err := o.cli.DescribeCDNDomain(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.DomainStatus
		sm = cdnDomainStatusPaths.Walk(st)
	// RDS.
	case "alicloud_db_instance":
		res, err := o.cli.DescribeRDSDBInstance(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.DBInstanceStatus
		sm = rdsDBInstanceStatusPaths.Walk(st)

	case "alicloud_polardb_cluster":
		res, err := o.cli.DescribePolarDBCluster(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.DBClusterStatus
		sm = polarDBClusterStatusPaths.Walk(st)
	// SLB.
	case "alicloud_slb":
		res, err := o.cli.DescribeSLBLoadBalancer(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.LoadBalancerStatus
		sm = slbLoadBalancerStatusPaths.Walk(st)
	// Network.
	case "alicloud_vpc":
		res, err := o.cli.DescribeVpc(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.Status
		sm = vpcStatusPaths.Walk(st)
	case "alicloud_vswitch":
		res, err := o.cli.DescribeVSwitch(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.Status
		sm = vSwitchStatusPaths.Walk(st)
	case "alicloud_eip":
		res, err := o.cli.DescribeEipAddress(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = res.Status
		sm = eipStatusPaths.Walk(st)
	// Kubernetes.
	case "alicloud_cs_kubernetes":
		res, err := o.cli.DescribeCSCluster(resource.Name)
		if err != nil {
			return st, err
		}
		st.SummaryStatus = string(res.State)
		sm = csClusterStatusPaths.Walk(st)
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
