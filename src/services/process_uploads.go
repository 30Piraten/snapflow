package services

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var order *PhotoOrder

// Error returns a string representation of the error in the format "type: code - message".
func (e *ProcessingError) Error() string {
	return fmt.Sprintf("%s: %s - %s", e.Type, e.Code, e.Message)
}

// ProcessUploadedFiles parses the uploaded files and processes them accordingly.
// If there is a single file, it is processed and a JSON response is returned
// containing the file path. If there are multiple files, they are processed
// concurrently and a JSON response is returned containing the paths of the processed
// files. If there were any errors while processing the files, they are collected and returned
// in the response as well. The function returns an error if there were any issues while processing the files.
func ProcessUploadedFiles(c *fiber.Ctx) error {

	// Parse the uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		// return utils.HandleError(c, fiber.StatusBadGateway, "Failed to parse multipart form", err)
		return fmt.Errorf("failed to parse multipart form: %w", err)
	}

	files := form.File["photos"]
	if len(files) == 0 {
		// return utils.HandleError(c, fiber.StatusBadRequest, "No files uploaded", nil)
		return fmt.Errorf("no files uploaded")
	}

	// Validate total file count
	if len(files) > MaxFileCount {
		// return utils.HandleError(c, fiber.StatusBadRequest, fmt.Sprintf("Too many files uploaded. Maximum allowed is %d", MaxFileCount), nil)
		return fmt.Errorf("Too many files uploaded. Maximum allowed is %d", MaxFileCount)
	}

	// Process single or multiple files
	opts := ProcessingOptions{
		Quality:         HighQuality,
		TargetSizeBytes: TargetFileSize,
		Format:          "jpeg",
		MaxDimensions: Dimensions{
			// Need to review: TODO
			Width:  5000,
			Height: 5000,
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
	_, errors := ProcessMultipleFiles(files, opts) // results should be added here
	// If there are errors, return the first one
	if len(errors) > 0 {
		// Collect all errors into one
		var errMsg []string
		for _, e := range errors {
			errMsg = append(errMsg, e.Error())
		}
		// return utils.HandleError(c, fiber.StatusInternalServerError, "Some files failed to process", errors[0])
		return fmt.Errorf("one or more files failed processing: %s", strings.Join(errMsg, "; "))
	}

	// // Prepare the response for successfully processed files
	// var filePaths []string
	// for _, result := range results {
	// 	filePaths = append(filePaths, result.Path)
	// }

	// return c.JSON(fiber.Map{
	// 	"message":   "Files processed succesfully",
	// 	"filePaths": filePaths,
	// })

	return nil
}

// generateUniqueFileName generates a unique filename by appending a timestamp
// to the base name of the original file name.
func generateUniqueFileName(originalName string) string {
	extension := filepath.Ext(originalName)
	basename := strings.TrimSuffix(originalName, extension)

	// Sanitize the basename
	basename = Sanitize(basename)

	randomNumber := rand.Intn(1000)
	timestamp := time.Now().UnixNano() // -> Make random
	return fmt.Sprintf("%s_%d_%d%s", basename, timestamp, randomNumber, extension)
}
