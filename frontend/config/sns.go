package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	s3Client = s3.NewFromConfig(config)
}

// Get photo file names from S3
func getPhotoFileNamesFromS3(customerName, orderID string) ([]string, error) {
	prefix := fmt.Sprintf("%s/%s/", customerName, orderID)
	var fileNames []string

	paginator := s3.NewListObjectsV2Paginator(s3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("listing objects: %w", err)
		}

		for _, obj := range page.Contents {
			key := *obj.Key
			// Extract the filename (after the prefix)
			fileName := key[len(prefix):]
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames, nil
}

// / Modified
func ProcessedPhotoHandler(customerEmail, orderID, customerName string) error {

	// Get photo filenames
	photoFileNames, err := getPhotoFileNamesFromS3(customerEmail, orderID)
	if err != nil {
		return fmt.Errorf("getting photo filenames: %w", err)
	}

	if len(photoFileNames) == 0 {
		log.Println("No photos found for this order")
	}

	for _, photoFileName := range photoFileNames {
		objectKey := fmt.Sprintf("%s/%s/%s", customerName, orderID, photoFileName)
		expires := time.Now().Add(5 * time.Minute).Unix()

		policy := fmt.Sprintf(`{
			"Statement":[{
				"Resource":"https://%s/%s", 
				"Condition":{"DateLessThan":{"AWS:EpochTime":%d}}}]
		}`, cloudFrontDomain, objectKey, expires)

		signedInfo := SignedURLInfo{
			CloudFrontDomain: cloudFrontDomain,
			ObjectKey:        objectKey,
			Expires:          expires,
			KeyPairID:        keyPairID,
			Policy:           policy,
			CustomerName:     customerName,
			OrderID:          orderID,
		}

		message, err := json.Marshal(signedInfo)
		if err != nil {
			return fmt.Errorf("marshaling signed URL info: %w", err)
		}

		err = SendNotification(customerEmail, message)
		if err != nil {
			return fmt.Errorf("sending via SNS: %w", err)
		}
	}

	// // Generate a single signed URL for the folder
	// folderPath := fmt.Sprintf("https://%s/%s/%s", cloudFrontDomain, customerName, orderID)

	// // Generate the signed URL for the entire folder -> TODO: CHANGE!
	// signedURL, err := GenerateSignedURL(folderPath)
	// if err != nil {
	// 	log.Printf("failed to generate signed URL for order %s: %v", orderID, err)
	// 	return err
	// }

	// // Send a single notification
	// return SendNotification(customerEmail, orderID, signedURL)

	return nil
}

// removed photoIDs & signedURL -> review!
// SendNotification sends a notification with the signed URL
func SendNotification(customerEmail string, message []byte) error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	// Get .env variables
	snsTopicArn := os.Getenv("SNS_TOPIC_ARN")

	// Delay sending SNS notification for 10 secs
	log.Println("Delaying SNS notification for 10 seconds...")
	time.Sleep(10 * time.Second)

	// message := "Your photo is ready! You can download it from the following link: " + signedURL

	// Send SNS message to the email endpoint
	_, err = snsClient.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(string(message)),
		Subject:  aws.String("Your photo is ready!"),
		TopicArn: aws.String(snsTopicArn),
	})

	if err != nil {
		log.Printf("failed to send SNS notification message to %s: %v", customerEmail, err)
		return err
	}

	return nil
}
