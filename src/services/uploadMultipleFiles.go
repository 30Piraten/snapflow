package services

import (
	"fmt"
	"io"
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
			errors = append(errors, fmt.Errorf("failed to open file %s: %v", file.Filename, err))
			continue
		}

		fileData, err := io.ReadAll(source)
		source.Close() // Close the file after reading
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to read the file %s: %v", file.Filename, err))
			continue
		}
		processor := NewImageProcessor(utils.Logger)
		if _, err := processor.ValidateAndProcessImage(fileData, opts); err != nil {
			errors = append(errors, fmt.Errorf("file %s failed validation: %v", file.Filename, err))
		}
	}

	// Short-circuit if the validation fials
	if len(errors) > 0 {
		return nil, errors
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

	wg.Wait()
	close(resultsChan)
	close(errorsChan)

	// Collect results and errors
	for result := range resultsChan {
		results = append(results, result)
	}
	for err := range errorsChan {
		errors = append(errors, err)
	}

	return results, errors
}
