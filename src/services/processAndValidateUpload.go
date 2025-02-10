package services

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/30Piraten/snapflow/models"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

type NewProcessError struct {
	*models.ProcessingError
}

// Error returns a string representation of the
// error in the format "type: code - message".
func (e *NewProcessError) Error() string {
	return fmt.Sprintf("%s: %s - %s", e.Type, e.Code, e.Message)
}

// ProcessUploadedFiles parses the uploaded files and processes
// them accordingly. If there is a single or there are multiple
// files, a JSON response is returned.
func ProcessUploadedFiles(c *fiber.Ctx) error {

	// Parse the uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadGateway, "Failed to parse multipart form", err)
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return utils.HandleError(c, fiber.StatusBadRequest, "No files uploaded", nil)
	}

	// Validate total file count
	if len(files) > models.MaxFileCount {
		return utils.HandleError(c, fiber.StatusBadRequest, fmt.Sprintf("Too many files uploaded. Maximum allowed is %d", models.MaxFileCount), nil)
	}

	// Process single or multiple files
	opts := models.ProcessingOptions{
		Quality:         models.HighQuality,
		TargetSizeBytes: models.TargetFileSize,
		Format:          "jpeg",
		MaxDimensions: models.Dimensions{
			Width:  6000,
			Height: 6000,
		},
	}

	// Handle single file
	if len(files) == 1 {
		if err := handleSingleFile(c, files[0], opts); err != nil {
			return fmt.Errorf("failed to process file %s: %w", files[0].Filename, err)
		}
		return nil
	}

	// Validate and handle multiple files
	_, errors := ProcessMultipleFiles(c, files, opts)
	if len(errors) > 0 {
		// Collect all errors into one
		var errMsg []string
		for _, e := range errors {
			errMsg = append(errMsg, e.Error())
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Some files failed to process", errors[0])
	}

	return nil
}

// generateUniqueFileName generates a unique filename
func generateUniqueFileName(originalName string) string {
	extension := filepath.Ext(originalName)
	basename := strings.TrimSuffix(originalName, extension)

	// Sanitize the basename
	basename = utils.Sanitize(basename)

	randomNumber := rand.Intn(1000)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d_%d%s", basename, timestamp, randomNumber, extension)
}
