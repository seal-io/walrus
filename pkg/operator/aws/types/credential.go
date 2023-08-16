package types

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/seal-io/walrus/pkg/operator/types"
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
	cred, err := types.CredentialFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	wrap := Credential(*cred)

	return wrap.Config()
}
