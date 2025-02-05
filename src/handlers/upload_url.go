package handlers

import (
	"github.com/30Piraten/snapflow/models"
	"github.com/30Piraten/snapflow/url"
	"github.com/gofiber/fiber/v2"
)

// Upload configures the route for the presigned URL generation.
// It adds the /generate-upload-url endpoint to the provided fiber app.
// This endpoint expects a POST request with a JSON body containing the
// required fields for generating a presigned URL.
func Upload(app *fiber.App) {
	app.Post("/generate-upload-url", HandleGenerateUploadURL)
}

// HandleGenerateUploadURL handles the request to generate a presigned URL for uploading a photo
// to S3. It expects a JSON body containing the required fields for generating a presigned URL.
// If the parsing of the form data fails, it responds with a 400 status code and an error message.
// If the generation of the presigned URL fails, it responds with a 500 status code and an error message.
// Otherwise, it responds with a JSON object containing the presigned URL and the order ID.
func HandleGenerateUploadURL(c *fiber.Ctx) error {

	// Parse the form data from the webpage into the defined struct
	// PhotoOrder serves as a base for the new "order" struct
	order := new(models.PhotoOrder)
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	response, err := url.GeneratePresignedURL(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
