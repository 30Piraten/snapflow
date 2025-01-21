package handlers

import (
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

func Upload(app *fiber.App) {
	app.Post("/generate-upload-url", HandleGenerateUploadURL)
}

func HandleGenerateUploadURL(c *fiber.Ctx) error {

	// Parse the form data from the webpage into the defined struct
	// PhotoOrder serves as a base for the new "order" struct
	order := new(utils.PhotoOrder)
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	response, err := utils.GeneratePresignedURL(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
