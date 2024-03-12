package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/seal-io/walrus/pkg/dao/model"
)

type Manager struct {
	config      *Config
	minioClient *minio.Client
}

// NewManager creates a new storage manager.
func NewManager(conf *Config) (*Manager, error) {
	minioClient, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Secure: conf.Secure,
	})
	if err != nil {
		return nil, err
	}

	return &Manager{
		config:      conf,
		minioClient: minioClient,
	}, nil
}

func (m *Manager) GetAddress() string {
	return m.config.Endpoint
}

// SetRunPlan sets the run plan files to s3 storage.
func (m *Manager) SetRunPlan(ctx context.Context, run *model.ResourceRun, plan []byte) error {
	fileName := GetPlanFileName(run)

	planFile := bytes.NewReader(plan)
	bucketName := m.config.Bucket
	fileSize := int64(len(plan))

	if err := m.CheckValidBucketName(ctx, bucketName); err != nil {
		return err
	}

	_, err := m.minioClient.PutObject(ctx, bucketName, fileName, planFile, fileSize, minio.PutObjectOptions{})

	return err
}

func GetPlanFileName(run *model.ResourceRun) string {
	return fmt.Sprintf("walrus/project/%s/environment/%s/resource/%s.zip", run.ProjectID, run.EnvironmentID, run.ID)
}

func (m *Manager) GetRunPlan(ctx context.Context, run *model.ResourceRun) ([]byte, error) {
	fileName := GetPlanFileName(run)
	bucketName := m.config.Bucket

	object, err := m.minioClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(object)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *Manager) DeleteRunPlan(ctx context.Context, run *model.ResourceRun) error {
	fileName := GetPlanFileName(run)
	bucketName := m.config.Bucket

	err := m.minioClient.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil && minio.ToErrorResponse(err).Code != s3.ErrCodeNoSuchKey {
		return err
	}

	return nil
}

func (m *Manager) CheckValidBucketName(ctx context.Context, bucketName string) error {
	found, err := m.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if found {
		return nil
	}

	err = m.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: m.config.Region,
	})
	if err != nil {
		switch minio.ToErrorResponse(err).Code {
		case s3.ErrCodeBucketAlreadyOwnedByYou,
			s3.ErrCodeBucketAlreadyExists:
			return nil
		}

		return err
	}

	return nil
}
