package services

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

// Error returns a string representation of the error in the format "type: code - message".
func (e *ProcessingError) Error() string {
	return fmt.Sprintf("%s: %s - %s", e.Type, e.Code, e.Message)
}

// ProcessUploadedFiles parses the uploaded files and processes them accordingly.
// If there is a single file, it is processed and a JSON response is returned
// containing the file path.
// If there are multiple files, they are processed concurrently and a JSON
// response is returned containing the paths of the processed files. If there
// were any errors while processing the files, they are collected and returned
// in the response as well.
// The function returns an error if there were any issues while processing the
// files.
func ProcessUploadedFiles(c *fiber.Ctx) error {
	// Parse the uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadGateway, "Failed to parse multipart form", err)
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return utils.HandleError(c, fiber.StatusBadRequest, "No files uploaded", nil)
	}

	// Confirm whether its a single file or multiple files
	opts := ProcessingOptions{
		Quality:         HighQuality,
		TargetSizeBytes: TargetFileSize,
		Format:          "jpeg",
	}

	// Process single file
	if len(files) == 1 {
		result := ProcessFile(files[0], opts)
		if result.Error != nil {
			// return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to process file", result.Error)
			return &ProcessingError{
				Type:    "Validation",
				Code:    ErrCodeProcessingFailed,
				Message: "Failed to process file",
			}
		}

		return c.JSON(fiber.Map{
			"message":  "File processed successfully",
			"filePath": result.Path,
		})
	}

	// Process multiple files
	results, errors := ProcessMultipleFiles(files, opts)
	if len(errors) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some files failed to process",
			"errors":  errors,
		})
	}

	// Collect the results for the response
	var processedFiles []string
	for _, result := range results {
		processedFiles = append(processedFiles, result.Path)
	}

	return c.JSON(fiber.Map{
		"message":        "Files processed succesfully",
		"processedFiles": processedFiles,
	})
}

// ProcessFile validates and processes a single file. It takes a multipart.FileHeader
// and options for image processing. It first validates the file for security, then
// opens the file and reads its data. Afterwards it checks if the file is above the target
// size and if so, it resizes the image to the target size. The processed image is
// then saved to the uploads directory with a unique filename. It returns a
// FileProcessingResult containing the path of the saved image, the filename and
// size of the original file. If an error occurs during processing, it returns a
// ProcessingError with the appropriate code and message.
func ProcessFile(file *multipart.FileHeader, opts ProcessingOptions) FileProcessingResult {

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

	// If file exceeds 50MB, reject it
	if file.Size > MaxFileSize {
		return FileProcessingResult{
			Error: &ProcessingError{
				Type:    "Processing",
				Code:    ErrCodeFileTooLarge,
				Message: fmt.Sprintf("file size exceeds the 50MB limit: %v bytes", file.Size),
			},
		}
	}

	var processedImage image.Image
	// If the file is large (1MB - 50MB), we resize it to <1MB with strict setting
	if file.Size > TargetFileSize && file.Size <= MaxFileSize {
		opts.TargetSizeBytes = TargetFileSize
		// Process the image with size validation
		processedImage, err = processor.ValidateAndProcessImage(img, opts)
		if err != nil {
			return FileProcessingResult{
				Error: &ProcessingError{
					Type:    "Processing",
					Code:    ErrCodeProcessingFailed,
					Message: fmt.Sprintf("file to process image: %v", err),
				},
			}
		}

	} else {
		// If the file is already under 1MB, there's no need to resize
		processedImage, err = processor.ValidateAndProcessImage(img, opts)
		if err != nil {
			return FileProcessingResult{
				Error: &ProcessingError{
					Type:    "Processing",
					Code:    ErrCodeProcessingFailed,
					Message: fmt.Sprintf("failed to process image: %v", err),
				},
			}
		}
	}

	// Since most photos are often uploaded with the same camera name
	// Lets generate a unique filename to avoid collisions
	// filename := generateUniqueFileName(file.Filename)
	// outputPath := filepath.Join("./uploads", filename)

	// Then we save the processed image with a unique name
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

// ProcessMultipleFiles processes multiple files concurrently and returns the results
// and any errors encountered while processing the files.
//
// The function will process the files in chunks of MaxConcurrentProcessing, so if
// you have 10 files and MaxConcurrentProcessing is 3, the function will process
// 3 files at a time until all files are processed.
//
// The function returns two channels, one for the results and one for any errors
// encountered while processing the files. The results channel will contain all
// the results of the processed files, and the errors channel will contain all the
// errors encountered while processing the files.
//
// The function will block until all files have been processed and the results and
// errors have been collected.
func ProcessMultipleFiles(files []*multipart.FileHeader, opts ProcessingOptions) ([]FileProcessingResult, []error) {

	var (
		results   []FileProcessingResult
		errors    []error
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, MaxConcurrentProcessing)
	)

	// Create channels for results
	resultsChan := make(chan FileProcessingResult, len(files))
	errorsChan := make(chan error, len(files))

	// Iterate over all files and process them concurrently
	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader) {
			defer wg.Done()
			// Acquire semaphore for concurrency control
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			// Process the file
			result := ProcessFile(file, opts)
			if result.Error != nil {
				errorsChan <- result.Error
			} else {
				resultsChan <- result
			}
		}(file)
	}

	go func() {
		// Wait for all goroutines to finish
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	// Collect results and errors
	for result := range resultsChan {
		results = append(results, result)
	}

	for err := range errorsChan {
		errors = append(errors, err)
	}

	return results, errors
}

// generateUniqueFileName generates a unique filename by appending a timestamp
// to the base name of the original file name.
func generateUniqueFileName(originalName string) string {
	extension := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, extension)
	timestamp := time.Now().Format("20060102150405") // -> Make random
	return fmt.Sprintf("%s_%s%s", baseName, timestamp, extension)
}
