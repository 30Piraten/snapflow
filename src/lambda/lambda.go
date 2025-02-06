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

// PrintJob holds the print request
type PrintJob struct {
	CustomerEmail       string `json:"customer_email"`
	PhotoID             string `json:"photo_id"`
	ProcessedS3Location string `json:"processed_s3_location"`
}

// Initialize AWS clients -> DynamoDB and SNS
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

// SimulatedPrint handles a dummy printer using a delay sequence
func SimulatedPrint(job PrintJob) {
	fmt.Printf("🖨️ Printing photo: %s for %s\n", job.PhotoID, job.CustomerEmail)
	// Simulate a 10-second print delay
	time.Sleep(10 * time.Second)
	fmt.Println("✅ Print completed!")
}

// ProcessPrintJob processes a print job by simulating
// the printing process, updating the photo status in
// DynamoDB to "printed", and sending a notification via SNS.
func ProcessPrintJob(ctx context.Context, job PrintJob) error {

	// Access .env variables
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	snsTopicArn := os.Getenv("SNS_TOPIC_ARN")

	// Simulate printing -> add 5 seconds delay here
	SimulatedPrint(job)

	// Update DynamoDB status
	_, err := dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
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

	time.Sleep(10 * time.Second)
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

// Handler function processes a print job by unmarshaling
// the message body into a PrintJob, then calls ProcessPrintJob
// to simulate the printing process.
func Handler(ctx context.Context, event SQSEvent) error {
	var printJob PrintJob

	if snsClient == nil || dynamoClient == nil {
		InitAWS()
	}

	for _, record := range event.Records {
		err := json.Unmarshal([]byte(record.Body), &printJob)
		if err != nil {
			log.Printf("❌ failed to unmarshal SQS message: %v", err)
			continue
		}

		// Process the job
		if err := ProcessPrintJob(ctx, printJob); err != nil {
			log.Printf("❌ Error processing print job: %v", err)
			continue
		}
	}
	return nil
}

func main() {
	InitAWS()
	lambda.Start(Handler)
}
