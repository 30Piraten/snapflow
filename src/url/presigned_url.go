package url

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	cfg "github.com/30Piraten/snapflow/config"
	"github.com/30Piraten/snapflow/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// PhotoOrder represents the structure of o÷ur form data

type PresignedURLResponse struct {
	URL     string            `json:"url"`
	OrderID string            `json:"order_id"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// GeneratePresignedURL generates a presigned URL for the given order details.
// The generated presigned URL is valid for 15 minutes.
// The generated presigned URL will contain the defined metadata
func GeneratePresignedURL(order *models.PhotoOrder) (*PresignedURLResponse, error) {

	// Initialize DyanmoDB client
	cfg.InitDynamoDB()

	if order.FullName == "" || order.Email == "" {
		return nil, fmt.Errorf("missing required fields for presigned URL generation")
	}

	orderID := uuid.New().String()
	uploadTimestamp := time.Now().Unix()
	folderKey := fmt.Sprintf("%s/%s", order.FullName, orderID) // -> Review

	// Insert metadata into DynamoDB
	err := cfg.InsertMetadata(order.FullName, order.Email, order.PaperType, order.Size, orderID, uploadTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to insert metadata into DynamoDB: %v", err)
	}

	// Initialize S3 client
	s3Client, err := cfg.S3Client()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize s3 client: %v", err)
	}

	// Get bucket name from .env file
	bucketName := os.Getenv("BUCKET_NAME")

	presignedClient := s3.NewPresignClient(s3Client)

	// Generate presigned URLs for each photo
	var presignedURLs []string
	for _, photo := range order.Photos {
		photoKey := fmt.Sprintf("%s/%s", folderKey, photo.Filename)

		presignedPut, err := presignedClient.PresignPutObject(context.TODO(),
			&s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(photoKey),
				Metadata: map[string]string{
					"full_name":  order.FullName,
					"location":   order.Location,
					"size":       order.Size,
					"paper_type": order.PaperType,
					"order_id":   orderID,
					"email":      order.Email,
				},
			},
			s3.WithPresignExpires(time.Minute*15),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
		}
		presignedURLs = append(presignedURLs, presignedPut.URL)
	}

	// Send SQS print job
	err = cfg.SendPrintRequest(order.Email, orderID, order.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to send SQS print job: %v", err)
	}

	// Send notification via SNS after print job completes
	err = cfg.SendSNSNotification(orderID, order.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to send SNS notification to (email: %s): %v", order.Email, err)
	}

	log.Printf("Values from presignedURL: %s : %s", order.FullName, order.Location)

	return &PresignedURLResponse{
		URL:     folderKey,
		OrderID: orderID,
	}, nil
}
