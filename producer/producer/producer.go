package producer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/EWK20/event-processor/producer/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rs/zerolog/log"
)

var (
	ErrFailedToCreateClient = errors.New("failed to create SQS client")
	ErrFailedToGetQueueURL  = errors.New("failed to get queue URL")
)

type Event struct {
	EventType string `json:"event_type"`
	ClientID  string `json:"client_id"`
	Payload   any    `json:"payload"`
	Timestamp string `json:"timestamp"`
}

type Producer struct {
	sqsClient *sqs.Client
	queueURL  *sqs.GetQueueUrlOutput
}

func New(cfg config.Config) (*Producer, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithRegion(cfg.AWSRegion),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWSAccessKeyID,
				cfg.AWSSecretAccessKey,
				"",
			),
		),
		awsConfig.WithBaseEndpoint(cfg.SQSEndpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateClient, err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)

	queueURL, err := sqsClient.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: &cfg.SQSQueueName,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateClient, err)
	}

	return &Producer{
		sqsClient,
		queueURL,
	}, nil
}

func (p *Producer) Run() {
	log.Info().Msg("Producing messages...")

	clientIDs := []string{"client_123", "client_456", "client_789"}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random transactions
	for {
		event := Event{
			EventType: "transaction_approved",
			ClientID:  clientIDs[rand.Intn(len(clientIDs))],
			Payload: map[string]any{
				"transaction_id": "txn_" + fmt.Sprint(rand.Intn(1000)),
				"amount":         fmt.Sprintf("%d.%d", rand.Intn(1000), rand.Intn(99)),
				"currency":       "GBP",
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		msg, err := json.Marshal(&event)
		if err != nil {
			panic(err)
		}

		_, err = p.sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
			QueueUrl:    p.queueURL.QueueUrl,
			MessageBody: aws.String(string(msg)),
		})
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 15)
	}
}
