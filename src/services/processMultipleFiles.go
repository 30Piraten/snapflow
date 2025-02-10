package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"sync"

	"github.com/30Piraten/snapflow/models"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

// ProcessMultipleFiles processes multiple uploaded files
// concurrently using the given processing options.
func ProcessMultipleFiles(c *fiber.Ctx, files []*multipart.FileHeader, opts models.ProcessingOptions) ([]models.FileProcessingResult, []error) {

	// Define required variables
	var (
		results   []models.FileProcessingResult
		errors    []error
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, models.MaxConcurrentProcessing)
	)

	resultsChan := make(chan models.FileProcessingResult, len(files))
	errorsChan := make(chan error, len(files))

	// Validate all files upfront
	for _, file := range files {
		source, err := file.Open()
		if err != nil {
			errorsChan <- fmt.Errorf("failed to open file %s: %v", file.Filename, err)
			continue
		}
		defer source.Close() // Close the file after reading

		// fileData, err := io.ReadAll(source)
		fileData, _, err := image.Decode(source)
		if err != nil {
			errorsChan <- fmt.Errorf("file %s failed decoding: %v", file.Filename, err)
			continue
		}

		// Convert back to []byte
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, fileData, nil); err != nil {
			errorsChan <- fmt.Errorf("file %s failed encoding; %v", file.Filename, err)
			continue
		}

		processor := NewImageProcessor(utils.Logger)
		if _, err := processor.ValidateAndProcessImage(buf.Bytes(), opts); err != nil {
			errorsChan <- fmt.Errorf("file %s failed validation: %v", file.Filename, err)
		}
	}

	// // Short-circuit if the validation fails
	close(errorsChan)

	var validationErrors []error
	for err := range errorsChan {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	// Concurrent processing of validated files
	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Create a separate PhotoOrder for each files
			order, err := ParseOrderDetails(c)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to parse form fields for file %s: %v", file.Filename, err)
				return
			}

			// Pass the order to ProcessFile
			result := ProcessFile(c, file, opts, order)
			if result.Error != nil {
				errorsChan <- fmt.Errorf("file: %s processing failed %v", file.Filename, result.Error)
			} else {
				resultsChan <- result
			}
		}(file)
	}

	// Close channels only after all goroutines finish
	go func() {
		wg.Wait()
		close(resultsChan)
		// close(errorsChan)
	}()

	// Collect results and errors
	for result := range resultsChan {
		results = append(results, result)
	}
	// var processErrors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	return results, errors
}
