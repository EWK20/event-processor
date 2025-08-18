#!/bin/bash
set -e

echo "🚀 Creating SQS queue: events..."
awslocal sqs create-queue --queue-name events
echo "✅ SQS queue created."

echo "🚀 Creating SQS queue: test-queue..."
awslocal sqs create-queue --queue-name test-queue
echo "✅ SQS test queue created."