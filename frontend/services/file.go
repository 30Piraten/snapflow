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
)

// Constants for file size
const (
	MaxConcurrentProcessing = 5
)

type ResponseData struct {
	Message      string            `json:"message"`
	Order        *utils.PhotoOrder `json:"order"`
	PresignedURL string            `json:"presigned_url"`
	OrderID      string            `json:"order_id"`
}

// FileProcessingResult holds the result of processing a single file
type FileProcessingResult struct {
	Filename string
	Path     string
	Size     int64
	Error    error
}

func ProcessUploadedFiles(c *fiber.Ctx, order *utils.PhotoOrder) error {
	form, err := c.MultipartForm()
	if err != nil {
		return fmt.Errorf("failed to process multipart form: %w", err)
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return fmt.Errorf("no files uploaded")
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
		return fmt.Errorf("file processing errors: %v", errors)
	}

	return nil
}

func processFile(file *multipart.FileHeader) FileProcessingResult {

	// Check the file size first
	if file.Size > MaxFileSize {
		return FileProcessingResult{
			Error: fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", file.Size, MaxFileSize),
		}
	}

	// Initialise processor
	processor := NewImageProcessor(utils.Logger)

	// Define processing options with target size for large files
	opts := ProcessingOptions{
		Quality: HighQuality,
		Format:  "jpeg",
	}

	// set the target if file is larger than 10MB
	if file.Size > TargetFileSize {
		opts.TargetSizeBytes = TargetFileSize
	}

	// Open the file
	source, err := file.Open()
	if err != nil {
		return FileProcessingResult{
			Error: fmt.Errorf("failed to open the file: %w", err),
		}
	}
	defer source.Close()

	// Decode file data
	// img, _, err := image.Decode(source)
	// if err != nil {
	// 	return FileProcessingResult{
	// 		Error: err,
	// 	}
	// }

	// Read file data
	img, err := io.ReadAll(source)
	if err != nil {
		return FileProcessingResult{
			Error: fmt.Errorf("failed to read file data: %w", err),
		}
	}

	// Process the image with size validation
	processedImage, err := processor.ValidateAndProcessImage(img, opts)
	if err != nil {
		if strings.Contains(err.Error(), "exceeds maximum allowed size") {
			return FileProcessingResult{
				Error: fmt.Errorf("file too large: %w", err),
			}
		}
		return FileProcessingResult{
			Error: fmt.Errorf("filed to process image: %w", err),
		}
	}

	// // Create uploads directory if it doesn't exist
	// if err := os.MkdirAll("./uploads", 0775); err != nil {
	// 	return FileProcessingResult{
	// 		Error: fmt.Errorf("failed to create uploads directory: %w", err),
	// 	}
	// }

	// Since most photos are often uploaded with the same camera name
	// Lets generate a unique filename to avoid collisions
	filename := generateUniqueFileName(file.Filename)
	// outputPath := filepath.Join("./uploads", filename)

	// Then we save the processed image
	outputPath := fmt.Sprintf("./uploads/%s", filename) // file.Filename
	if err := processor.SaveImage(processedImage, outputPath, opts); err != nil {
		return FileProcessingResult{
			Error: fmt.Errorf("failed to save processed image: %w", err),
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
