package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *Config
		wantErr bool
	}{
		{
			name:    "empty",
			config:  "",
			wantErr: true,
		},
		{
			name:    "invalid",
			config:  "invalid",
			wantErr: true,
		},
		{
			name:   "valid",
			config: `s3://accessKey:secretAccessKey@endpoint:9000/bucketName?sslmode=disable`,
			want: &Config{
				Endpoint:        "endpoint:9000",
				Bucket:          "bucketName",
				AccessKeyID:     "accessKey",
				SecretAccessKey: "secretAccessKey",
				Secure:          false,
			},
		},
		{
			name:   "valid-secure",
			config: `s3://accessKey1:xdrlT7a2x*sd34s@endpoint:9000/bucketName?sslmode=enable`,
			want: &Config{
				Endpoint:        "endpoint:9000",
				Bucket:          "bucketName",
				AccessKeyID:     "accessKey1",
				SecretAccessKey: "xdrlT7a2x*sd34s",
				Secure:          true,
			},
		},
		{
			name:   "without-port",
			config: `s3://ak1:sk2sk2sk2@yourdomain.com/bucketName`,
			want: &Config{
				Endpoint:        "yourdomain.com",
				Bucket:          "bucketName",
				AccessKeyID:     "ak1",
				SecretAccessKey: "sk2sk2sk2",
				Secure:          true,
			},
		},
		{
			name:   "with-region",
			config: `s3://accessKey:secret@yourdomain/bucket?region=ap-northeast-1`,
			want: &Config{
				Endpoint:        "yourdomain",
				Bucket:          "bucket",
				AccessKeyID:     "accessKey",
				SecretAccessKey: "secret",
				Region:          "ap-northeast-1",
				Secure:          true,
			},
		},
	}
	for _, tt := range tests {
		got, err := ParseConfig(tt.config)
		if err != nil {
			if !tt.wantErr {
				t.Errorf("%s: ParseConfig() unexpected error: %v", tt.name, err)
			}

			continue
		}

		if !assert.Equal(t, tt.want, got) {
			t.Errorf("%s: ParseConfig() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
