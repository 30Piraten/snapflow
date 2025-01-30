package config

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var sqsClient *sqs.Client
var queueURL = "https://sqs.us-east-1.amazonaws.com/YOUR_ACCOUNT_ID/photo-print-queue"

func InitSQS() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	sqsClient = sqs.NewFromConfig(cfg)
}

// PrintJob represents a print request message
type PrintJob struct {
	CustomerEmail       string `json:"customer_email"`
	PhotoID             string `json:"photo_id"`
	ProcessedS3Location string `json:"processed_s3_location"`
}

// SendPrintRequest sends a print job request to SQS
func SendPrintRequest(customerEmail, photoID, processedS3Location string) error {
	InitSQS()

	job := PrintJob{
		CustomerEmail:       customerEmail,
		PhotoID:             photoID,
		ProcessedS3Location: processedS3Location,
	}

	jobBytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(jobBytes)),
	})

	if err != nil {
		log.Printf("Failed to send print request to SQS: %v", err)
	}

	return err
}

func SendToSQS(queueURL string, messageBody string) error {
	config, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := sqs.NewFromConfig(config)

	_, err = client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:     aws.String(queueURL),
		MessageBody:  aws.String(messageBody),
		DelaySeconds: 10, // Delay message for 10 seconds
	})

	if err != nil {
		log.Printf("Failed to send message to SQS: %v", err)
	}

	return err
}
