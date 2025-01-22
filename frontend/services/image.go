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

	"github.com/nfnt/resize"
	"go.uber.org/zap"
)

// NewImageProcessor creates a new processor instance
func NewImageProcessor(logger *zap.Logger) *ImageProcessor {
	return &ImageProcessor{
		logger: logger,
		cache:  &sync.Map{},
	}
}

func (p *ImageProcessor) ProcessImageWithSizeTarget(originalImage image.Image, opts ProcessingOptions) (image.Image, error) {

	var buf bytes.Buffer

	// Initial compression with high quality
	err := jpeg.Encode(&buf, originalImage, &jpeg.Options{Quality: opts.Quality})
	if err != nil {
		p.logger.Error("Failed to encode image", zap.Error(err))
		return nil, fmt.Errorf("initial encoding failed: %w", err)
	}

	currentSize := int64(buf.Len())

	// Check if the current size already meets the target
	if currentSize <= opts.TargetSizeBytes {
		// Return the compressed image if it meets the target size
		img, _, err := image.Decode(&buf)
		if err != nil {
			p.logger.Info("Image already meets target size",
				zap.Int64("current_size", currentSize),
				zap.Int("quality", opts.Quality),
			)
		}
		return img, nil
	}

	// We must calculate the dimension reduction ratio
	reductionRatio := math.Sqrt(float64(opts.TargetSizeBytes) / float64(currentSize))

	bounds := originalImage.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// Then we update the dimensions
	newWidth := int(float64(originalWidth) * reductionRatio)
	newHeight := int(float64(originalHeight) * reductionRatio)

	// Resize the image
	resizedImage := resize.Resize(uint(newWidth), uint(newHeight), originalImage, resize.Lanczos3)

	// Encode the resized image with reduced quality if necessary
	for {
		buf.Reset() // Clear the buffer for each encoding attempt

		err := jpeg.Encode(&buf, resizedImage, &jpeg.Options{Quality: opts.Quality})
		if err != nil {
			p.logger.Error("Failed to encode resized image", zap.Error(err))
			return nil, fmt.Errorf("resized encoding failed: %w", err)
		}

		finalSize := int64(buf.Len())
		if finalSize <= opts.TargetSizeBytes || opts.Quality <= LowQuality {
			p.logger.Info("Image processing results",
				zap.Int64("original_size", currentSize),
				zap.Int64("final_size", finalSize),
				zap.Float64("reduction_ratio", float64(finalSize)/float64(currentSize)),
				zap.Int("final_quality", opts.Quality),
			)
			break
		}

		// Reduce quality and retry
		opts.Quality -= QualityStep
	}

	img, _, err := image.Decode(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode final image: %w", err)
	}

	return img, nil
}

// SaveImage saves the processed image with appropriate format and quality
func (p *ImageProcessor) SaveImage(img image.Image, path string, opts ProcessingOptions) error {

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
