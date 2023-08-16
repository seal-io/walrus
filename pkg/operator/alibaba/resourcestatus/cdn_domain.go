package resourcestatus

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getCdnDomain(cred types.Credential, resourceType, name string) (*status.Status, error) {
	cli, err := cdnClient(cred)
	if err != nil {
		return nil, err
	}

	req := cdn.CreateDescribeCdnDomainDetailRequest()
	req.Scheme = schemeHttps
	req.DomainName = name

	resp, err := cli.DescribeCdnDomainDetail(req)
	if err != nil {
		return nil, fmt.Errorf("error describe alibaba resource %s %s: %w", resourceType, name, err)
	}

	st := cdnDomainStatusConverter.Convert(resp.GetDomainDetailModel.DomainStatus, "")

	return st, nil
}
