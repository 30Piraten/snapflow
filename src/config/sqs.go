package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	maxRetries = 3
	retryDelay = 1 * time.Second
)

// PrintJob represents a print request
type PrintJob struct {
	CustomerEmail       string `json:"customer_email"`
	PhotoID             string `json:"photo_id"`
	ProcessedS3Location string `json:"processed_s3_location"`
}

var sqsClient *sqs.Client

// SendPrintRequest sends a print job request to SQS
func SendPrintRequest(customerEmail, photoID, processedS3Location string) error {
	queueURL := os.Getenv("SQS_QUEUE_URL")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1")) // Explicit region
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}
	client := sqs.NewFromConfig(cfg)

	job := PrintJob{
		CustomerEmail:       customerEmail,
		PhotoID:             photoID,
		ProcessedS3Location: processedS3Location,
	}

	jobBytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal print job: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5-second timeout
	defer cancel()

	var previousError error
	for attempt := 0; attempt < maxRetries; attempt++ {
		output, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:     aws.String(queueURL),
			MessageBody:  aws.String(string(jobBytes)),
			DelaySeconds: 5,
		})

		if err == nil {
			log.Printf("Successfully sent print request to SQS for photos %s (attempt %d), Output: %+v", photoID, attempt+1, output)
			return nil // Success!
		}

		// tautological condition spoted -> err != nil
		if previousError != nil {
			// Log the original error
			log.Printf("Original error: %v", err)

			// Try to unwrap the error (Go 1.13 and later)
			if unwrappedErr := errors.Unwrap(err); unwrappedErr != nil {
				log.Printf("Unwrapped error: %v", unwrappedErr)

				// Check if it's a net.Error (for network issues)
				if netErr, ok := unwrappedErr.(net.Error); ok {
					log.Printf("Network error: %v, Timeout: %v", netErr, netErr.Timeout())
				}
			}
			previousError = err // Capture the error
			log.Printf("Failed to send print request (attempt %d): %v", attempt+1, err)
		} else {
			log.Printf("SendMessage returned nil error, but failed (attempt %d)", attempt+1)
			previousError = fmt.Errorf("SendMessage returned nil error") // Create a placeholder error
		}
		previousError = err // Capture the error
		log.Printf("Failed to send print request (attempt %d): %v", attempt+1, err)

		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to send print request after %d attempts: %w", maxRetries, previousError)
}
