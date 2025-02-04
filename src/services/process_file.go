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
	"github.com/30Piraten/snapflow/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ProcessFile validates and processes a single file. It takes a multipart.FileHeader
// and options for image processing. It first validates the file for security, then
// opens the file and reads its data. Afterwards it checks if the file is above the target
// size and if so, it resizes the image to the target size. The processed image is
// then saved to the uploads directory with a unique filename. It returns a
// FileProcessingResult containing the path of the saved image, the filename and
// size of the original file. If an error occurs during processing, it returns a
// ProcessingError with the appropriate code and message.
func ProcessFile(file *multipart.FileHeader, opts ProcessingOptions, order *PhotoOrder) FileProcessingResult {

	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("BUCKET_NAME")

	// Open the file

	if file == nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "FileError",
				Code:    ErrCodeFileOpen,
				Message: "file is nil",
			},
		}
	}
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
	// cannot use order (variable of type *PhotoOrder) as *fiber.Ctx value in argument to ParseOrderDetailscompiler
	order, err = ParseOrderDetails(order)
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "Validation",
				Code:    ErrCodeProcessingFailed,
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
	userFolder := Sanitize(order.FullName) // R
	uploadDate := time.Now().Format("jan_02")
	uniqueFileName := generateUniqueFileName(file.Filename)
	s3key := path.Join("uploads", userFolder, uploadDate, uniqueFileName)

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

	log.Printf("foldername: %s", userFolder)

	return FileProcessingResult{
		Path:     s3key, // <- Changed the name from s3Path to s3Key
		Filename: file.Filename,
		Size:     file.Size,
	}

}
