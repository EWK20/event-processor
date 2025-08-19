package processor_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/EWK20/event-processor/processor/internal/config"
	"github.com/EWK20/event-processor/processor/internal/models"
	"github.com/EWK20/event-processor/processor/internal/processor"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	type Test struct {
		input models.Event
	}

	testCases := map[string]Test{
		"Successful Run": {
			input: models.Event{
				EventType: "transaction_approved",
				ClientID:  "client_789",
				Payload: map[string]any{
					"transaction_id": "txn_123",
					"amount":         "125.33",
					"currency":       "GBP",
				},
				Timestamp: time.Now().UTC(),
			},
		},
	}

	for scenario, test := range testCases {
		t.Run(scenario, func(t *testing.T) {
			awsCfg := config.AWS{
				AWSRegion:          "us-east-1",
				AWSAccessKeyID:     "test",
				AWSSecretAccessKey: "test",
				SQSEndpoint:        "http://localhost:4566",
				SQSQueueName:       "test-queue",
				SQSDLQName:         "test-queue-dlq",
			}

			fakeDB := NewFakeDB()

			processor, err := processor.New(awsCfg, fakeDB)
			require.NoError(t, err)

			body, err := json.Marshal(test.input)
			require.NoError(t, err)

			_, err = processor.Client.SendMessage(t.Context(), &sqs.SendMessageInput{
				QueueUrl:    processor.QueueURL,
				MessageBody: aws.String(string(body)),
			})
			require.NoError(t, err)

			// Run processor in a goroutine so it consumes the message
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			go func() {
				processor.Run(ctx) // blocks forever
			}()

			// Wait until the message is processed (poll the mock DB)
			require.Eventually(t, func() bool {
				return len(fakeDB.events) > 0
			}, 10*time.Second, 500*time.Millisecond, "event was not processed in time")

			// Check that event persisted correctly
			require.Equal(t, test.input.ClientID, fakeDB.events[len(fakeDB.events)-1].ClientID)
			require.Equal(t, test.input.EventType, fakeDB.events[len(fakeDB.events)-1].EventType)

		})
	}
}
