package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"net/http"
	"regexp"
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

// ValidateOrder validates a PhotoOrder instance, returning an error if required fields are missing.
func ValidateOrder(c *fiber.Ctx, order *utils.PhotoOrder) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	var missingFields []string

	if strings.TrimSpace(order.FullName) == "" {
		// missingFields = append(missingFields, "full name")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Full name is required",
			"errorField": "FullName",
		})
	}

	if strings.TrimSpace(order.Email) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Email is required",
			"errorField": "Email",
		})
	}
	if strings.TrimSpace(order.Location) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Print location is required",
			"errorField": "Location",
		})
	}
	if strings.TrimSpace(order.Size) == "" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":      "Print szie is required",
			"errorField": "Size",
		})
	}
	if strings.TrimSpace(order.PaperType) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Paper type is required",
			"errorField": "PaperType",
		})
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	// Email validation using regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(order.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid email format",
			"errorField": "Email",
		})
	}

	return nil
}

// ValidateAndProcessImage validates and processes an image. It checks the file size, decodes the image, ensures the file
// extension is allowed, validates the MIME type, and enforces maximum dimensions. If the file size exceeds the target
// size, it resizes the image to meet the target size. The processed image is returned, or an error is returned if any
// validation or processing fails.
func (p *ImageProcessor) ValidateAndProcessImage(imgData []byte, opts ProcessingOptions) (image.Image, error) {
	// Validate file size
	fileSize := int64(len(imgData))
	if fileSize > MaxFileSize {
		return nil, fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", fileSize, MaxFileSize)
	}

	// Decode the image securely
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Check if the file extension is allowed
	var extension string
	switch format {
	case "jpeg":
		extension = ".jpg" // Treat "jpeg" as ".jpg"
	case "png":
		extension = ".png"
	default:
		extension = "." + format // Unrecognized format
	}

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
	if fileSize > TargetFileSize {
		opts.TargetSizeBytes = TargetFileSize
		return p.ProcessImageWithSizeTarget(img, opts)
	}

	// Accept the image without resizing if <= TargetFileSize
	return img, nil
}
