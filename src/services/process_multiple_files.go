package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"sync"

	"github.com/30Piraten/snapflow/utils"
)

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

			result := ProcessFile(file, opts, order)
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
