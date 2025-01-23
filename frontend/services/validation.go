package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"path/filepath"
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

// ValidateUpload checks the upload constraints for the given files.
// It validates the total request size against MaxTotalUploadSize,
// the number of files against MaxFileCount, and each file's size
// against MaxFileSize. If any of these validations fail, it returns
// an appropriate error.
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

// ValidateUploadedFile checks the validity of a given file based on its
// extension and MIME type. It ensures the file has a permitted extension
// (JPG or PNG) and that its MIME type is an image type. If any validation
// fails, it returns an error detailing the issue.
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
	// ValidateAndProcessImage validates the size and decodes the given image data, and if the image is above the target size, it
	// resizes the image to meet the target size. If the image is below the target size, it is returned as is. The function
	// returns the processed image and any error that occurred during processing.
	// We must validate the file size first
	fileSize := int64(len(imgData))

	// Reject single file > 50MB
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

	// Next, if the file is between 1MB and 50MB, resize
	if fileSize > TargetFileSize {
		opts.TargetSizeBytes = TargetFileSize
		return p.ProcessImageWithSizeTarget(img, opts)
	}

	// Accept file <=1MB without resizing
	return img, nil
}
