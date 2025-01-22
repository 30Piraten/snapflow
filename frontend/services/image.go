package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"mime/multipart"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// NewImageProcessor creates a new processor instance
func NewImageProcessor(logger *zap.Logger) *ImageProcessor {
	return &ImageProcessor{
		logger: logger,
		cache:  &sync.Map{},
	}
}

func ValidateUpload(c *fiber.Ctx, files []*multipart.FileHeader) error {

	// Check the total request size
	if c.Request().Header.ContentLength() > MaxTotalUploadSize {
		return fmt.Errorf("total upload size exceeds %d bytes", MaxTotalUploadSize)
	}

	// Validate file count
	if len(files) > MaxFileCount {
		return fmt.Errorf("too many files upload, max allowed is %d", MaxFileCount)
	}

	// Validate each file size
	for _, file := range files {
		if file.Size > MaxFileSize {
			return fmt.Errorf("file %s exceeds max size of %d bytes", file.Filename, MaxFileSize)
		}
	}

	return nil
}

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

	// Next, if the file || image is between 10MB and 100MB, we set the target size
	if fileSize >= TargetFileSize { // TODO > || >=
		opts.TargetSizeBytes = TargetFileSize
		return p.ProcessImageWithSizeTarget(img, opts)
	}

	return img, nil
}

func (p *ImageProcessor) ProcessImageWithSizeTarget(originalImage image.Image, opts ProcessingOptions) (image.Image, error) {
	var buf bytes.Buffer

	// First we will try an initial compression with high quality
	initialOpts := opts
	initialOpts.Quality = HighQuality

	currentSize := buf.Len()

	// If the file is still too large, calculate reduction ratio and resize
	if int64(currentSize) > opts.TargetSizeBytes {
		// We must calculate the dimension reduction ratio
		reductionRatio := math.Sqrt(float64(opts.TargetSizeBytes) / float64(currentSize))

		bounds := originalImage.Bounds()
		originalWidth := bounds.Dx()
		originalHeight := bounds.Dy()

		// Then we update the dimensions
		opts.MaxWidth = int(float64(originalWidth) * reductionRatio)
		opts.MaxHeight = int(float64(originalHeight) * reductionRatio)

		// If the file || image is still too large, reduce quality
		if opts.Quality > LowQuality {
			opts.Quality = LowQuality
		}
	}

	finalSize := buf.Len()
	p.logger.Info("Image processing results",
		zap.Int64("original_size", int64(currentSize)),
		zap.Int64("final_size", int64(finalSize)),
		zap.Float64("reduction_ratio", float64(finalSize)/float64(currentSize)),
		zap.Int("final_quality", opts.Quality),
	)

	// return img, nil

	return originalImage, nil
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
