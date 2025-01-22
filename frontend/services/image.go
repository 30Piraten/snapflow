package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"runtime"
	"sync"

	"go.uber.org/zap"
)

// Quality settings for image processing
const (
	HighQuality    = 90
	MediumQuality  = 75
	LowQuality     = 60
	MaxFileSize    = 100 * 100 * 1024 // 100MB in bytes
	TargetFileSize = 10 * 1024 * 1024 // 10MB in bytes
)

// ProcessingOptions defines configuration for image processing
type ProcessingOptions struct {
	MaxWidth         int
	MaxHeight        int
	Quality          int
	Sharpen          bool
	Format           string
	PreserveMetadata bool
	OptimiseSizeOnly bool
	TargetSizeBytes  int64 // -> New field for target file size
}

// ImageProcessor handles all image processing operations
type ImageProcessor struct {
	logger *zap.Logger
	cache  *sync.Map // -> Allow cache for processed images
}

// NewImageProcessor creates a new processor instance
func NewImageProcessor(logger *zap.Logger) *ImageProcessor {
	return &ImageProcessor{
		logger: logger,
		cache:  &sync.Map{},
	}
}

// ProcessImage performs all necessary image processing operations
func (p *ImageProcessor) ProcessImage(img image.Image, opts ProcessingOptions) (image.Image, error) {

	// Start with resizing if needed
	if opts.MaxWidth > 0 || opts.MaxHeight > 0 {
		resized, err := p.resizeLanczos(img, opts.MaxWidth, opts.MaxHeight)
		if err != nil {
			return nil, fmt.Errorf("resize error: %w", err)
		}
		img = resized
	}

	// Apply sharpening if requested
	if opts.Sharpen {
		img = p.sharpenImage(img)
	}

	// Optimise size if requested
	if opts.OptimiseSizeOnly {
		img = p.OptimiseSize(img)
	}

	return img, nil
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

	// For files || images under the 10MB mark, we just process normally
	return p.ProcessImage(img, opts)
}

func (p *ImageProcessor) ProcessImageWithSizeTarget(originalImage image.Image, opts ProcessingOptions) (image.Image, error) {

	// First try: just compression
	var buf bytes.Buffer
	img, err := p.ProcessImage(originalImage, opts)
	if err != nil {
		return nil, err
	}

	if opts.Format == "jpeg" {
		err = jpeg.Encode(&buf, img, &jpeg.Options{
			Quality: opts.Quality,
		})
	} else {
		err = png.Encode(&buf, img)
	}

	if err != nil {
		return nil, err
	}

	currentSize := buf.Len()

	// if we're already under target size, return
	if int64(currentSize) <= opts.TargetSizeBytes {
		return img, nil
	}

	// Calculate necessary reduction ratio
	reductionRatio := float64(opts.TargetSizeBytes) / float64(currentSize)

	// Adjust dimensions to meet target size while preserving aspect ratio
	bounds := originalImage.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// Calculate new dimensions based on reduction ratio
	newWidth := int(math.Sqrt(reductionRatio) * float64(originalWidth))
	newHeight := int(math.Sqrt(reductionRatio) * float64(originalHeight))

	// Update options with new dimensions
	opts.MaxHeight = newHeight
	opts.MaxWidth = newWidth

	// Try processing with new dimensions
	img, err = p.ProcessImage(originalImage, opts)
	if err != nil {
		return nil, err
	}

	// Verify final size
	buf.Reset()
	if opts.Format == "jpeg" {
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: opts.Quality})
	} else {
		err = png.Encode(&buf, img)
	}

	if err != nil {
		return nil, err
	}

	finalSize := buf.Len()
	p.logger.Info("Image processing results",
		zap.Int("original_size", currentSize),
		zap.Int("final_size", finalSize),
		zap.Float64("reduction_ratio", float64(finalSize)/float64(currentSize)),
	)

	return img, nil
}

