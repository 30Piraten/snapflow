package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func isValidMimeType(mimeType string) bool {
	return mimeType == "image/jpeg" || mimeType == "image/png"
}

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
