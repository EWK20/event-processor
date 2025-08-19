package config_test

import (
	"testing"

	"github.com/EWK20/event-processor/processor/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Input struct {
	user               string
	password           string
	host               string
	port               string
	dbName             string
	sslMode            string
	sqsQueueName       string
	sqsDLQName         string
	sqsEndpoint        string
	awsRegion          string
	awsAccessKeyID     string
	awsSecretAccessKey string
}

type Test struct {
	input  Input
	output *config.Config
	err    error
}

func TestConfig(t *testing.T) {
	testCases := map[string]Test{
		"All Required Envs Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: &config.Config{
				DB: config.DB{
					User:     "user",
					Password: "password",
					Host:     "localhost",
					Port:     "5432",
					DBName:   "test",
					SSLMode:  "disable",
				},
				AWS: config.AWS{
					SQSQueueName:       "test-queue",
					SQSDLQName:         "test-queue-dlq",
					SQSEndpoint:        "test-endpoint:1234",
					AWSRegion:          "aws-region",
					AWSAccessKeyID:     "aws-access-key-id",
					AWSSecretAccessKey: "aws-secret-access-key",
				},
			},
			err: nil,
		},
		"DB User Not Set": {
			input: Input{
				user:               "",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"DB Password Not Set": {
			input: Input{
				user:               "user",
				password:           "",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"DB Host Not Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"DB Port Set To 9432": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "9432",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: &config.Config{
				DB: config.DB{
					User:     "user",
					Password: "password",
					Host:     "localhost",
					Port:     "9432",
					DBName:   "test",
					SSLMode:  "disable",
				},
				AWS: config.AWS{
					SQSQueueName:       "test-queue",
					SQSDLQName:         "test-queue-dlq",
					SQSEndpoint:        "test-endpoint:1234",
					AWSRegion:          "aws-region",
					AWSAccessKeyID:     "aws-access-key-id",
					AWSSecretAccessKey: "aws-secret-access-key",
				},
			},
			err: nil,
		},
		"DB Name Not Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"DB SSL Mode Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "allow",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: &config.Config{
				DB: config.DB{
					User:     "user",
					Password: "password",
					Host:     "localhost",
					Port:     "5432",
					DBName:   "test",
					SSLMode:  "allow",
				},
				AWS: config.AWS{
					SQSQueueName:       "test-queue",
					SQSDLQName:         "test-queue-dlq",
					SQSEndpoint:        "test-endpoint:1234",
					AWSRegion:          "aws-region",
					AWSAccessKeyID:     "aws-access-key-id",
					AWSSecretAccessKey: "aws-secret-access-key",
				},
			},
			err: nil,
		},
		"SQS Queue Name Not Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"SQS DLQ Name Not Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"SQS Endpoint Not Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"AWS Region Not Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"AWS Access Key ID Not Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "",
				awsSecretAccessKey: "aws-secret-access-key",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
		"AWS Secret Access Key Not Set": {
			input: Input{
				user:               "user",
				password:           "password",
				host:               "localhost",
				port:               "",
				dbName:             "test",
				sslMode:            "",
				sqsQueueName:       "test-queue",
				sqsDLQName:         "test-queue-dlq",
				sqsEndpoint:        "test-endpoint:1234",
				awsRegion:          "aws-region",
				awsAccessKeyID:     "aws-access-key-id",
				awsSecretAccessKey: "",
			},
			output: nil,
			err:    config.ErrMissingCfg,
		},
	}

	for scenario, test := range testCases {
		t.Run(scenario, func(t *testing.T) {
			setEnvs(t, test.input)

			cfg, err := config.New()

			if test.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.output, cfg)
		})
	}
}

func setEnvs(t *testing.T, input Input) {
	t.Helper()

	t.Setenv("DB_USER", input.user)
	t.Setenv("DB_PASSWORD", input.password)
	t.Setenv("DB_HOST", input.host)
	t.Setenv("DB_PORT", input.port)
	t.Setenv("DB_NAME", input.dbName)
	t.Setenv("DB_SSLMODE", input.sslMode)
	t.Setenv("AWS_REGION", input.awsRegion)
	t.Setenv("AWS_ACCESS_KEY_ID", input.awsAccessKeyID)
	t.Setenv("AWS_SECRET_ACCESS_KEY", input.awsSecretAccessKey)
	t.Setenv("SQS_ENDPOINT", input.sqsEndpoint)
	t.Setenv("SQS_QUEUE_NAME", input.sqsQueueName)
}
