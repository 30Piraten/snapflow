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

// NewImageProcessor creates a new instance of ImageProcessor with the provided logger.
// The logger is used for logging messages, and a new cache is initialized for caching processed images.
func NewImageProcessor(logger *zap.Logger) *ImageProcessor {
	return &ImageProcessor{
		logger: logger,
		cache:  &sync.Map{},
	}
}

// ProcessImageWithSizeTarget takes an original image and processes it to meet the target size specified in
// the ProcessingOptions. If the original image is already below the target size, it is returned as is.
// Otherwise, the image is resized to the target size and then encoded with a quality that is reduced
// incrementally until the target size is met. The final quality of the image is returned in the
// ProcessingOptions struct.
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

		// Display finalSize to terminal
		// utils.Logger.Info("final_size:", zap.String("final_size", fmt.Sprintf("%d", finalSize)))
	}

	img, _, err := image.Decode(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode final image: %w", err)
	}

	return img, nil
}

// SaveImage saves the given image to the specified file path using the format and quality
// options provided in ProcessingOptions. The function supports JPEG and PNG formats.
// An error is returned if the file cannot be created or if the image format is unsupported.
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
