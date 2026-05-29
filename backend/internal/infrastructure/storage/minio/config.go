package minio

import (
	"errors"
	"strings"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	PublicURL string
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Endpoint) == "" {
		return errors.New("minio endpoint is required")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return errors.New("minio bucket is required")
	}
	if strings.TrimSpace(c.AccessKey) == "" {
		return errors.New("minio access key is required")
	}
	if strings.TrimSpace(c.SecretKey) == "" {
		return errors.New("minio secret key is required")
	}
	return nil
}
