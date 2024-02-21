package types

import (
	"errors"

	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/property"
)

const (
	Project     = "project"
	Region      = "region"
	Zone        = "zone"
	Credentials = "credentials"
)

type Credential struct {
	Project     string
	Region      string
	Zone        string
	Credentials string
}

func GetCredential(configData crypto.Properties) (*Credential, error) {
	var (
		cred = &Credential{}
		ok   bool
		err  error
	)

	cred.Project, ok, err = property.GetString(configData[Project].Value)
	if !ok || err != nil || cred.Project == "" {
		return nil, errors.New("project is empty")
	}

	cred.Region, ok, err = property.GetString(configData[Region].Value)
	if !ok || err != nil || cred.Region == "" {
		return nil, errors.New("region is empty")
	}

	cred.Zone, ok, err = property.GetString(configData[Zone].Value)
	if !ok || err != nil || cred.Zone == "" {
		return nil, errors.New("zone is empty")
	}

	cred.Credentials, ok, err = property.GetString(configData[Credentials].Value)
	if !ok || err != nil || cred.Credentials == "" {
		return nil, errors.New("credentials is empty")
	}

	return cred, nil
}
