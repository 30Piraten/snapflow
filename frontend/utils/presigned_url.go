package utils

import (
	"context"
	"fmt"
	"time"

	cfg "github.com/30Piraten/snapflow/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PhotoOrder represents the structure of our form data
type PhotoOrder struct {
	FullName  string            `form:"fullName"`
	Location  string            `form:"location"`
	Size      string            `form:"size"`
	PaperType string            `form:"paperType"`
	Email     string            `form:"email"`
	Photos    []*fiber.FormFile `form:"photos"`
}

type PresignedURLResponse struct {
	URL     string            `json:"url"`
	OrderID string            `json:"order_id"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func GeneratePresignedURL(order *PhotoOrder) (*PresignedURLResponse, error) {
	orderID := uuid.New().String()
	s3key := fmt.Sprintf("%s/%s/%s", order.FullName, order.Email, orderID)

	s3Client, err := cfg.S3Client()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize s3 client: %v", err)
	}

	presignedClient := s3.NewPresignClient(s3Client)
	presignedPut, err := presignedClient.PresignPutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String("originalS3"),
			Key:    aws.String(s3key),
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
