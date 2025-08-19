package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/EWK20/event-processor/processor/internal/config"
	"github.com/EWK20/event-processor/processor/internal/models"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rs/zerolog/log"
)

var (
	ErrFailedToCreateClient = errors.New("failed to create SQS client")
	ErrFailedToGetQueueURL  = errors.New("failed to get queue URL")
	ErrFailedToGetDLQURL    = errors.New("failed to get DLQ URL")
)

type DB interface {
	Save(ctx context.Context, event models.Event) error
}

type Processor struct {
	Client   *sqs.Client
	QueueURL *string
	DLQURL   *string
	db       DB
}

func New(cfg config.AWS, db DB) (*Processor, error) {
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
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetQueueURL, err)
	}

	dlqQueueURL, err := sqsClient.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: &cfg.SQSDLQName,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetDLQURL, err)
	}

	return &Processor{
		Client:   sqsClient,
		QueueURL: queueURL.QueueUrl,
		DLQURL:   dlqQueueURL.QueueUrl,
		db:       db,
	}, nil
}

func (p *Processor) Run(ctx context.Context) {
	for {
		msgOutput, err := p.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            p.QueueURL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     5,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to receive messages")

			continue
		}

		for _, msg := range msgOutput.Messages {
			var event models.Event

			if err := json.Unmarshal([]byte(*msg.Body), &event); err != nil {
				log.Error().Err(err).Msg("event is invalid")

				if err := p.sendMsgToDLQ(&msg); err != nil {
					log.Error().Err(err).Msg("failed to send event to dead letter queue")
				}

				continue
			}

			if err := p.db.Save(context.Background(), event); err != nil {
				log.Error().Err(err).Msg("failed to save event to database")

				continue
			}

			// Delete message after successful insert
			_, err = p.Client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      p.QueueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to delete message from queue")
			}

			log.Info().Any("event", event).Msg("persisted an event")
		}
	}
}

func (p *Processor) sendMsgToDLQ(msg *types.Message) error {
	// Send message to DLQ
	_, err := p.Client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    p.DLQURL,
		MessageBody: msg.Body,
	})
	if err != nil {
		return fmt.Errorf("failed to send invalid message: %w", err)
	}

	// Delete message after sending to dlq
	_, err = p.Client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      p.QueueURL,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		return fmt.Errorf("failed to delete message from queue: %w", err)
	}

	return nil
}
