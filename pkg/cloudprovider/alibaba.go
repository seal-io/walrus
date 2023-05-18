package cloudprovider

import (
	"errors"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/cs"
)

const (
	schemeHttps = "HTTPS"
)

var errNotFound = errors.New("not found")

func NewAlibaba(name string, cred *Credential) (*Alibaba, error) {
	var (
		err     error
		alibaba = &Alibaba{
			name:       name,
			Credential: *cred,
		}
	)

	// ECS.
	alibaba.ecsCli, err = ecs.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error alibaba ecs client %s: %w", cred.AccessKey, err)
	}

	// CDN.
	alibaba.cdnCli, err = cdn.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error alibaba cdn client %s: %w", cred.AccessKey, err)
	}

	// RDS.
	alibaba.rdsCli, err = rds.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error alibaba rds client %s: %w", cred.AccessKey, err)
	}

	alibaba.polarDBCli, err = polardb.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error alibaba polar db client %s: %w", cred.AccessKey, err)
	}

	// SLB.
	alibaba.slbCli, err = slb.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error alibaba slb client %s: %w", cred.AccessKey, err)
	}

	// VPC.
	alibaba.vpcCli, err = vpc.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error alibaba vpc client %s: %w", cred.AccessKey, err)
	}

	// Kubernetes.
	alibaba.csCli = cs.NewClient(cred.AccessKey, cred.SecretKey)

	return alibaba, nil
}

type Alibaba struct {
	Credential

	name       string
	ecsCli     *ecs.Client
	cdnCli     *cdn.Client
	rdsCli     *rds.Client
	polarDBCli *polardb.Client
	slbCli     *slb.Client
	vpcCli     *vpc.Client
	csCli      *cs.Client
}

func (a *Alibaba) IsConnected() error {
	client, err := ram.NewClientWithAccessKey(a.Region, a.AccessKey, a.SecretKey)
	if err != nil {
		return err
	}

	// Use ListUsers API to check reachable.
	req := ram.CreateListUsersRequest()
	req.Scheme = schemeHttps

	_, err = client.ListUsers(req)
	if err != nil {
		return fmt.Errorf("error connect to alibaba %s: %w", a.name, err)
	}

	return nil
}

func (a *Alibaba) DescribeEcsInstance(id string) (*ecs.Instance, error) {
	req := ecs.CreateDescribeInstancesRequest()
	req.Scheme = schemeHttps
	req.InstanceIds = a.toReqIds(id)

	resp, err := a.ecsCli.DescribeInstances(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s ecs instance %s: %w", a.name, id, err)
	}

	if len(resp.Instances.Instance) == 0 {
		return nil, errNotFound
	}

	return &resp.Instances.Instance[0], nil
}

func (a *Alibaba) DescribeEcsImage(id string) (*ecs.Image, error) {
	req := ecs.CreateDescribeImagesRequest()
	req.Scheme = schemeHttps
	req.ImageId = id

	resp, err := a.ecsCli.DescribeImages(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s ecs image %s: %w", a.name, id, err)
	}

	if len(resp.Images.Image) == 0 {
		return nil, errNotFound
	}

	return &resp.Images.Image[0], nil
}

func (a *Alibaba) DescribeEcsDisk(id string) (*ecs.Disk, error) {
	req := ecs.CreateDescribeDisksRequest()
	req.Scheme = schemeHttps
	req.DiskIds = a.toReqIds(id)

	resp, err := a.ecsCli.DescribeDisks(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s ecs disk %s: %w", a.name, id, err)
	}

	if len(resp.Disks.Disk) == 0 {
		return nil, errNotFound
	}

	return &resp.Disks.Disk[0], nil
}

func (a *Alibaba) DescribeEcsSnapshot(id string) (*ecs.Snapshot, error) {
	req := ecs.CreateDescribeSnapshotsRequest()
	req.Scheme = schemeHttps
	req.SnapshotIds = a.toReqIds(id)

	resp, err := a.ecsCli.DescribeSnapshots(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s ecs snapshot %s: %w", a.name, id, err)
	}

	if len(resp.Snapshots.Snapshot) == 0 {
		return nil, errNotFound
	}

	return &resp.Snapshots.Snapshot[0], nil
}

