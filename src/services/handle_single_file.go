package services

import (
	"io"
	"mime/multipart"

	"github.com/30Piraten/snapflow/models"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

// handleSingleFile validates and processes a single file
func handleSingleFile(c *fiber.Ctx, file *multipart.FileHeader, opts models.ProcessingOptions) error {

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
		return utils.HandleError(c, fiber.StatusBadRequest, "File validation failed", err)
	}

	order := new(models.PhotoOrder)

	// Process the file
	result := ProcessFile(c, file, opts, order)
	if result.Error != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to process file", &result.Error.Error)
	}

	return c.JSON(fiber.Map{
		"message":  "File processed successfully",
		"filePath": result.Path,
	})
}
