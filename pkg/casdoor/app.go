package casdoor

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/utils/req"
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
	var app ApplicationCredential
	err := req.HTTPRequest().
		WithCookies(adminSession...).
		GetWithContext(ctx, getApplicationURL).
		BodyJSON(&app)
	if err != nil {
		return nil, fmt.Errorf("error getting app admin/%s: %w", appz, err)
	}
	if app.ClientID == "" || app.ClientSecret == "" {
		return nil, fmt.Errorf("failed to get app admin/%s: blank client id/secret", appz)
	}
	return &app, nil
}
