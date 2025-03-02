package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func SendEmail(ctx context.Context, recipient, body string) error {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS SDK config: %v", err)
	}

	subject := "📸 Your Photo is Ready for Pickup!"
	senderEmail := os.Getenv("SENDER_EMAIL")

	// Validate sender email
	if senderEmail == "" {
		log.Fatal("SENDER_EMAIL environment variable is not set")
	}

	// Initialize SES client
	sesClient := ses.NewFromConfig(cfg)

	// Create the email request
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
		},
		Source: aws.String(senderEmail),
	}

	// Send the email
	_, err = sesClient.SendEmail(ctx, input)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}
