package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var snsClient *sns.Client

func InitSNS() {
	// Decided to set an explicit region!
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Failed to load SNS config: %v", err)
	}
	snsClient = sns.NewFromConfig(cfg)
}

func SendSNSNotification(orderID, orderEmail string) error {
	if snsClient == nil {
		InitSNS()
	}

	snsTopicArn := os.Getenv("SNS_TOPIC_ARN")
	if snsTopicArn == "" {
		return fmt.Errorf("SNS_TOPIC_ARN environment variable not set")
	}

	message := fmt.Sprintf("Your order (ID: %s) has been processed! Thank you for your order. We will notify you when your order is ready for pickup.", orderID) // Improved message

	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(snsTopicArn),
		Subject:  aws.String("Order Update"),
	}

	result, err := snsClient.Publish(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("failed to send SNS notification: %w", err)
	}

	log.Printf("Successfully sent SNS notification: %v", result.MessageId)
	return nil
}
