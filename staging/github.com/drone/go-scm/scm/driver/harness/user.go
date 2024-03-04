// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"fmt"
	"strings"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(harnessUser)
	// the following is for the corporate version of Harness code
	tempUserService := *s
	// get the basepath
	basePath := tempUserService.client.BaseURL.Path
	// use the NG user endpoint
	basePath = strings.Replace(basePath, "code", "ng", 1)
	// set the new basepath
	tempUserService.client.BaseURL.Path = basePath
	// set the path
	path := fmt.Sprintf("api/user/currentUser")
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertHarnessUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	return "", nil, scm.ErrNotSupported
}

func (s *userService) ListEmail(context.Context, scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

//
// native data structures
//

type harnessUser struct {
	Status string `json:"status"`
	Data   struct {
		UUID             string      `json:"uuid"`
		Name             string      `json:"name"`
		Email            string      `json:"email"`
		Token            interface{} `json:"token"`
		Defaultaccountid string      `json:"defaultAccountId"`
		Intent           interface{} `json:"intent"`
		Accounts         []struct {
			UUID              string `json:"uuid"`
			Accountname       string `json:"accountName"`
			Companyname       string `json:"companyName"`
			Defaultexperience string `json:"defaultExperience"`
			Createdfromng     bool   `json:"createdFromNG"`
			Nextgenenabled    bool   `json:"nextGenEnabled"`
		} `json:"accounts"`
		Admin                          bool        `json:"admin"`
		Twofactorauthenticationenabled bool        `json:"twoFactorAuthenticationEnabled"`
		Emailverified                  bool        `json:"emailVerified"`
		Locked                         bool        `json:"locked"`
		Disabled                       bool        `json:"disabled"`
		Signupaction                   interface{} `json:"signupAction"`
		Edition                        interface{} `json:"edition"`
		Billingfrequency               interface{} `json:"billingFrequency"`
		Utminfo                        struct {
			Utmsource   interface{} `json:"utmSource"`
			Utmcontent  interface{} `json:"utmContent"`
			Utmmedium   interface{} `json:"utmMedium"`
			Utmterm     interface{} `json:"utmTerm"`
			Utmcampaign interface{} `json:"utmCampaign"`
		} `json:"utmInfo"`
		Externallymanaged bool `json:"externallyManaged"`
	} `json:"data"`
	Metadata      interface{} `json:"metaData"`
	Correlationid string      `json:"correlationId"`
}

//
// native data structure conversion
//

func convertHarnessUser(src *harnessUser) *scm.User {
	return &scm.User{
		Login: src.Data.Email,
		Email: src.Data.Email,
		Name:  src.Data.Name,
		ID:    src.Data.UUID,
	}
}
