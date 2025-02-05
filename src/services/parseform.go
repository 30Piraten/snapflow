package services

import (
	"github.com/30Piraten/snapflow/models"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// parseOrderDetails parses the order details from the request body and validates
// the required fields. If the parsing or validation fails, it returns an error.
// If the parsing and validation succeed, it returns the parsed order details.
func ParseOrderDetails(c *fiber.Ctx) (*models.PhotoOrder, error) {
	order := new(models.PhotoOrder)

	if err := c.BodyParser(order); err != nil {
		utils.Logger.Error("Form parsing failed", zap.Error(err))
		return nil, err
	}

	// Validate required fields
	if err := ValidateOrder(order); err != nil {
		return nil, err
	}

	utils.Logger.Info("Order details parsed successfully",
		zap.String("fullName", order.FullName),
		zap.String("email", order.Email))

	return order, nil
}
