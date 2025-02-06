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

	var (
		results   []models.FileProcessingResult
		errors    []error
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, models.MaxConcurrentProcessing)
	)

	processor := NewImageProcessor(utils.Logger)
	order := new(models.PhotoOrder)

	// Validate all files upfront
	for _, file := range files {
		source, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to open file %s: %v", file.Filename, err))
			continue
		}

		fileData, err := io.ReadAll(source)
		// Close the file after reading
		source.Close()
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
	resultsChan := make(chan models.FileProcessingResult, len(files))
	errorsChan := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := ProcessFile(c, file, opts, order)
			if result.Error != nil {
				errorsChan <- &result.Error.Error // Review
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
