package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/joho/godotenv"
)

type SQSEvent struct {
	Records []struct {
		Body string `json:"body"`
	} `json:"Records"`
}

var (
	dynamoClient *dynamodb.Client
	snsClient    *sns.Client
	s3Client     *s3.Client
	snsTopicArn  string
)

// PrintJob represents a print request
type PrintJob struct {
	CustomerEmail       string `json:"customer_email"`
	PhotoID             string `json:"photo_id"`
	ProcessedS3Location string `json:"processed_s3_location"`
}

// Initialize AWS clients
func InitAWS() {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Initialize clients
	dynamoClient = dynamodb.NewFromConfig(config)
	snsClient = sns.NewFromConfig(config)

	// Get SNS topic ARN from environment variable
	snsTopicArn = os.Getenv("SNS_TOPIC_ARN")
	if snsTopicArn == "" {
		log.Fatalf("SNS_TOPIC_ARN environment variable is not set")
	}
}

// Simulated Print Function
func SimulatedPrint(job PrintJob) {
	fmt.Printf("🖨️ Printing photo: %s for %s\n", job.PhotoID, job.CustomerEmail)
	time.Sleep(5 * time.Second) // Simulate a 5-second print delay
	fmt.Println("✅ Print completed!")
}

// Process a single print job
func ProcessPrintJob(ctx context.Context, job PrintJob) error {
	// Load .env files
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ unable to load .env file: %v", err)
	}

	// Access .env variables
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	snsTopicArn := os.Getenv("SNS_TPOIC_ARN")

	// Simulate printing -> addd 5 seconds delay here
	SimulatedPrint(job)

	// Update DynamoDB status
	_, err = dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"customer_email": &types.AttributeValueMemberS{
				Value: job.CustomerEmail,
			},
			"photo_id": &types.AttributeValueMemberS{
				Value: job.PhotoID,
			},
		},
		UpdateExpression: aws.String("SET photo_status = :s"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{
				Value: "printed",
			},
		},
	})

	if err != nil {
		log.Printf("❌ Failed to update DynamoDB: %v", err)
		return err
	}

	// Send SNS notification
	message := fmt.Sprintf("📣 Hello %s, your photo (ID: %s) had been printed and it is ready for pickup!", job.CustomerEmail, job.PhotoID)
	_, err = snsClient.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(snsTopicArn),
	})

	if err != nil {
		log.Printf("❌ Failed to send SNS: %v", err)
		return err
	}

	log.Printf("✅ print job completed for %s", job.PhotoID)
	return nil
}

func Handler(ctx context.Context, event SQSEvent) error { // Add the error return
	var printJob PrintJob

	if snsClient == nil || dynamoClient == nil {
		InitAWS()
	}

	for _, record := range event.Records {
		err := json.Unmarshal([]byte(record.Body), &printJob)

		if err != nil {
			log.Printf("❌ failed to unmarshal SQS message: %v", err) // More specific message
			continue                                                 // Important: Continue to the next message
		}

		// Process the job
		if err := ProcessPrintJob(ctx, printJob); err != nil { // Check for errors
			log.Printf("❌ Error processing print job: %v", err)
			// Important: Decide how you want to handle errors.
			// Options:
			// 1. Continue to the next message (like below)
			// 2. Return the error to stop processing and cause a retry
			continue // Continue to the next message
		}
	}

	return nil // Return nil if all messages were processed (even with some errors)
}

func main() {
	lambda.Start(Handler)
}
