package routes

import (
	"os"
	"strings"

	"github.com/30Piraten/snapflow/models"
	"github.com/30Piraten/snapflow/services"
	svc "github.com/30Piraten/snapflow/services"
	"github.com/30Piraten/snapflow/url"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

// HandleOrderSubmission is the main entry point for the order
// submission process. It will parse order details, generate
// a presigned URL, process uploaded files and return a successful
// response containing the order details, presigned URL and order ID.
func HandleOrderSubmission(c *fiber.Ctx) error {

	// trustedOrigin defines the trusted frontend domain
	trustedOrigin := os.Getenv("TRUSTED_ORIGIN")

	// Get the Referer and Origin headers
	referer := c.Get("Referer")
	origin := c.Get("Origin")

	// Check if the Referer or origin header matches the trusted domain
	if !strings.HasPrefix(referer, trustedOrigin) && origin != trustedOrigin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden: Invalid origin or referer",
		})
	}

	// Limit the upload size to ensure large files are rejected early
	c.Locals("limit", "50MB") // -> Review

	utils.Logger.Info("Starting order submission processing")

	// Parse the order details
	order, err := services.ParseOrderDetails(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Failed to parse order details", err)
	}

	// Generate presigned URL
	presignedResponse, err := url.GeneratePresignedURL(order)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "File validation or processing failed", err)
	}

	// Process uploaded photos
	if err := svc.ProcessUploadedFiles(c); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Failed to process files", err)
	}

	// Return a successful response
	return c.JSON(models.ResponseData{
		Message:      "Order received successfully",
		Order:        order,
		PresignedURL: presignedResponse.URLs,
		OrderID:      presignedResponse.OrderID,
	})
}
