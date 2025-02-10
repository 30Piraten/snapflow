package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"sync"

	"github.com/30Piraten/snapflow/models"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
)

// ProcessImageWithSizeTarget takes an original image and processes
// it to meet the target size specified in the ProcessingOptions.
func (p *ImageProcessor) ProcessImageWithSizeTarget(originalImage image.Image, opts models.ProcessingOptions) (image.Image, error) {

	var buf bytes.Buffer
	quality := opts.Quality

	// First, encode at lower quality if file size is large
	if err := jpeg.Encode(&buf, originalImage, &jpeg.Options{Quality: quality}); err != nil {
		p.Logger.Error("Failed to encode image", zap.Error(err))
		return nil, fmt.Errorf("encoding failed: %w", err)
	}

	// Get current size
	currentSize := int64(buf.Len())

	// Check if the current size already meets the target
	if currentSize <= opts.TargetSizeBytes {
		img, _, err := image.Decode(&buf)
		if err != nil {
			p.Logger.Error("Failed to decode image after compression", zap.Error(err))
			return nil, fmt.Errorf("failed to decode after compression: %w", err)
		}
		p.Logger.Info("Image already meets target size",
			zap.Int64("current_size", currentSize),
			zap.Int("quality", opts.Quality),
		)
		return img, nil
	}

	// Reduce the quality first before resizing
	for quality > models.LowQuality && currentSize > opts.TargetSizeBytes {
		buf.Reset()
		quality -= models.QualityStep
		if err := jpeg.Encode(&buf, originalImage, &jpeg.Options{Quality: quality}); err != nil {
			return nil, fmt.Errorf("failed during quality reduction: %w", err)
		}
		currentSize = int64(buf.Len())
		p.Logger.Info("Reduced quality",
			zap.Int("quality", quality),
			zap.Int64("current_size", currentSize),
		)
	}

	// Preserve original image before resizing
	resizedImage := originalImage

	// If the file is still too large, resize
	if currentSize > opts.TargetSizeBytes {
		reductionRatio := math.Sqrt(float64(opts.TargetSizeBytes) / float64(currentSize))

		bounds := originalImage.Bounds()
		originalWidth := bounds.Dx()
		originalHeight := bounds.Dy()

		newWidth := int(float64(originalWidth) * reductionRatio)
		newHeight := int(float64(originalHeight) * reductionRatio)

		// Resize with high-quality Lancozos filter
		resizedImage := resize.Resize(uint(newWidth), uint(newHeight), originalImage, resize.Lanczos3)

		// Encode resized image
		buf.Reset()
		if err := jpeg.Encode(&buf, resizedImage, &jpeg.Options{Quality: quality}); err != nil {
			return nil, fmt.Errorf("failed to encode resized image: %w", err)
		}
		p.Logger.Info("Resized image",
			zap.Int("new_width", newWidth),
			zap.Int("new_height", newHeight),
			zap.Float64("reduction_ratio", reductionRatio),
		)
	}

	// Encode the final resized image with reduced quality
	for {
		buf.Reset()
		quality -= models.QualityStep

		if err := jpeg.Encode(&buf, resizedImage, &jpeg.Options{Quality: opts.Quality}); err != nil {
			p.Logger.Error("Failed to encode resized image", zap.Error(err))
			return nil, fmt.Errorf("resized encoding failed: %w", err)
		}

		finalSize := int64(buf.Len())

		p.Logger.Info("Encoding attempt",
			zap.Int64("original_size", currentSize),
			zap.Int64("final_size", finalSize),
			zap.Int("final_quality", quality),
			zap.Float64("reduction_ratio", float64(finalSize)/float64(currentSize)),
		)
		break
	}

	// Return the final decoded image
	img, _, err := image.Decode(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode final image: %w", err)
	}

	return img, nil
}

// SaveImage saves the given image to the specified file path
// using the format and quality options provided in ProcessingOptions.
func (p *ImageProcessor) SaveImage(img image.Image, path string, opts models.ProcessingOptions) error {

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	switch opts.Format {
	case "jpeg", "jpg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: opts.Quality})
	case "png":
		return png.Encode(file, img)
	default:
		return errors.New("unsupported image format")
	}
}

// ConcucurrentProcessImages processes multiple images in parallel
func (p *ImageProcessor) ConcurrentProcessImages(images []image.Image, opts models.ProcessingOptions) []image.Image {

	var wg sync.WaitGroup
	processedImages := make([]image.Image, len(images))
	errors := make([]error, len(images))

	for i, img := range images {
		wg.Add(1)
		go func(index int, img image.Image) {
			defer wg.Done()
			result, err := p.ProcessImageWithSizeTarget(img, opts)
			if err != nil {
				errors[index] = err
				p.Logger.Error("Error processing image", zap.Int("index", i), zap.Error(err))
				return
			}
			processedImages[index] = result
		}(i, img)
	}

	wg.Wait()

	return processedImages
}
