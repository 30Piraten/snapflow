package services

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"time"

	cfg "github.com/30Piraten/snapflow/config"
	"github.com/30Piraten/snapflow/models"
	"github.com/30Piraten/snapflow/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// ProcessFile validates and processes a single file
func ProcessFile(c *fiber.Ctx, file *multipart.FileHeader, opts models.ProcessingOptions, order *models.PhotoOrder) models.FileProcessingResult {

	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("BUCKET_NAME")

	if file == nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "FileError",
				Code:    models.ErrCodeFileOpen,
				Message: "file is nil",
			},
		}
	}

	// Open the file
	source, err := file.Open()
	if err != nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "FileError",
				Code:    models.ErrCodeFileOpen,
				Message: fmt.Sprintf("failed to open file: %v", err),
			},
		}
	}
	defer source.Close()

	// Read file data
	imgData, err := io.ReadAll(source)
	if err != nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "FileError",
				Code:    models.ErrCodeFileRead,
				Message: fmt.Sprintf("failed to read file data: %v", err),
			},
		}
	}

	// Validate and process the image
	processor := NewImageProcessor(utils.Logger)
	processedImage, err := processor.ValidateAndProcessImage(imgData, opts)
	if err != nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "Validation",
				Code:    models.ErrCodeProcessingFailed,
				Message: fmt.Sprintf("validation or processing failed: %v", err),
			},
		}
	}

	// Parse the order details
	order, err = ParseOrderDetails(c)
	if err != nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "Validation",
				Code:    models.ErrCodeProcessingFailed,
				Message: fmt.Sprintf("failed to parse form fields: %v", err),
			},
		}
	}

	// Initialise S3 Client
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(config)

	// Construct the the S3 key with the user's folder and date
	userFolder := utils.Sanitize(order.FullName)
	uploadDate := time.Now().Format("jan_02")
	uniqueFileName := generateUniqueFileName(file.Filename)
	s3key := path.Join("uploads", userFolder, uploadDate, uniqueFileName)

	// Convert processedImage to []byte
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, processedImage, nil)
	if err != nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "ImageEncodingError",
				Code:    models.ErrCodeProcessingFailed,
				Message: fmt.Sprintf("failed to encode image: %v", err),
			},
		}
	}

	imageBytes := buf.Bytes()

	// Upload processed image to S3
	err = cfg.UploadToS3(s3Client, bucketName, s3key, imageBytes, region)
	if err != nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "S3Error",
				Code:    models.ErrCodeFileSave,
				Message: fmt.Sprintf("failed to upload image to S3: %v", err),
			},
		}
	}

	log.Printf("foldername: %s", userFolder)

	return models.FileProcessingResult{
		Path:     s3key,
		Filename: file.Filename,
		Size:     file.Size,
	}
}