func (a *Alibaba) DescribeEcsNetworkInterface(id string) (*ecs.NetworkInterfaceSet, error) {
	req := ecs.CreateDescribeNetworkInterfacesRequest()
	req.Scheme = schemeHttps
	req.NetworkInterfaceId = &[]string{id}

	resp, err := a.ecsCli.DescribeNetworkInterfaces(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s ecs network interface %s: %w", a.name, id, err)
	}

	if len(resp.NetworkInterfaceSets.NetworkInterfaceSet) == 0 {
		return nil, errNotFound
	}

	return &resp.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
}

func (a *Alibaba) DescribeCDNDomain(id string) (*cdn.GetDomainDetailModel, error) {
	req := cdn.CreateDescribeCdnDomainDetailRequest()
	req.Scheme = schemeHttps
	req.DomainName = id

	resp, err := a.cdnCli.DescribeCdnDomainDetail(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s cdn domain %s: %w", a.name, id, err)
	}

	return &resp.GetDomainDetailModel, nil
}

func (a *Alibaba) DescribeRDSDBInstance(id string) (*rds.DBInstanceAttribute, error) {
	req := rds.CreateDescribeDBInstanceAttributeRequest()
	req.Scheme = schemeHttps
	req.DBInstanceId = id

	resp, err := a.rdsCli.DescribeDBInstanceAttribute(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s rds db instance %s: %w", a.name, id, err)
	}

	if len(resp.Items.DBInstanceAttribute) == 0 {
		return nil, errNotFound
	}

	return &resp.Items.DBInstanceAttribute[0], nil
}

func (a *Alibaba) DescribePolarDBCluster(id string) (*polardb.DBCluster, error) {
	req := polardb.CreateDescribeDBClustersRequest()
	req.Scheme = schemeHttps
	req.DBClusterIds = a.toReqIds(id)

	resp, err := a.polarDBCli.DescribeDBClusters(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s polar db instance %s: %w", a.name, id, err)
	}

	if len(resp.Items.DBCluster) == 0 {
		return nil, errNotFound
	}

	return &resp.Items.DBCluster[0], nil
}

func (a *Alibaba) DescribeVpc(id string) (*vpc.Vpc, error) {
	req := vpc.CreateDescribeVpcsRequest()
	req.Scheme = schemeHttps
	req.VpcId = id

	resp, err := a.vpcCli.DescribeVpcs(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s vpc %s: %w", a.name, id, err)
	}

	if len(resp.Vpcs.Vpc) == 0 {
		return nil, errNotFound
	}

	return &resp.Vpcs.Vpc[0], nil
}

func (a *Alibaba) DescribeVSwitch(id string) (*vpc.VSwitch, error) {
	req := vpc.CreateDescribeVSwitchesRequest()
	req.Scheme = schemeHttps
	req.VSwitchId = id

	resp, err := a.vpcCli.DescribeVSwitches(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s vSwitch %s: %w", a.name, id, err)
	}

	if len(resp.VSwitches.VSwitch) == 0 {
		return nil, errNotFound
	}

	return &resp.VSwitches.VSwitch[0], nil
}

func (a *Alibaba) DescribeEipAddress(id string) (*vpc.EipAddress, error) {
	req := vpc.CreateDescribeEipAddressesRequest()
	req.Scheme = schemeHttps
	req.AllocationId = id

	resp, err := a.vpcCli.DescribeEipAddresses(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s vSwitch %s: %w", a.name, id, err)
	}

	if len(resp.EipAddresses.EipAddress) == 0 {
		return nil, errNotFound
	}

	return &resp.EipAddresses.EipAddress[0], nil
}

func (a *Alibaba) DescribeSLBLoadBalancer(id string) (*slb.LoadBalancer, error) {
	req := slb.CreateDescribeLoadBalancersRequest()
	req.Scheme = schemeHttps
	req.LoadBalancerId = id

	resp, err := a.slbCli.DescribeLoadBalancers(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s rds db instance %s: %w", a.name, id, err)
	}

	if len(resp.LoadBalancers.LoadBalancer) == 0 {
		return nil, errNotFound
	}

	return &resp.LoadBalancers.LoadBalancer[0], nil
}

func (a *Alibaba) DescribeCSCluster(id string) (*cs.ClusterType, error) {
	resp, err := a.csCli.DescribeCluster(id)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba %s container cluster %s: %w", a.name, id, err)
	}

	return &resp, nil
}

func (a *Alibaba) toReqIds(id string) string {
	return fmt.Sprintf(`["%s"]`, id)
}
