package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func Init() {

	// Set up SNS client
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS configuration, %v", err)
	}

	snsClient = sns.NewFromConfig(config)
	s3Client = s3.NewFromConfig(config)
}

type NotificationMessage struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendSNSNotification(photoID, email string) error {
	Init()

	// Replace with your actual SNS topic ARN
	snsTopicArn := os.Getenv("SNS_TOPIC_ARN")

	// Create message payload
	message := NotificationMessage{
		Subject: "Print Job Completed",
		Body:    fmt.Sprintf("Your photo print job (%s) has been completed successfully!", photoID),
	}

	msgBytes, err := json.MarshalIndent(message, "", "  ")
	messageString := string(msgBytes)
	fmt.Println("DEBUG: SQS Message ->", string(msgBytes))

	// Publish to SNS
	resp, err := snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(snsTopicArn),
		Message:  aws.String(messageString),
		Subject:  aws.String("Print Job Completed"),
	})

	if err != nil {
		log.Printf("❌ Failed to send SNS notification: %v", err)
		return err
	}

	log.Printf("✅ SNS Notification sent! Message ID: %s", *resp.MessageId)
	return nil
}
