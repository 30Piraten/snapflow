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

// Error returns a string representation of the error in the format "type: code - message".
func (e *ProcessingError) Error() string {
	return fmt.Sprintf("%s: %s - %s", e.Type, e.Code, e.Message)
}

// ProcessUploadedFiles parses the uploaded files and processes them accordingly.
// If there is a single file, it is processed and a JSON response is returned
// containing the file path. If there are multiple files, they are processed
// concurrently and a JSON response is returned containing the paths of the processed
// files. If there were any errors while processing the files, they are collected and returned
// in the response as well. The function returns an error if there were any issues while processing the files.
func ProcessUploadedFiles(c *fiber.Ctx) error {

	// Parse the uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		// return utils.HandleError(c, fiber.StatusBadGateway, "Failed to parse multipart form", err)
		return fmt.Errorf("failed to parse multipart form: %w", err)
	}

	files := form.File["photos"]
	if len(files) == 0 {
		// return utils.HandleError(c, fiber.StatusBadRequest, "No files uploaded", nil)
		return fmt.Errorf("no files uploaded")
	}

	// Validate total file count
	if len(files) > MaxFileCount {
		// return utils.HandleError(c, fiber.StatusBadRequest, fmt.Sprintf("Too many files uploaded. Maximum allowed is %d", MaxFileCount), nil)
		return fmt.Errorf("Too many files uploaded. Maximum allowed is %d", MaxFileCount)
	}

	// Process single or multiple files
	opts := ProcessingOptions{
		Quality:         HighQuality,
		TargetSizeBytes: TargetFileSize,
		Format:          "jpeg",
		MaxDimensions: Dimensions{
			// Need to review: TODO
			Width:  5000,
			Height: 5000,
		},
	}

	// Handle single file
	if len(files) == 1 {
		if err := handleSingleFile(c, files[0], opts); err != nil {
			return fmt.Errorf("failed to process file %s: %w", files[0].Filename, err)
		}
		return nil
	}

	// Validate and handle multiple files
	_, errors := ProcessMultipleFiles(files, opts) // results should be added here
	// If there are errors, return the first one
	if len(errors) > 0 {
		// Collect all errors into one
		var errMsg []string
		for _, e := range errors {
			errMsg = append(errMsg, e.Error())
		}
		// return utils.HandleError(c, fiber.StatusInternalServerError, "Some files failed to process", errors[0])
		return fmt.Errorf("one or more files failed processing: %s", strings.Join(errMsg, "; "))
	}

	// // Prepare the response for successfully processed files
	// var filePaths []string
	// for _, result := range results {
	// 	filePaths = append(filePaths, result.Path)
	// }

	// return c.JSON(fiber.Map{
	// 	"message":   "Files processed succesfully",
	// 	"filePaths": filePaths,
	// })

	return nil
}

// handleSingleFile validates and processes a single file. It first validates the
// file using the provided ProcessingOptions, and if the validation fails, it
// returns a 400 Bad Request error. If the validation succeeds, it processes the
// file and returns a JSON response with the path of the processed file if
// successful. If there are any errors during processing, it returns a 500 Internal
// Server Error with the error message.
func handleSingleFile(c *fiber.Ctx, file *multipart.FileHeader, opts ProcessingOptions) error {

	// Open the file
	source, err := file.Open()
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer source.Close()

	// Read the file data into a byte slice
	fileData, err := io.ReadAll(source)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to read file data", err)
	}

	// Validate the file before processing
	processor := NewImageProcessor(utils.Logger)
	if _, err = processor.ValidateAndProcessImage(fileData, opts); err != nil {
		// utils.Logger.Error("File validation failed", zap.Error(err))
		return utils.HandleError(c, fiber.StatusBadRequest, "File validation failed", err)
	}

	// Process the file
	result := ProcessFile(file, opts)
	if result.Error != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to process file", result.Error)
	}

	return c.JSON(fiber.Map{
		"message":  "File processed successfully",
		"filePath": result.Path,
	})
}

// ProcessMultipleFiles processes multiple uploaded files concurrently using the given
// processing options. It first validates all files and returns any validation errors
// upfront. If all files are valid, they are processed concurrently with a limit on
// the number of concurrent operations. The function returns a slice of
// FileProcessingResult for each successfully processed file and a slice of errors
// for any failed processing attempts. The results and errors are collected and
// returned once all processing is complete.
func ProcessMultipleFiles(files []*multipart.FileHeader, opts ProcessingOptions) ([]FileProcessingResult, []error) {

	var (
		results   []FileProcessingResult
		errors    []error
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, MaxConcurrentProcessing)
	)

	processor := NewImageProcessor(utils.Logger)

	// Validate all files upfront
	for _, file := range files {
		source, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to open file %s: %v", file.Filename, err))
			continue
		}

		fileData, err := io.ReadAll(source)
		source.Close() // close the file after reading
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to read the file %s: %v", file.Filename, err))
			continue
		}

		if _, err := processor.ValidateAndProcessImage(fileData, opts); err != nil {
			errors = append(errors, fmt.Errorf("file %s failed validation: %v", file.Filename, err))
		}
	}

	// Short-circuit if there are validation errors
	if len(errors) > 0 {
		return nil, errors
	}

	// Concurrent processing of validated files
	resultsChan := make(chan FileProcessingResult, len(files))
	errorsChan := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := ProcessFile(file, opts)
			if result.Error != nil {
				errorsChan <- result.Error
			} else {
				resultsChan <- result
			}
		}(file)
	}

	// Collect results and errors
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

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
