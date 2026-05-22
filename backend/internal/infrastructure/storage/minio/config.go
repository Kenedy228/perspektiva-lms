package minio

import (
	"errors"
	"net/url"
	"strings"
)

type Config struct {
	Endpoint  string
	Bucket    string
	PublicURL string
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Endpoint) == "" {
		return errors.New("minio endpoint is required")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return errors.New("minio bucket is required")
	}
	if _, err := url.ParseRequestURI(c.Endpoint); err != nil {
		return err
	}
	if c.PublicURL != "" {
		if _, err := url.ParseRequestURI(c.PublicURL); err != nil {
			return err
		}
	}
	return nil
}
