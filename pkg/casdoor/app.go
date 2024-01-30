package casdoor

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/req"
)

type ApplicationCredential struct {
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

func GetApplicationCredential(
	ctx context.Context,
	adminSession []*req.HttpCookie,
	appz string,
) (*ApplicationCredential, error) {
	getApplicationURL := fmt.Sprintf("%s/api/get-application?id=admin/%s", endpoint.Get(), appz)

	var getAppCredResp struct {
		Status string                `json:"status"`
		Msg    string                `json:"msg"`
		Data   ApplicationCredential `json:"data"`
	}

	err := req.HTTPRequest().
		WithCookies(adminSession...).
		GetWithContext(ctx, getApplicationURL).
		BodyJSON(&getAppCredResp)
	if err != nil {
		return nil, errorx.Errorf("error getting app admin/%s: %v", appz, err)
	}

	if getAppCredResp.Status == statusError {
		return nil, errorx.Errorf("failed to get app credential: %s", getAppCredResp.Msg)
	}

	if getAppCredResp.Data.ClientID == "" || getAppCredResp.Data.ClientSecret == "" {
		return nil, errorx.Errorf("failed to get app admin/%s: blank client id/secret", appz)
	}

	return &getAppCredResp.Data, nil
}
