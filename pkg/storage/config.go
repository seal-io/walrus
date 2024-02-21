package storage

import (
	"errors"
	"net/url"
	"strings"
)

// ErrInvalidFormat is returned when the provided s3 credential is not in the
// expected format.
var ErrInvalidFormat = errors.New("invalid s3 credential format")

// Config contains all configuration necessary to connect to an s3 compatible
// server.
type Config struct {
	Endpoint        string
	Bucket          string
	Secure          bool
	AccessKeyID     string
	SecretAccessKey string
}

func (c *Config) GetAddress() string {
	return c.Endpoint
}

// ParseConfig parses the string s and extracts the s3 credentail.
// The supported format is: s3://ak:sk@region.endpoint/bucket?sslMode=disable.
func ParseConfig(s string) (*Config, error) {
	if !strings.HasPrefix(s, "s3://") {
		return nil, ErrInvalidFormat
	}

	s = strings.TrimPrefix(s, "s3://")

	// Parse the URL-style string.
	u, err := url.Parse("//" + s)
	if err != nil {
		return nil, err
	}

	conf := &Config{
		Endpoint: u.Host,
		Bucket:   strings.TrimPrefix(u.Path, "/"),
	}

	// Extract access key and secret key from userinfo.
	if u.User != nil {
		conf.AccessKeyID = u.User.Username()
		if password, passwordSet := u.User.Password(); passwordSet {
			conf.SecretAccessKey = password
		}
	}

	// Check query params for SSL mode.
	q := u.Query()
	if q.Get("sslMode") == "disable" {
		conf.Secure = false
	} else {
		conf.Secure = true
	}

	return conf, nil
}

func NewConfig(conf string) (*Config, error) {
	return ParseConfig(conf)
}
