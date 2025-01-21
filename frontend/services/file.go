package services

import (
	"fmt"
	"image"
	"mime/multipart"
	"sync"

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
	Error    error
	Path     string
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

	// Initialise processor
	processor := NewImageProcessor(utils.Logger)

	// Define processing options
	opts := ProcessingOptions{
		Quality: HighQuality,
		Format:  "jpeg",
	}

	// Open and decode the image
	source, err := file.Open()
	if err != nil {
		return FileProcessingResult{
			Error: err,
		}
	}
	defer source.Close()

	img, _, err := image.Decode(source)
	if err != nil {
		return FileProcessingResult{
			Error: err,
		}
	}

	// Process the image
	processed, err := processor.ProcessImage(img, opts)
	if err != nil {
		return FileProcessingResult{
			Error: err,
		}
	}

	// Save processed image
	outputPath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := processor.SaveImage(processed, outputPath, opts); err != nil {
		return FileProcessingResult{
			Error: err,
		}
	}

	return FileProcessingResult{
		Path:     outputPath,
		Filename: file.Filename,
	}
}
