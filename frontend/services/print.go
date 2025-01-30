package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/services/sns"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/config"
)

// DUMMY PRINT SERVICE
var sqsClient *sqs.Client
var queueURL = "https://sqs.us-east-1.amazonaws.com/YOUR_ACCOUNT_ID/photo-print-queue"

func InitSQS() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	sqsClient = sqs.NewFromConfig(cfg)
}

// PrintJob represents a print request
type PrintJob struct {
	CustomerEmail       string `json:"customer_email"`
	PhotoID             string `json:"photo_id"`
	ProcessedS3Location string `json:"processed_s3_location"`
}

// PollPrintQueue listens for new messages and processes them
func PollPrintQueue() {
	InitSQS()
	for {
		// Poll for messages
		output, err := sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     5,
		})

		if err != nil {
			log.Printf("Error polling SQS: %v", err)
			continue
		}

		for _, message := range output.Messages {
			var job PrintJob
			err := json.Unmarshal([]byte(*message.Body), &job)
			if err != nil {
				log.Printf("Error decoding SQS message: %v", err)
				continue
			}

			// Simulate printing
			log.Printf("Printing photo %s for customer %s...", job.PhotoID, job.CustomerEmail)
			time.Sleep(5 * time.Second) // Simulate processing time

			// Update DynamoDB with print status
			err = dynamodb.UpdatePrintStatus(job.CustomerEmail, job.PhotoID, "printed")
			if err != nil {
				log.Printf("Error updating DynamoDB: %v", err)
			}

			// Send print completion notification
			err = sns.SendPrintNotification(job.CustomerEmail, job.PhotoID)
			if err != nil {
				log.Printf("Error sending notification: %v", err)
			}

			// Delete message from SQS queue
			_, err = sqsClient.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("Error deleting message from SQS: %v", err)
			}
		}
	}
}

func main() {
	PollPrintQueue()
}
