package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

var (
	ErrMissingCfg = errors.New("required config missing")
)

type DB struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

type AWS struct {
	SQSQueueName       string
	SQSEndpoint        string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
}

type Config struct {
	DB  DB
	AWS AWS
}

func New() (*Config, error) {
	var cfg Config
	if err := getDatabaseCfg(&cfg.DB); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if err := getAWSCfg(&cfg.AWS); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &cfg, nil
}

func getDatabaseCfg(cfg *DB) error {
	if cfg.User = os.Getenv("DB_USER"); cfg.User == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "DB_USER")
	}
	if cfg.Password = os.Getenv("DB_PASSWORD"); cfg.Password == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "DB_PASSWORD")
	}
	if cfg.Host = os.Getenv("DB_HOST"); cfg.Host == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "DB_HOST")
	}
	if cfg.Port = os.Getenv("DB_PORT"); cfg.Port == "" {
		log.Warn().Msg("DB_PORT is not set. Using default")

		cfg.Port = "5432"
	}
	if cfg.DBName = os.Getenv("DB_NAME"); cfg.DBName == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "DB_NAME")
	}
	if cfg.SSLMode = os.Getenv("DB_SSLMODE"); cfg.SSLMode == "" {
		log.Warn().Msg("DB_SSLMODE is not set. Using default")

		cfg.SSLMode = "disable"
	}

	return nil
}

func getAWSCfg(cfg *AWS) error {
	if cfg.SQSQueueName = os.Getenv("SQS_QUEUE_NAME"); cfg.SQSQueueName == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "SQS_QUEUE_NAME")
	}

	if cfg.SQSEndpoint = os.Getenv("SQS_ENDPOINT"); cfg.SQSEndpoint == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "SQS_ENDPOINT")
	}

	if cfg.AWSRegion = os.Getenv("AWS_REGION"); cfg.AWSRegion == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "AWS_REGION")
	}

	if cfg.AWSAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID"); cfg.AWSAccessKeyID == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "AWS_ACCESS_KEY_ID")
	}

	if cfg.AWSSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY"); cfg.AWSSecretAccessKey == "" {
		return fmt.Errorf("%w: %s", ErrMissingCfg, "AWS_SECRET_ACCESS_KEY")
	}

	return nil
}
