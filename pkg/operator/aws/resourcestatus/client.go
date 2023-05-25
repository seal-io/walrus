package resourcestatus

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func ec2Client(ctx context.Context) (*ec2.Client, error) {
	cfg, err := ConfigFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return ec2.NewFromConfig(*cfg), nil
}

func rdsClient(ctx context.Context) (*rds.Client, error) {
	cfg, err := ConfigFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return rds.NewFromConfig(*cfg), nil
}

func cloudfrontClient(ctx context.Context) (*cloudfront.Client, error) {
	cfg, err := ConfigFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return cloudfront.NewFromConfig(*cfg), nil
}

func elasticCacheClient(ctx context.Context) (*elasticache.Client, error) {
	cfg, err := ConfigFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return elasticache.NewFromConfig(*cfg), nil
}

func elbClient(ctx context.Context) (*elbv2.ELBV2, error) {
	cred, err := credentialFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&awsv1.Config{
		Region:      awsv1.String(cred.Region),
		Credentials: credentials.NewStaticCredentials(cred.AccessKey, cred.AccessSecret, ""),
	})
	if err != nil {
		return nil, err
	}

	return elbv2.New(sess), nil
}

func eksClient(ctx context.Context) (*eks.Client, error) {
	cfg, err := ConfigFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return eks.NewFromConfig(*cfg), nil
}
