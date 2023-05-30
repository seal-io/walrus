package types

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/dao/types/crypto"
)

type CredentialKeyType string

const (
	CredentialKey CredentialKeyType = "credential"
)

const (
	AccessKey    = "access_key"
	AccessSecret = "secret_key"
	Region       = "region"
)

type Credential struct {
	AccessKey    string
	AccessSecret string
	Region       string
}

func GetCredential(configData crypto.Properties) (*Credential, error) {
	var (
		cred = &Credential{}
		ok   bool
		err  error
	)

	cred.AccessKey, ok, err = configData[AccessKey].GetString()
	if !ok || err != nil || cred.AccessKey == "" {
		return nil, errors.New("accessKey is empty")
	}

	cred.AccessSecret, ok, err = configData[AccessSecret].GetString()
	if !ok || err != nil || cred.AccessSecret == "" {
		return nil, errors.New("accessSecret is empty")
	}

	cred.Region, ok, err = configData[Region].GetString()
	if !ok || err != nil || cred.Region == "" {
		return nil, errors.New("region is empty")
	}

	return cred, nil
}

func CredentialFromCtx(ctx context.Context) (*Credential, error) {
	cred, ok := ctx.Value(CredentialKey).(*Credential)
	if !ok {
		return nil, errors.New("not found credential from context")
	}

	if cred.AccessKey == "" {
		return nil, errors.New("accessKey is empty")
	}

	if cred.AccessSecret == "" {
		return nil, errors.New("secretKey is empty")
	}

	if cred.Region == "" {
		return nil, errors.New("region is empty")
	}

	return cred, nil
}
