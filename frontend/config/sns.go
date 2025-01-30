package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/joho/godotenv"
)

func init() {
	InitCloudFront()

	// Set up SNS client
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS configuration, %v", err)
	}

	snsClient = sns.NewFromConfig(config)
}

// SendNotification sends a notification with the signed URL
func SendNotification(customerEmail, photoIDs, signedURL string) error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	// Get .env variables
	snsTopicArn := os.Getenv("SNS_TOPIC_ARN")

	message := "Your photo is ready! You can download it from the following link: " + signedURL

	// Send SNS message to the email endpoint
	_, err = snsClient.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(message),
		Subject:  aws.String("Your photo is ready!"),
		TopicArn: aws.String(snsTopicArn),
	})

	if err != nil {
		log.Printf("failed to send SNS notification message to %s: %v", customerEmail, err)
		return err
	}

	return nil
}

// / Modified
func ProcessedPhotoHandler(customerEmail string, orderID string, customerName string) error {

	// Generate a single signed URL for the folder
	folderPath := fmt.Sprintf("https://%s/%s/%s", cloudFrontDomain, customerName, orderID)

	// Generate the signed URL for the entire folder
	signedURL, err := GenerateSignedURL(folderPath)
	if err != nil {
		log.Printf("failed to generate signed URL for order %s: %v", orderID, err)
		return err
	}

	// Send a single notification
	return SendNotification(customerEmail, orderID, signedURL)
}

// // Simulate a delay before sending the SNS notification
// func DelayedSendNotification(email string, message string) {
// 	log.Println("Delaying SNS notification for 10 seconds...")
// 	time.Sleep(10 * time.Second) // Delay for 10 seconds
// 	err := SendNotification(email, message)
// 	if err != nil {
// 		log.Printf("Failed to send SNS notification: %v", err)
// 	}
// }

// ProcessedPhotoHandler handles photo processing and sends notification
// func ProcessedPhotoHandler(customerEmail string, photoIDs []string, processedS3Location string) error {
// 	// Generate the CloudFront signed URL for the processed photo
// 	signedURLs := make([]string, len(photoIDs))
// 	for i, id := range photoIDs {
// 		signedURL, err := GenerateSignedURL(processedS3Location + "/" + id)
// 		if err != nil {
// 			log.Printf("failed to generate signed URL for photo %s: %v", id, err)
// 			return err
// 		}
// 		signedURLs[i] = signedURL
// 	}

// 	// Convert array of URLS into a single message
// 	fullMessage := strings.Join(signedURLs, "\n")

// 	// Send the notification email with signed URL
// 	return SendNotification(customerEmail, photoIDs, fullMessage)
// }
