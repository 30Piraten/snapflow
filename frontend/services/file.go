package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ProcessUploadedFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		utils.Logger.Error("failed to process multipart form: %w", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(ProcessingError{
			Type:    "Validation",
			Code:    ErrCodeInvalidRequest,
			Message: "Failed to parse multipart form",
		})
	}

	files := form.File["photos"]
	if err := ValidateUpload(c, files); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ProcessingError{
			Type:    "Validation",
			Code:    ErrCodeFailedFileUpload,
			Message: "No files uploaded",
		})
	}

	// Create buffered channels for results and limiting concurrency
	results := make(chan FileProcessingResult, len(files))
	semaphore := make(chan struct{}, MaxConcurrentProcessing)

	// Start a worker pool for file processing
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			// Get the semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := processFile(file)
			results <- result
		}(file)
	}

	// Close results channel after all processing is complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and process results
	var errors []string
	for result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", result.Filename, result.Error))
		}
	}

	if len(errors) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": errors,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All files processed successfully",
	})
}

func processFile(file *multipart.FileHeader) FileProcessingResult {

	// Validate file for security
	if err := ValidateUploadedFile(file); err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "Validation",
				Code:    ErrCodeInvalidFormat,
				Message: err.Error(),
			},
		}
	}

	// Open the file
	source, err := file.Open()
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "FileError",
				Code:    ErrCodeFileOpen,
				Message: fmt.Sprintf("failed to open the file: %v", err),
			},
		}
	}
	defer source.Close()

	// Read file data
	img, err := io.ReadAll(source)
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "FileError",
				Code:    ErrCodeFileRead,
				Message: fmt.Sprintf("failed to read file data: %v", err),
			},
		}
	}

	// Initialise processor
	processor := NewImageProcessor(utils.Logger)
	opts := ProcessingOptions{
		Quality:         HighQuality,
		TargetSizeBytes: TargetFileSize,
		Format:          "jpeg",
	}

	// Process the image with size validation
	processedImage, err := processor.ValidateAndProcessImage(img, opts)
	if err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "Processing",
				Code:    ErrCodeProcessingFailed,
				Message: fmt.Sprintf("file to process image: %v", err),
			},
		}
	}

	// Since most photos are often uploaded with the same camera name
	// Lets generate a unique filename to avoid collisions
	// filename := generateUniqueFileName(file.Filename)
	// outputPath := filepath.Join("./uploads", filename)

	// Then we save the processed image
	outputPath := fmt.Sprintf("./%s/%s", ProcessedImageDir, generateUniqueFileName(file.Filename))
	if err := processor.SaveImage(processedImage, outputPath, opts); err != nil {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "FileError",
				Code:    ErrCodeFileSave,
				Message: fmt.Sprintf("failed to save processed image: %v", err),
			},
		}
	}

	return FileProcessingResult{
		Path:     outputPath,
		Filename: file.Filename,
		Size:     file.Size,
	}
}

// Helper function to generate unique filename
func generateUniqueFileName(originalName string) string {
	extension := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, extension)
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s_%s%s", baseName, timestamp, extension)
}
