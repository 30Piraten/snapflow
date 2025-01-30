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

	cfg "github.com/30Piraten/snapflow/config"
	"github.com/30Piraten/snapflow/handlers"
	"github.com/30Piraten/snapflow/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

var printjob cfg.PrintJob

// ProcessFile validates and processes a single file. It takes a multipart.FileHeader
// and options for image processing. It first validates the file for security, then
// opens the file and reads its data. Afterwards it checks if the file is above the target
// size and if so, it resizes the image to the target size. The processed image is
// then saved to the uploads directory with a unique filename. It returns a
// FileProcessingResult containing the path of the saved image, the filename and
// size of the original file. If an error occurs during processing, it returns a
// ProcessingError with the appropriate code and message.
func ProcessFile(file *multipart.FileHeader, opts ProcessingOptions, order PhotoOrder) FileProcessingResult {
	// Open the file
	source, err := file.Open()
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "FileError",
				Code:    ErrCodeFileOpen,
				Message: fmt.Sprintf("failed to open file: %v", err),
			},
		}
	}
	defer source.Close()

	// Read file data
	imgData, err := io.ReadAll(source)
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "FileError",
				Code:    ErrCodeFileRead,
				Message: fmt.Sprintf("failed to read file data: %v", err),
			},
		}
	}

	// Validate and process the image
	processor := NewImageProcessor(utils.Logger)
	processedImage, err := processor.ValidateAndProcessImage(imgData, opts)
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "Validation",
				Code:    ErrCodeProcessingFailed,
				Message: fmt.Sprintf("validation or processing failed: %v", err),
			},
		}
	}

	////////////////////////////////////////////////////////////////////
	// Save photos to s3 bucket here
	// Save the processed image
	// outputPath := fmt.Sprintf("./%s/%s", ProcessedImageDir, generateUniqueFileName(file.Filename))

	// Load .env variables and make use of them
	err = godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file", err)
	}

	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("BUCKET_NAME")

	// Initialise S3 Client
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// var order utils.PhotoOrder
	s3Client := s3.NewFromConfig(config)

	// Here we generate a unique file path in S3 under user's folder
	uniqueFileName := generateUniqueFileName(file.Filename)
	s3key := fmt.Sprintf("uploads/%s/%s", order.FullName, uniqueFileName)

	// Convert processedImage to []byte
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, processedImage, nil)
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "ImageEncodingError",
				Code:    ErrCodeProcessingFailed,
				Message: fmt.Sprintf("failed to encode image: %v", err),
			},
		}
	}
	imageBytes := buf.Bytes()

	// Upload processed image to S3
	err = cfg.UploadToS3(s3Client, bucketName, s3key, imageBytes, region)
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "S3Error",
				Code:    ErrCodeFileSave,
				Message: fmt.Sprintf("failed to upload image to S3: %v", err),
			},
		}
	}

	// Simulate PrintJob
	err = handlers.InitiatePrintJob(printjob.CustomerEmail, printjob.PhotoID, printjob.ProcessedS3Location)
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "PrintJobError",
				Code:    ErrCodeProcessingFailed,
				Message: fmt.Sprintf("failed to initiate print job: %v", err),
			},
		}
	}
	////////////////////////////////////////////////////////////////////

	return FileProcessingResult{
		Path:     s3key, // Changed the name from s3Path to s3Key
		Filename: file.Filename,
		Size:     file.Size,
	}
}
