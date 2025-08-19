#!/bin/bash
set -e

echo "🚀 Creating SQS queue: events..."
awslocal sqs create-queue --queue-name events
echo "✅ SQS queue created."

echo "🚀 Creating SQS queue: events-dlq..."
awslocal sqs create-queue --queue-name events-dlq
echo "✅ SQS dead letter queue created."

echo "🚀 Creating SQS queue: test-queue..."
awslocal sqs create-queue --queue-name test-queue
echo "✅ SQS test queue created."

echo "🚀 Creating SQS queue: test-queue-dlq..."
awslocal sqs create-queue --queue-name test-queue-dlq
echo "✅ SQS test queue dlq created."