// resizeLanczos implements high-quality Lanczos resampling
func (p *ImageProcessor) resizeLanczos(img image.Image, maxWidth, maxHeight int) (image.Image, error) {

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate new dimensions while maintaining aspect ratio
	ratio := float64(width) / float64(height)
	if maxWidth > 0 && maxHeight > 0 {
		newRatio := float64(maxWidth) / float64(maxHeight)
		if ratio > newRatio {
			maxHeight = int(float64(maxWidth) / ratio)
		} else {
			maxWidth = int(float64(maxHeight) * ratio)
		}
	}

	// Create output image
	destination := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))

	// Lanczos kernel
	kernel := func(x float64) float64 {
		if x == 0 {
			return 1
		}
		if x < -3 || x > 3 {
			return 0
		}
		return (math.Sin(math.Pi*x) * math.Sin(math.Pi*x/3)) / (math.Pi * math.Pi * x * x / 3)
	}

	// Process each pixel with Lanczos algorithm
	scaleX := float64(width) / float64(maxWidth)
	scaleY := float64(height) / float64(maxHeight)

	// Use parallel processing for better performance
	var wg sync.WaitGroup
	workers := runtime.NumCPU()
	rowsPerWorker := maxHeight / workers

	for w := 0; w < workers; w++ {
		wg.Add(1)
		startY := w * rowsPerWorker
		endY := startY + rowsPerWorker
		if w == workers-1 {
			endY = maxHeight
		}

		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := 0; x < maxWidth; x++ {
					// Calculate source position
					sourceX := (float64(x) + 0.5) * scaleX
					sourceY := (float64(y) + 0.5) * scaleY

					// Accumalate weighted colors
					var r, g, b, a float64
					var totalWeight float64

					// Sample neighbouring pixels
					for ix := -3; ix <= 3; ix++ {
						for iy := -3; iy <= 3; iy++ {
							px := int(sourceX) + ix
							py := int(sourceY) + iy

							if px < 0 || px >= width || py < 0 || py >= height {
								continue
							}

							// Calculate weight using Lanczos kernel
							dx := sourceX - float64(px)
							dy := sourceY - float64(py)
							weight := kernel(dx) * kernel(dy)

							// Get source color
							sourceColour := img.At(px, py)
							sr, sg, sb, sa := sourceColour.RGBA()

							// Accumalate weighted components
							r += float64(sr) * weight
							g += float64(sg) * weight
							b += float64(sb) * weight
							a += float64(sa) * weight
							totalWeight += weight
						}
					}

					// Normalise and set pixel
					if totalWeight != 0 {
						r /= totalWeight
						g /= totalWeight
						b /= totalWeight
						a /= totalWeight
					}

					destination.Set(x, y, color.RGBA64{
						R: uint16(r),
						G: uint16(g),
						B: uint16(b),
						A: uint16(a),
					})
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return destination, nil
}

// sharpenImage applies an unsharp mask for image sharpening
func (p *ImageProcessor) sharpenImage(img image.Image) image.Image {
	bounds := img.Bounds()
	destination := image.NewNRGBA(bounds)

	// Sharpening kernel
	kernel := [][]float64{
		{-1, -1, -1},
		{-1, 9, -1},
		{-1, -1, -1},
	}

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var r, g, b float64

			// Apply convolution
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					sourceColour := img.At(x+kx, y+ky)
					sr, sg, sb, _ := sourceColour.RGBA()
					weight := kernel[ky+1][kx+1]

					r += float64(sr) * weight
					g += float64(sg) * weight
					b += float64(sb) * weight
				}
			}

			// Clamp values
			r = math.Max(0, math.Min(65535, r))
			g = math.Max(0, math.Min(65535, g))
			b = math.Max(0, math.Min(65535, b))

			destination.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: 65535,
			})
		}
	}
	return destination
}

// optimiseSize reduces image size while maintaining quality
func (p *ImageProcessor) OptimiseSize(img image.Image) image.Image {
	bounds := img.Bounds()
	destination := image.NewRGBA(bounds)

	// Simple colour quantization
	colourMap := make(map[color.Color]color.Color)
	threshold := uint32(2000) // -> Color difference threshold

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			sourceColour := img.At(x, y)
			if mappedColour, exists := colourMap[sourceColour]; exists {
				destination.Set(x, y, mappedColour)
				continue
			}

			// Find similar color in map
			found := false
			for existingColour, mappedColour := range colourMap {
				if colourDifference(sourceColour, existingColour) < threshold {
					colourMap[sourceColour] = mappedColour
					destination.Set(x, y, mappedColour)
					found = true
					break
				}
			}

			if !found {
				colourMap[sourceColour] = sourceColour
				destination.Set(x, y, sourceColour)
			}
		}
	}

	return destination
}

// colourDifference calculates the difference between two colours
func colourDifference(colour1, colour2 color.Color) uint32 {
	r1, g1, b1, _ := colour1.RGBA()
	r2, g2, b2, _ := colour2.RGBA()

	return abs(r1-r2) + abs(g1-g2) + abs(b1-b2)
}

func abs(x uint32) uint32 {
	if x < 0 {
		return -x
	}

	return x
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
