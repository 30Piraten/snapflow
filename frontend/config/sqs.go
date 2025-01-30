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

// func InitSQS() {
// 	if sqsClient != nil {
// 		cfg, err := config.LoadDefaultConfig(context.TODO())
// 		if err != nil {
// 			log.Fatalf("unable to load SDK config, %v", err)
// 		}
// 		sqsClient = sqs.NewFromConfig(cfg)
// 	}
// }

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
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(jobBytes)),
	})

	if err != nil {
		log.Printf("Failed to send print request to SQS: %v", err)
	}

	log.Println("Print request sent to SQS", photoID)
	return err
}

// Send an SQS message first (with a delay).
// // Trigger Lambda from SQS to send SNS after the delay.
// func SendToSQS(queueURL string, messageBody string) error {

// 	// Load AWS config
// 	config, err := config.LoadDefaultConfig(context.Background())
// 	if err != nil {
// 		log.Fatalf("unable to load SDK config, %v", err)
// 	}
// 	client := sqs.NewFromConfig(config)

// 	_, err = client.SendMessage(context.Background(), &sqs.SendMessageInput{
// 		QueueUrl:     aws.String(queueURL),
// 		MessageBody:  aws.String(messageBody),
// 		DelaySeconds: 10, // Delay message for 10 seconds
// 	})

// 	if err != nil {
// 		log.Printf("Failed to send message to SQS: %v", err)
// 	}

// 	return err
// }
