package routes

import (
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
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate presigned URL", err)
	}

	// Process uploaded photos -> TODO
	err = svc.ProcessUploadedFiles(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to process files", err)
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
