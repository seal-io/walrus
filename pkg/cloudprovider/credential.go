package cloudprovider

import (
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"
)

type Credential struct {
	AccessKey string
	SecretKey string
	Region    string
}

func CredentialFromConnector(conn *model.Connector) (*Credential, error) {
	var (
		cred Credential
		ok   bool
		err  error
	)

	for k, v := range conn.ConfigData {
		switch k {
		case "accessKey":
			cred.AccessKey, ok, err = v.GetString()
		case "accessSecret":
			cred.SecretKey, ok, err = v.GetString()
		case "region":
			cred.Region, ok, err = v.GetString()
		}

		if !ok || err != nil {
			return nil, fmt.Errorf("invalid connector config: %w", err)
		}
	}

	if cred.AccessKey == "" || cred.SecretKey == "" {
		return nil, errors.New("credential is empty")
	}

	if cred.Region == "" {
		return nil, errors.New("region is empty")
	}

	return &cred, nil
}
