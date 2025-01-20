package handlers

import (
	"fmt"

	"github.com/30Piraten/snapflow/routes"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

// ResponseData to structure response
type ResponseData struct {
	Message      string            `json:"message"`
	Order        *utils.PhotoOrder `json:"order"`
	PresignedURL string            `json:"presigned_url"`
	OrderID      string            `json:"order_id"`
}

func upload(app *fiber.App) {
	app.Post("/generate-upload-url", routes.HandleGenerateUploadURL)
}

// Handler configures the routes for our photo upload service
func Handler(app *fiber.App) {

	// Render the upload form
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Photo Upload Service",
		})
	})

	// Handle form submission
	app.Post("/submit-order", func(c *fiber.Ctx) error {
		// Parse the multipart form
		order := new(utils.PhotoOrder)
		if err := c.BodyParser(order); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse form",
			})
		}

		// Make internal call to generate the presigned URL
		// app.Post("/generate-upload-url", r.HandleGenerateUploadURL)

		// // create new context with the form data
		// ctx := c.Context()
		// ctx.SetBody(c.Body())

		// // call the upload URL generator
		// if err := r.HandleGenerateUploadURL(c); err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 		"error": "Failed to generate upload URL",
		// 	})
		// }

		// Use the shared presigned URL generation function
		presignedResponse, err := utils.GeneratePresignedURL(order)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Validatde the form fields
		// if order.FullName == "" || len(order.Photos) == 0 {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"error": "Full name and photos are required!",
		// 	})
		// }

		// Handle file uploads
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to process uploaded files",
			})
		}

		// Get the files from form
		files := form.File["photos"]

		// Process each uploaded file
		for _, file := range files {

			// TODO!
			// Save file to disk or cloud storage
			// logic to store photos in S3 bucket [original & processed]!
			err := c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to save file",
				})
			}
		}

		return c.JSON(fiber.Map{
			"message":       "Order received successfully",
			"order":         order,
			"presigned_url": presignedResponse.URL,
			"order_id":      presignedResponse.OrderID,
		})
	})

	// Register the presigned URL route
	upload(app)
}
