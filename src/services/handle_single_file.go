package services

import (
	"io"
	"mime/multipart"

	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

// handleSingleFile validates and processes a single file. It first validates the
// file using the provided ProcessingOptions, and if the validation fails, it
// returns a 400 Bad Request error. If the validation succeeds, it processes the
// file and returns a JSON response with the path of the processed file if
// successful. If there are any errors during processing, it returns a 500 Internal
// Server Error with the error message.
func handleSingleFile(c *fiber.Ctx, file *multipart.FileHeader, opts ProcessingOptions) error {

	// Open the file
	source, err := file.Open()
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer source.Close()

	// Read the file data into a byte slice
	fileData, err := io.ReadAll(source)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to read file data", err)
	}

	// Validate the file before processing
	processor := NewImageProcessor(utils.Logger)
	if _, err = processor.ValidateAndProcessImage(fileData, opts); err != nil {
		// utils.Logger.Error("File validation failed", zap.Error(err))
		return utils.HandleError(c, fiber.StatusBadRequest, "File validation failed", err)
	}

	// Process the file
	result := ProcessFile(file, opts, order)
	if result.Error != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to process file", result.Error)
	}

	return c.JSON(fiber.Map{
		"message":  "File processed successfully",
		"filePath": result.Path,
	})
}
