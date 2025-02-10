package services

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
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
	"go.uber.org/zap"
)

// ProcessFile validates and processes a file
func ProcessFile(c *fiber.Ctx, file *multipart.FileHeader, opts models.ProcessingOptions, order *models.PhotoOrder) models.FileProcessingResult {

	region, bucketName := os.Getenv("AWS_REGION"), os.Getenv("BUCKET_NAME")
	if region == "" || bucketName == "" {
		utils.Logger.Error("Missing required environment variables")
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "ConfigError",
				Code:    models.ErrCodeInvalidConfig,
				Message: "Env config is empty or nil",
			},
		}
	}

	if file == nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "FileError",
				Code:    models.ErrCodeFileOpen,
				Message: "file is nil",
			},
		}
	}

	// Parse the order details
	parsedOrder, err := ParseOrderDetails(c)
	if err != nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "Validation",
				Code:    models.ErrCodeProcessingFailed,
				Message: fmt.Sprintf("failed to parse form fields: %v", err),
			},
		}
	}

	// order is passed as a pointer, but parsedOrder assigns
	// a new parsed object instead of modifying the existing reference
	// this might not update the original object outside the function
	// hence the pointer reference here!
	*order = *parsedOrder

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

	// Stream file data into a buffer
	bufImg := new(bytes.Buffer)
	_, err = io.Copy(bufImg, source)
	if err != nil {
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "FileError",
				Code:    models.ErrCodeFileRead,
				Message: fmt.Sprintf("failed to read file: %v", err),
			},
		}
	}

	imgData := bufImg.Bytes()

	// Process the image without reading it to memory unnecessarily
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

	// Initialise S3 Client
	s3Config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		utils.Logger.Error("Failed to load AWS config", zap.Error(err))
		return models.FileProcessingResult{
			Error: &models.ProcessingError{
				Type:    "S3ConfigurationError",
				Code:    models.ErrCodeProcessingFailed,
				Message: fmt.Sprintf("failed to load or process AWS config: %v", err),
			},
		}
	}
	s3Client := s3.NewFromConfig(s3Config)

	// Construct the the S3 key with the user's folder and date
	userFolder := utils.Sanitize(order.FullName)
	uploadDate := time.Now().Format("Jan_02")
	uniqueFileName := generateUniqueFileName(file.Filename)
	s3key := path.Join("uploads", userFolder, uploadDate, uniqueFileName)

	// Convert processedImage to []byte
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, processedImage, nil); err != nil {
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

	utils.Logger.Info("Successfully processed and uplaoded image",
		zap.String("folder_name", userFolder),
		zap.String("s3_key", s3key),
		zap.String("file_name", file.Filename),
		zap.Int64("file_size", file.Size),
	)

	return models.FileProcessingResult{
		Path:     s3key,
		Filename: file.Filename,
		Size:     file.Size,
	}
}
