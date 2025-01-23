package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// isValidMimeType returns true if the given MIME type is either "image/jpeg" or "image/png", both of which are valid image formats.
func isValidMimeType(mimeType string) bool {
	return mimeType == "image/jpeg" || mimeType == "image/png"
}

// HandleError logs an error message with the given status and error details,
// then sends a JSON response with the error message and details.
// It initializes the logger if it hasn't been initialized yet.
// If the provided error is nil, it defaults to an "unspecified error".
func HandleError(c *fiber.Ctx, status int, message string, err error) error {
	if Logger == nil {
		InitLogger()
	}

	if err == nil {
		err = fmt.Errorf("unspecified error")
	}

	Logger.Error(message,
		zap.Error(err),
		zap.Int("status", status))

	return c.Status(status).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}
