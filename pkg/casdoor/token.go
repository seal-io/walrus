package casdoor

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/seal-io/seal/utils/req"
)

const (
	// NeverExpiresInSeconds gives a large number to simulate "never expires",
	// ref to https://github.com/casdoor/casdoor/issues/803.
	neverExpiresInSeconds = 50 * 365 * 24 * 60 * 60
)

type Token struct {
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

func CreateToken(ctx context.Context, clientID, clientSecret string, usr string, expiryInSecondsPtr *int) (*Token, error) {
	var createTokenURL = fmt.Sprintf("%s/api/login/oauth/access_token", endpoint.Get())
	var expiryInSeconds int
	if expiryInSecondsPtr != nil {
		expiryInSeconds = *expiryInSecondsPtr
	}
	if expiryInSeconds <= 0 || expiryInSeconds > neverExpiresInSeconds {
		expiryInSeconds = neverExpiresInSeconds
	}
	var createTokenReq = url.Values{
		"grant_type":        []string{"client_credentials"},
		"username":          []string{usr},
		"expiry_in_seconds": []string{strconv.Itoa(expiryInSeconds)},
	}
	var createTokenResp struct {
		Owner        string `json:"owner"`
		Name         string `json:"name"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	var err = req.HTTPRequest().
		WithBasicAuth(clientID, clientSecret).
		WithBodyForm(createTokenReq).
		PostWithContext(ctx, createTokenURL).
		BodyJSON(&createTokenResp)
	if err != nil {
		return nil, fmt.Errorf("error creating token: %w", err)
	}
	if createTokenResp.AccessToken == "" {
		return nil, errors.New("failed to create token: blank access token")
	}
	return &Token{
		Owner:        createTokenResp.Owner,
		Name:         createTokenResp.Name,
		AccessToken:  createTokenResp.AccessToken,
		RefreshToken: createTokenResp.RefreshToken,
		ExpiresIn:    createTokenResp.ExpiresIn,
	}, nil
}

func DeleteToken(ctx context.Context, clientID, clientSecret string, owner, name string) error {
	var deleteTokenURL = fmt.Sprintf("%s/api/delete-token", endpoint.Get())
	var deleteTokenReq = Token{
		Owner: owner,
		Name:  name,
	}
	var deleteTokenResp struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	var err = req.HTTPRequest().
		WithBasicAuth(clientID, clientSecret).
		WithBodyJSON(deleteTokenReq).
		PostWithContext(ctx, deleteTokenURL).
		BodyJSON(&deleteTokenResp)
	if err != nil {
		return fmt.Errorf("error deleting token: %w", err)
	}
	if deleteTokenResp.Status == "error" {
		return fmt.Errorf("failed to delete token: %s", deleteTokenResp.Msg)
	}
	return nil
}

type Introspection struct {
	Organization string `json:"organization"`
	UserName     string `json:"username"`
	Active       bool   `json:"active"`
	Exp          int64  `json:"exp"`
}

func IntrospectToken(ctx context.Context, clientID, clientSecret string, token string) (*Introspection, error) {
	var introspectTokenURL = fmt.Sprintf("%s/api/login/oauth/introspect", endpoint.Get())
	var introspectTokenReq = url.Values{
		"token": []string{token},
	}
	var introspectTokenResp struct {
		Introspection `json:",inline"`
		Status        string `json:"status"`
		Msg           string `json:"msg"`
	}
	var err = req.HTTPRequest().
		WithBasicAuth(clientID, clientSecret).
		WithBodyForm(introspectTokenReq).
		PostWithContext(ctx, introspectTokenURL).
		BodyJSON(&introspectTokenResp)
	if err != nil {
		return nil, fmt.Errorf("error introspecting token: %w", err)
	}
	if introspectTokenResp.Status == "error" {
		return nil, fmt.Errorf("failed to introspect token: %s", introspectTokenResp.Msg)
	}
	return &introspectTokenResp.Introspection, nil
}
