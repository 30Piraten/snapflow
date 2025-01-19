package routes

import (
	"context"
	"fmt"

	"time"

	cfg "github.com/30Piraten/snapflow/config"
	"github.com/30Piraten/snapflow/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PresignedURLResponse struct {
	URL     string            `json:"url"`
	OrderID string            `json:"order_id"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func Upload(app *fiber.App) {
	app.Post("/generate-upload-url", handleGenerateUploadURL)
}

func handleGenerateUploadURL(c *fiber.Ctx) error {

	// Parse the form data from the webpage into the defined struct
	// PhotoOrder serves as a base for the new "order" struct
	order := new(utils.PhotoOrder)
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	// Validate required fields
	if order.FullName == "" || order.Location == "" || order.Size == "" || order.PaperType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	// Generate unique order ID
	orderID := uuid.New().String()

	// Create the S3 key for the file
	s3key := fmt.Sprintf("%s/%s/%s", order.FullName, order.Email, orderID)

	// Get S3 client
	s3Client, err := cfg.S3Client()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initialize S3 client",
		})
	}

	// Generate presigned URL
	presignClient := s3.NewPresignClient(s3Client)
	presignedPut, err := presignClient.PresignPutObject(context.TODO(),
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate upload URL",
		})
	}

	// Return the presigned URL and order information
	return c.JSON(PresignedURLResponse{
		URL:     presignedPut.URL,
		OrderID: orderID,
	})
}
