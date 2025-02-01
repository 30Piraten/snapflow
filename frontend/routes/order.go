package routes

import (
	"os"
	"strings"

	svc "github.com/30Piraten/snapflow/services"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// HandleOrderSubmission is the main entry point for the order submission process.
// It will parse order details, generate a presigned URL, process uploaded files
// and return a successful response containing the order details, presigned URL
// and order ID.
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
	c.Locals("limit", "50MB")

	utils.Logger.Info("Starting order submission processing")

	// Parse the order details
	order, err := parseOrderDetails(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Failed to parse order details", err)
	}

	// Generate presigned URL
	presignedResponse, err := utils.GeneratePresignedURL(order)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "File validation or processing failed", err)
	}

	// Log the generated presigned URL for debugging
	// utils.Logger.Info("Generated Presigned URL:", zap.String("URL", presignedResponse.URL))x

	// Process uploaded photos
	if err := svc.ProcessUploadedFiles(c); err != nil {
		// // Return an error response if the file processing failes
		// utils.Logger.Error("File validation/processing failed", zap.Error(err))
		return utils.HandleError(c, fiber.StatusBadRequest, "Failed to process files", err)
	}

	// Return a successful response
	return c.JSON(svc.ResponseData{
		Message:      "Order received successfully",
		Order:        order,
		PresignedURL: presignedResponse.URL,
		OrderID:      presignedResponse.OrderID,
	})
}

// parseOrderDetails parses the order details from the request body and validates
// the required fields. If the parsing or validation fails, it returns an error.
// If the parsing and validation succeed, it returns the parsed order details.
func parseOrderDetails(c *fiber.Ctx) (*utils.PhotoOrder, error) {
	order := new(utils.PhotoOrder)
	if err := c.BodyParser(order); err != nil {
		utils.Logger.Error("Form parsing failed", zap.Error(err))
		return nil, err
	}

	// Validate required fields
	if err := svc.ValidateOrder(c, order); err != nil {
		return nil, err
	}

	utils.Logger.Info("Order details parsed successfully",
		zap.String("fullName", order.FullName),
		zap.String("email", order.Email))

	return order, nil
}
