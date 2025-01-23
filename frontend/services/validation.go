package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

// AllowedFileExtensions defines permitted image file extensions
var AllowedFileExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
}

// ValidateOrder checks if all required fields are present and valid
func ValidateOrder(order *utils.PhotoOrder) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	var missingFields []string

	if strings.TrimSpace(order.FullName) == "" {
		missingFields = append(missingFields, "full name")
	}
	if strings.TrimSpace(order.Email) == "" {
		missingFields = append(missingFields, "email")
	}
	if strings.TrimSpace(order.Location) == "" {
		missingFields = append(missingFields, "location")
	}
	if strings.TrimSpace(order.Size) == "" {
		missingFields = append(missingFields, "size")
	}
	if strings.TrimSpace(order.PaperType) == "" {
		missingFields = append(missingFields, "paper type")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	// Basic email validation
	if !strings.Contains(order.Email, "@") || !strings.Contains(order.Email, ".") {
		return errors.New("invalid email format")
	}

	return nil
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

// ValidateUploadedFile validates file name, extension and MIME type
func ValidateUploadedFile(file *multipart.FileHeader) error {
	// Extract file extension
	extension := strings.ToLower(filepath.Ext(file.Filename))

	// Validate file extension
	if _, allowed := AllowedFileExtensions[extension]; !allowed {
		return fmt.Errorf("invalid file type: %s, only JPG and PNG are allowed", extension)
	}

	// Open the file to check its MIME type
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for validation: %w", err)
	}
	defer src.Close()

	// Buffer the first 512 bytes for MIME type detection
	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return fmt.Errorf("failed to read file for MIME type validation: %w", err)
	}

	// Check MIME type
	mimeType := http.DetectContentType(buffer)
	if !strings.HasPrefix(mimeType, "image/") {
		return fmt.Errorf("invalid file type detected: %s", mimeType)
	}

	return nil
}

func (p *ImageProcessor) ValidateAndProcessImage(imgData []byte, opts ProcessingOptions) (image.Image, error) {

	// We must validate the file size first
	fileSize := int64(len(imgData))

	// Reject single file > 100MB
	if fileSize > MaxFileSize {
		return nil, fmt.Errorf("file szie %d bytes exceeds maximum allowed szie of %d bytes", fileSize, MaxFileSize)
	}

	// Next we securely decode the image
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// And set the format of the image if it's not specified
	if opts.Format == "" {
		opts.Format = format
	}

	// Next, if the file is between 1MB and 100MB, resize
	if fileSize > TargetFileSize {
		opts.TargetSizeBytes = TargetFileSize
		return p.ProcessImageWithSizeTarget(img, opts)
	}

	// Accept file <=1MB without resizing
	return img, nil
}
