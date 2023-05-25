package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/types"
)

func getEcsImage(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := ecsClient(cred)
	if err != nil {
		return nil, err
	}

	req := ecs.CreateDescribeImagesRequest()
	req.Scheme = schemeHttps
	req.ImageId = name

	resp, err := cli.DescribeImages(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Images.Image) == 0 {
		return nil, errNotFound
	}

	st := ecsImageStatusConverter.Convert(resp.Images.Image[0].Status, "")

	return st, nil
}
