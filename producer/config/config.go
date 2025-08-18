package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrMissingCfg = errors.New("required config missing")
)

type Config struct {
	SQSQueueName       string
	SQSEndpoint        string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
}

func New() (*Config, error) {
	var cfg Config

	if cfg.SQSQueueName = os.Getenv("SQS_QUEUE_NAME"); cfg.SQSQueueName == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingCfg, "SQS_QUEUE_NAME")
	}

	if cfg.SQSEndpoint = os.Getenv("SQS_ENDPOINT"); cfg.SQSEndpoint == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingCfg, "SQS_ENDPOINT")
	}

	if cfg.AWSRegion = os.Getenv("AWS_REGION"); cfg.AWSRegion == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingCfg, "AWS_REGION")
	}

	if cfg.AWSAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID"); cfg.AWSAccessKeyID == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingCfg, "AWS_ACCESS_KEY_ID")
	}

	if cfg.AWSSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY"); cfg.AWSSecretAccessKey == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingCfg, "AWS_SECRET_ACCESS_KEY")
	}

	return &cfg, nil
}
