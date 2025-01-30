package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SQSEvent struct {
	Records []struct {
		Body string `json:"body"`
	} `json:"Records"`
}

var snsTopicArn string

func InitAWS() {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	snsClient = sns.NewFromConfig(config)

	// Get SNS topic ARN from environment variable
	snsTopicArn = os.Getenv("SNS_TOPIC_ARN")
}

func Handler(ctx context.Context, event SQSEvent) {
	InitAWS()

	for _, record := range event.Records {
		message := record.Body

		_, err := snsClient.Publish(context.TODO(), &sns.PublishInput{
			Message:  aws.String(message),
			TopicArn: aws.String(snsTopicArn),
		})

		if err != nil {
			log.Fatalf("failed to send SNS message, %v", err)
		} else {
			log.Println("SNS notification sent successfully")
		}
	}
}
