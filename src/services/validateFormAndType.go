package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"net/http"
	"regexp"
	"strings"

	"github.com/30Piraten/snapflow/models"
)

// AllowedFileExtensions defines permitted image file extensions
var AllowedFileExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
}

// ValidateOrder validates a PhotoOrder instance,
// returning an error if required fields are missing.
func ValidateOrder(order *models.PhotoOrder) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	var missingFields []string

	if strings.TrimSpace(order.FullName) == "" {
		return errors.New("Full name is required")
	}
	if strings.TrimSpace(order.Email) == "" {
		return errors.New("Email is required")
	}
	if strings.TrimSpace(order.Location) == "" {
		return errors.New("Location is required")
	}
	if strings.TrimSpace(order.Size) == "" {
		return errors.New("Photo size is required")
	}
	if strings.TrimSpace(order.PaperType) == "" {
		return errors.New("Paper type is required")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	// Email validation using regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(order.Email) {
		return errors.New("Invalid email format")
	}

	return nil
}

// ValidateAndProcessImage validates and processes an image.
// It checks the file size, decodes the image, ensures the file
// extension is allowed. It validates the MIME type, and enforces
// maximum dimensions.
func (p *ImageProcessor) ValidateAndProcessImage(imgData []byte, opts models.ProcessingOptions) (image.Image, error) {

	// Validate file size
	fileSize := int64(len(imgData))
	if fileSize > models.MaxFileSize {
		return nil, fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", fileSize, models.MaxFileSize)
	}

	// Decode the image securely
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Validate file extension
	extension := fmt.Sprintf(".%s", format)
	if _, allowed := AllowedFileExtensions[extension]; !allowed {
		return nil, fmt.Errorf("invalid file type: %s, only JPG and PNG are allowed", extension)
	}

	// Validate MIME type using the first 512 bytes
	buffer := bytes.NewReader(imgData)
	header := make([]byte, 512)
	if _, err := buffer.Read(header); err != nil {
		return nil, fmt.Errorf("failed to read file data for MIME type validation: %w", err)
	}
	mimeType := http.DetectContentType(header)
	if !strings.HasPrefix(mimeType, "image/") {
		return nil, fmt.Errorf("invalid file type detected: %s", mimeType)
	}

	// Enforce maximum dimensions to prevent resource exhaustion
	maxWidth, maxHeight := opts.MaxDimensions.Width, opts.MaxDimensions.Height
	if maxWidth > 0 && maxHeight > 0 {
		if img.Bounds().Dx() > maxWidth || img.Bounds().Dy() > maxHeight {
			return nil, fmt.Errorf("image dimensions exceed maximum allowed size of %dx%d pixels", maxWidth, maxHeight)
		}
	}

	// Set the format if not already specified
	if opts.Format == "" {
		opts.Format = format
	}

	// Resize the image if the file size exceeds the target size
	if fileSize > models.TargetFileSize {
		opts.TargetSizeBytes = models.TargetFileSize
		return p.ProcessImageWithSizeTarget(img, opts)
	}

	// Accept the image without resizing if <= TargetFileSize
	return img, nil
}
