package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	cfg "github.com/30Piraten/snapflow/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// PhotoOrder represents the structure of o÷ur form data
type PhotoOrder struct {
	FullName  string                  `json:"fullName"`
	Location  string                  `json:"location"`
	Size      string                  `json:"size"`
	PaperType string                  `json:"paperType"`
	Email     string                  `json:"email"`
	Photos    []*multipart.FileHeader `json:"photos"`
}

type PresignedURLResponse struct {
	URL     string            `json:"url"`
	OrderID string            `json:"order_id"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// GeneratePresignedURL generates a presigned URL for the given order details.
// The generated presigned URL is valid for 15 minutes.
// The generated presigned URL will contain the defined metadata
func GeneratePresignedURL(order *PhotoOrder) (*PresignedURLResponse, error) {

	// Initialize DyanmoDB client
	cfg.InitDynamoDB()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	if order.FullName == "" || order.Email == "" {
		return nil, fmt.Errorf("missing required fields for presigned URL generation")
	}

	orderID := uuid.New().String()
	s3key := fmt.Sprintf("%s/%s", order.FullName, orderID)
	uploadTimestamp := time.Now().Unix()

	// Insert metadata into DynamoDB
	err = cfg.InsertMetadata(order.FullName, order.Email, orderID, uploadTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to insert metadata into DynamoDB: %v", err)
	}

	// Initialize S3 client
	s3Client, err := cfg.S3Client()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize s3 client: %v", err)
	}

	// Get bucket name from .envs
	bucketName := os.Getenv("BUCKET_NAME")

	presignedClient := s3.NewPresignClient(s3Client)
	presignedPut, err := presignedClient.PresignPutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(s3key),
			// ACL:    aws.String("private"),
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

	return &PresignedURLResponse{
		URL:     presignedPut.URL,
		OrderID: orderID,
	}, nil
}
