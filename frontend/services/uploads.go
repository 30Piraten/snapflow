package services

import (
	"bytes"
	"fmt"
	"image"
	"sync"
)

func (p *ImageProcessor) ValidateAndProcessImage(imgData []byte, opts ProcessingOptions) (image.Image, error) {

	// We must check the file size first before validating or processing
	fileSize := int64(len(imgData))

	if fileSize > MaxFileSize {
		return nil, fmt.Errorf("file szie %d bytes exceeds maximum allowed szie of %d bytes", fileSize, MaxFileSize)
	}

	// Next we decode the image
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// We set the format of the image if it's not specified || unknown
	if opts.Format == "" {
		opts.Format = format
	}

	// Next, if the file is between 1MB and 100MB, we set the target size
	if fileSize > TargetFileSize {
		opts.TargetSizeBytes = TargetFileSize
		return p.ProcessImageWithSizeTarget(img, opts)
	}

	return img, nil
}

func (p *ImageProcessor) ValidateAndProcessMultipleImages(files [][]byte, opts ProcessingOptions) ([]image.Image, error) {
	// Security: Ensure the request adheres to the file limit
	if len(files) > MaxFileCount {
		return nil, fmt.Errorf("too many files: %d exceeds maximum allowed %d", len(files), MaxFileCount)
	}

	var (
		processedImages = make([]image.Image, len(files))
		totalSize       int64
		mu              sync.Mutex
		wg              sync.WaitGroup
		errChan         = make(chan error, len(files))
	)

	// Calculate total size and validate each file
	for _, fileData := range files {
		fileSize := int64(len(fileData))
		if fileSize > MaxFileSize {
			return nil, fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", fileSize, MaxFileSize)
		}
		totalSize += fileSize
	}

	// Goroutine logic for resizing or accepting files
	for idx, fileData := range files {
		wg.Add(1)
		go func(i int, data []byte) {
			defer wg.Done()

			img, err := p.ValidateAndProcessImage(data, opts)
			if err != nil {
				errChan <- fmt.Errorf("error processing file %d: %w", i+1, err)
				return
			}

			// Strict resizing for total size > 1MB
			if totalSize > MaxTotalUploadSize {
				opts.TargetSizeBytes = TargetFileSize
				img, err = p.ProcessImageWithSizeTarget(img, opts)
				if err != nil {
					errChan <- fmt.Errorf("error resizing file %d: %w", i+1, err)
					return
				}
			}

			mu.Lock()
			processedImages[i] = img
			mu.Unlock()
		}(idx, fileData)
	}

	wg.Wait()
	close(errChan)

	// Check for any errors during processing
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return processedImages, nil
}
