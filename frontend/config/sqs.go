package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var sqsClient *sqs.Client
var queueURL = os.Getenv("SQS_QUEUE_URL")

// SendPrintRequest sends a print job request to SQS
func SendPrintRequest(customerEmail, photoID, processedS3Location string) error {

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
	_, err = client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:     aws.String(queueURL),
		MessageBody:  aws.String(string(jobBytes)),
		DelaySeconds: 5,
	})

	if err != nil {
		log.Printf("Failed to send print request to SQS: %v", err)
	}

	log.Println("Print request sent to SQS", photoID)
	return err
}
