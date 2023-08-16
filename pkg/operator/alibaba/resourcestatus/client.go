package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/cs"

	"github.com/seal-io/walrus/pkg/operator/types"
)

func ecsClient(cred types.Credential) (*ecs.Client, error) {
	cli, err := ecs.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("error create alibaba ecs client %s: %w", cred.AccessKey, err)
	}

	return cli, nil
}

func cdnClient(cred types.Credential) (*cdn.Client, error) {
	cli, err := cdn.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("error create alibaba cdn client %s: %w", cred.AccessKey, err)
	}

	return cli, nil
}

func rdsClient(cred types.Credential) (*rds.Client, error) {
	cli, err := rds.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("error create alibaba rds client %s: %w", cred.AccessKey, err)
	}

	return cli, nil
}

func polarDBClient(cred types.Credential) (*polardb.Client, error) {
	cli, err := polardb.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("error create alibaba polar db client %s: %w", cred.AccessKey, err)
	}

	return cli, nil
}

func slbClient(cred types.Credential) (*slb.Client, error) {
	cli, err := slb.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("error create alibaba slb client %s: %w", cred.AccessKey, err)
	}

	return cli, nil
}

func vpcClient(cred types.Credential) (*vpc.Client, error) {
	cli, err := vpc.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("error create alibaba vpc client %s: %w", cred.AccessKey, err)
	}

	return cli, nil
}

func csClient(cred types.Credential) *cs.Client {
	return cs.NewClient(cred.AccessKey, cred.AccessSecret)
}
