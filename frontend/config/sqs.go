package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

const (
	maxRetries = 3
	retryDelay = 1 * time.Second
)

var sqsClient *sqs.Client

// SendPrintRequest sends a print job request to SQS
func SendPrintRequest(customerEmail, photoID, processedS3Location string) error {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env", err)
	}

	// Load sqs url
	queueURL := os.Getenv("SQS_QUEUE_URL")

	// Load AWS config
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("failed to load configuration, ", err)
	}
	client := sqs.NewFromConfig(config)

	// Create message structure for SQS
	job := PrintJob{
		CustomerEmail:       customerEmail,
		PhotoID:             photoID,
		ProcessedS3Location: processedS3Location,
	}

	jobBytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// Send message to SQS
	// implement retry logic
	var previousError error
	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err = client.SendMessage(context.Background(), &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueURL),
			MessageBody: aws.String(string(jobBytes)),
			// DelaySeconds: 5,
		})

		if err == nil {
			log.Printf("Successfully sent print request to SQS for photos %s (attemtp - %d)", photoID, attempt+1)
			break
		}

		previousError = err
		log.Printf("Failed to send print request (attempt %d): %v", attempt+1, err)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to send print request after %d attempts: %w", maxRetries, previousError)
}
