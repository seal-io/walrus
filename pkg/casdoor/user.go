package casdoor

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/seal-io/seal/utils/req"
)

func SignInUser(ctx context.Context, app, org, usr, pwd string) ([]*req.HttpCookie, error) {
	var loginURL = fmt.Sprintf("%s/api/login", endpoint.Get())
	var loginReq = map[string]interface{}{
		"type":         "login",
		"application":  app,
		"organization": org,
		"username":     usr,
		"password":     pwd,
		"autoSignin":   true,
	}
	var loginRespBody struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	var loginResp = req.HTTPRequest().
		WithBodyJSON(loginReq).
		PostWithContext(ctx, loginURL)
	var err = loginResp.BodyJSON(&loginRespBody)
	if err != nil {
		return nil, fmt.Errorf("error signing in user %s/%s: %w", org, usr, err)
	}
	if loginRespBody.Status == "error" {
		return nil, fmt.Errorf("failed to sign in user %s/%s: %s", org, usr, loginRespBody.Msg)
	}
	var userSession = loginResp.Cookies()
	if len(userSession) == 0 {
		return nil, fmt.Errorf("faield to sign in user %s/%s", org, usr)
	}
	return userSession, nil
}

func SignOutUser(ctx context.Context, userSession []*req.HttpCookie) error {
	var logoutURL = fmt.Sprintf("%s/api/logout", endpoint.Get())
	var err = req.HTTPRequest().
		WithCookies(userSession...).
		PostWithContext(ctx, logoutURL).
		Error()
	if err != nil {
		return fmt.Errorf("error signing out out: %w", err)
	}
	return nil
}

func CreateUser(ctx context.Context, clientID, clientSecret string, app, org, usr, pwd string) error {
	var createUserURL = fmt.Sprintf("%s/api/add-user", endpoint.Get())
	var createUserReq = map[string]any{
		"owner":             org,
		"name":              usr,
		"type":              "normal-user",
		"password":          pwd,
		"displayName":       org + "/" + usr,
		"isAdmin":           true,
		"isGlobalAdmin":     true,
		"signupApplication": app,
	}
	var createUserResp struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	var err = req.HTTPRequest().
		WithBasicAuth(clientID, clientSecret).
		WithBodyJSON(createUserReq).
		PostWithContext(ctx, createUserURL).
		BodyJSON(&createUserResp)
	if err != nil {
		return fmt.Errorf("error creating user %s/%s: %w", org, usr, err)
	}
	if createUserResp.Status == "error" {
		return fmt.Errorf("failed to create the user %s/%s: %s", org, usr, createUserResp.Msg)
	}
	return nil
}

func DeleteUser(ctx context.Context, clientID, clientSecret string, org, usr string) error {
	var deleteUserURL = fmt.Sprintf("%s/api/delete-user", endpoint.Get())
	var deleteUserReq = map[string]any{
		"owner": org,
		"name":  usr,
	}
	var deleteUserResp struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	var err = req.HTTPRequest().
		WithBasicAuth(clientID, clientSecret).
		WithBodyJSON(deleteUserReq).
		PostWithContext(ctx, deleteUserURL).
		BodyJSON(&deleteUserResp)
	if err != nil {
		return fmt.Errorf("error deleting user %s/%s: %w", org, usr, err)
	}
	if deleteUserResp.Status == "error" {
		return fmt.Errorf("failed to delete user %s/%s: %s", deleteUserResp.Msg, org, usr)
	}
	return nil
}

func UpdateUserPassword(ctx context.Context, clientID, clientSecret string, org, usr, oldPwd, newPwd string) error {
	var setPwdURL = fmt.Sprintf("%s/api/set-password", endpoint.Get())
	var setPwdReq = url.Values{
		"userOwner":   []string{org},
		"userName":    []string{usr},
		"newPassword": []string{newPwd},
		"oldPassword": []string{oldPwd},
	}
	var setPwdResp struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	var err = req.HTTPRequest().
		WithBasicAuth(clientID, clientSecret).
		WithBodyForm(setPwdReq).
		PostWithContext(ctx, setPwdURL).
		BodyJSON(&setPwdResp)
	if err != nil {
		return fmt.Errorf("error setting password: %w", err)
	}
	if setPwdResp.Status == "error" {
		return fmt.Errorf("failed to set password: %s", setPwdResp.Msg)
	}
	return nil
}

type User struct {
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	DisplayName   string `json:"displayName"`
	IsAdmin       bool   `json:"isAdmin"`
	IsGlobalAdmin bool   `json:"isGlobalAdmin"`
}

func GetUser(ctx context.Context, clientID, clientSecret string, org, usr string) (*User, error) {
	var getUserURL = fmt.Sprintf("%s/api/get-user?id=%s/%s", endpoint.Get(), org, usr)
	var user User
	var err = req.HTTPRequest().
		WithBasicAuth(clientID, clientSecret).
		GetWithContext(ctx, getUserURL).
		BodyJSON(&user)
	if err != nil {
		return nil, fmt.Errorf("error getting user %s/%s: %w", org, usr, err)
	}
	if user.Owner == "" || user.Name == "" {
		return nil, fmt.Errorf("failed to get user %s/%s: not found", org, usr)
	}
	return &user, nil
}

type UserInfo struct {
	Organization string `json:"organization"`
	Name         string `json:"name"`
}

func GetUserInfo(ctx context.Context, userSession []*req.HttpCookie) (*UserInfo, error) {
	var getAccountURL = fmt.Sprintf("%s/api/get-account", endpoint.Get())
	var account struct {
		Sub          string `json:"sub"`
		Name         string `json:"name"`
		Organization struct {
			Name string ` json:"name"`
		} `json:"data2"`
	}
	var err = req.HTTPRequest().
		WithCookies(userSession...).
		GetWithContext(ctx, getAccountURL).
		BodyJSON(&account)
	if err != nil {
		return nil, fmt.Errorf("error getting user account: %w", err)
	}
	if account.Sub == "" {
		return nil, errors.New("failed to get user account")
	}
	return &UserInfo{
		Organization: account.Organization.Name,
		Name:         account.Name,
	}, nil
}
