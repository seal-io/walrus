package cloudprovider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfrontypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	elacticachetypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func NewAWS(name string, cred *Credential) (*AWS, error) {
	cfg, err := newAWSConfig(cred)
	if err != nil {
		return nil, err
	}

	aws := &AWS{
		name:   name,
		config: cfg,
	}

	// EC2.
	aws.ec2Cli = ec2.NewFromConfig(*cfg)

	// RDS.
	aws.rdsCli = rds.NewFromConfig(*cfg)

	// CloudFront.
	aws.cloudfrontCli = cloudfront.NewFromConfig(*cfg)

	// Auto Scaling.
	aws.autoScaling = autoscaling.NewFromConfig(*cfg)

	// Cache.
	aws.elasticCacheCli = elasticache.NewFromConfig(*cfg)

	// ELB.
	sess, err := session.NewSession(&awsv1.Config{
		Region:      awsv1.String(cred.Region),
		Credentials: credentials.NewStaticCredentials(cred.AccessKey, cred.SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}
	aws.elbCli = elbv2.New(sess)

	// EKS.
	aws.eksCli = eks.NewFromConfig(*cfg)

	return aws, nil
}

type AWS struct {
	Credential

	name            string
	config          *aws.Config
	ec2Cli          *ec2.Client
	rdsCli          *rds.Client
	cloudfrontCli   *cloudfront.Client
	autoScaling     *autoscaling.Client
	elasticCacheCli *elasticache.Client
	elbCli          *elbv2.ELBV2
	eksCli          *eks.Client
}

func (a *AWS) IsConnected() error {
	// Use GetCallerIdentity API to check reachable.
	svc := sts.NewFromConfig(*a.config)

	_, err := svc.GetCallerIdentity(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("error connect to aws %s: %w", a.name, err)
	}

	return nil
}

func (a *AWS) DescribeEc2Instance(ctx context.Context, id string) (*ec2types.Instance, error) {
	resp, err := a.ec2Cli.DescribeInstances(
		ctx,
		&ec2.DescribeInstancesInput{
			InstanceIds: []string{id},
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s ec2 instance %s: %w", a.name, id, err)
	}

	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return nil, errNotFound
	}

	return &resp.Reservations[0].Instances[0], nil
}

func (a *AWS) DescribeEc2Image(ctx context.Context, id string) (*ec2types.Image, error) {
	resp, err := a.ec2Cli.DescribeImages(
		ctx,
		&ec2.DescribeImagesInput{
			ImageIds: []string{id},
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s ec2 image %s: %w", a.name, id, err)
	}

	if len(resp.Images) == 0 {
		return nil, errNotFound
	}

	return &resp.Images[0], nil
}

func (a *AWS) DescribeEc2Volume(ctx context.Context, id string) (*ec2types.Volume, error) {
	resp, err := a.ec2Cli.DescribeVolumes(
		ctx,
		&ec2.DescribeVolumesInput{
			VolumeIds: []string{id},
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s ec2 volume %s: %w", a.name, id, err)
	}

	if len(resp.Volumes) == 0 {
		return nil, errNotFound
	}

	return &resp.Volumes[0], nil
}

func (a *AWS) DescribeEc2Snapshot(ctx context.Context, id string) (*ec2types.Snapshot, error) {
	resp, err := a.ec2Cli.DescribeSnapshots(
		ctx,
		&ec2.DescribeSnapshotsInput{
			SnapshotIds: []string{id},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s ec2 snapshot %s: %w", a.name, id, err)
	}

	if len(resp.Snapshots) == 0 {
		return nil, errNotFound
	}

	return &resp.Snapshots[0], nil
}

func (a *AWS) DescribeEc2NetworkInterface(ctx context.Context, id string) (*ec2types.NetworkInterface, error) {
	resp, err := a.ec2Cli.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s ec2 network interface %s: %w", a.name, id, err)
	}

	if len(resp.NetworkInterfaces) == 0 {
		return nil, errNotFound
	}

	return &resp.NetworkInterfaces[0], nil
}

func (a *AWS) DescribeRDSDBInstance(ctx context.Context, id string) (*rdstypes.DBInstance, error) {
	resp, err := a.rdsCli.DescribeDBInstances(
		ctx,
		&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: &id,
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s rds db instance %s: %w", a.name, id, err)
	}

	if len(resp.DBInstances) == 0 {
		return nil, errNotFound
	}

	return &resp.DBInstances[0], nil
}

func (a *AWS) DescribeRDSDBCluster(ctx context.Context, id string) (*rdstypes.DBCluster, error) {
	resp, err := a.rdsCli.DescribeDBClusters(
		ctx,
		&rds.DescribeDBClustersInput{
			DBClusterIdentifier: &id,
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s rds db cluster %s: %w", a.name, id, err)
	}

	if len(resp.DBClusters) == 0 {
		return nil, errNotFound
	}

	return &resp.DBClusters[0], nil
}

func (a *AWS) DescribeVpc(ctx context.Context, id string) (*ec2types.Vpc, error) {
	resp, err := a.ec2Cli.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		VpcIds: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s vpc %s: %w", a.name, id, err)
	}

	if len(resp.Vpcs) == 0 {
		return nil, errNotFound
	}

	return &resp.Vpcs[0], nil
}

func (a *AWS) DescribeSubnet(ctx context.Context, id string) (*ec2types.Subnet, error) {
	resp, err := a.ec2Cli.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		SubnetIds: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s subnet %s: %w", a.name, id, err)
	}

	if len(resp.Subnets) == 0 {
		return nil, errNotFound
	}

	return &resp.Subnets[0], nil
}

func (a *AWS) DescribeCloudFrontDistribute(ctx context.Context, id string) (*cloudfrontypes.Distribution, error) {
	resp, err := a.cloudfrontCli.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: &id,
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s cloud front distribute %s: %w", a.name, id, err)
	}

	return resp.Distribution, nil
}

func (a *AWS) DescribeElasticCache(ctx context.Context, id string) (*elacticachetypes.CacheCluster, error) {
	resp, err := a.elasticCacheCli.DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{
		CacheClusterId: &id,
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s elastic cache %s: %w", a.name, id, err)
	}

	if len(resp.CacheClusters) == 0 {
		return nil, errNotFound
	}

	return &resp.CacheClusters[0], nil
}

func (a *AWS) DescribeLoadBalancer(_ context.Context, id string) (*elbv2.LoadBalancer, error) {
	resp, err := a.elbCli.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		Names: []*string{&id},
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s load balancer %s: %w", a.name, id, err)
	}

	if len(resp.LoadBalancers) == 0 {
		return nil, errNotFound
	}

	return resp.LoadBalancers[0], nil
}

func (a *AWS) DescribeEKSCluster(ctx context.Context, id string) (*ekstypes.Cluster, error) {
	resp, err := a.eksCli.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: &id,
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws %s eks cluster %s: %w", a.name, id, err)
	}

	return resp.Cluster, nil
}

// newAWSConfig create aws config from connector.
func newAWSConfig(cred *Credential) (*aws.Config, error) {
	ac := awsCredential{*cred}
	return ac.Config()
}

// awsCredential is a struct use to implement aws credentials.
type awsCredential struct {
	Credential
}

// Retrieve implement aws.CredentialsProvider interface.
func (a awsCredential) Retrieve(_ context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     a.AccessKey,
		SecretAccessKey: a.SecretKey,
	}, nil
}

// Config creates an AWS SDK V2 Config.
func (a awsCredential) Config() (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(a), config.WithRegion(a.Region))
	if err != nil {
		return &cfg,
			fmt.Errorf(
				"error initialize AWS SDK config with accessKey %s region %s: %w",
				a.AccessKey,
				a.Region,
				err,
			)
	}

	return &cfg, nil
}
