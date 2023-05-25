package resourcestatus

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/seal-io/seal/pkg/operator/types"
)

// Credential is a struct use to implement aws credentials.
type Credential types.Credential

// Retrieve implement aws.CredentialsProvider interface.
func (a Credential) Retrieve(_ context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     a.AccessKey,
		SecretAccessKey: a.AccessSecret,
	}, nil
}

// Config creates an AWS SDK V2 Config.
func (a Credential) Config() (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(a), config.WithRegion(a.Region))
	if err != nil {
		return &cfg,
			fmt.Errorf(
				"error create aws config with accessKey %s region %s: %w",
				a.AccessKey,
				a.Region,
				err,
			)
	}

	return &cfg, nil
}

// ConfigFromCtx get credential from context and creates an AWS SDK V2 Config.
func ConfigFromCtx(ctx context.Context) (*aws.Config, error) {
	cred, err := credentialFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	wrap := Credential(*cred)

	return wrap.Config()
}

func credentialFromCtx(ctx context.Context) (*types.Credential, error) {
	cred, ok := ctx.Value(types.CredentialKey).(*types.Credential)
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
