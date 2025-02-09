package services

import (
	"fmt"

	"github.com/30Piraten/snapflow/models"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// parseOrderDetails parses the order details from
// the request body and validates the required fields.
func ParseOrderDetails(c *fiber.Ctx) (*models.PhotoOrder, error) {
	order := new(models.PhotoOrder)

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		utils.Logger.Error("Multipart form parsing failed", zap.Error(err))
		return nil, err
	}

	// Get form values
	order.FullName = c.FormValue("fullName")
	order.Location = c.FormValue("location")
	order.Size = c.FormValue("size")
	order.PaperType = c.FormValue("paperType")
	order.Email = c.FormValue("email")

	// Get files from form
	files := form.File["photos"]
	if len(files) == 0 {
		utils.Logger.Error("No photos found in the form data")
		return nil, fmt.Errorf("no photos provided in the order")
	}

	// Assign files to order
	order.Photos = files

	// Log the number of files received
	utils.Logger.Info("Files received",
		zap.Int("fileCount", len(files)),
		zap.String("fullName", order.FullName),
		zap.String("email", order.Email))

	// Validate the required fields
	if err := ValidateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}
