package handlers

import (
	"fmt"

	"github.com/30Piraten/snapflow/routes"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
)

// // PhotoOrder represents the structure of our form data
// type PhotoOrder struct {
// 	FullName  string            `form:"fullName"`
// 	Location  string            `form:"location"`
// 	Size      string            `form:"size"`
// 	PaperType string            `form:"paperType"`
// 	Email     string            `form:"email"`
// 	Photos    []*fiber.FormFile `form:"photos"`
// }

// SetupRoutes configures the routes for our photo upload service
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

		// Validatde the form fields
		if order.FullName == "" || len(order.Photos) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Full name and photos are required!",
			})
		}

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

		// Return success response
		// return c.JSON(fiber.Map{
		// 	"message": "Order received successfully",
		// 	"order":   order,
		// })

		return c.Redirect("/generate-upload-url")
	})

	// Register the presigned URL route
	routes.Upload(app)
}
