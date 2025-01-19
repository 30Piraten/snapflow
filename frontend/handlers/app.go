package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// PhotoOrder represents the structure of our form data
type PhotoOrder struct {
	FullName  string            `form:"fullName"`
	Location  string            `form:"location"`
	Size      string            `form:"size"`
	PaperType string            `form:"paperType"`
	Photos    []*fiber.FormFile `form:"photos"`
}

// SetupRoutes configures the routes for our photo upload service
func Handler(app *fiber.App) {
	// // Initialize template engine
	// engine := html.New("./views", ".html")

	// app = fiber.New(fiber.Config{
	// 	Views: engine,
	// })

	// Render the upload form
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Photo Upload Service",
		})
	})

	// Handle form submission
	app.Post("/upload", func(c *fiber.Ctx) error {
		// Parse the multipart form
		order := new(PhotoOrder)
		if err := c.BodyParser(order); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse form",
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
			// Save file to disk or cloud storage
			// logic to store photos in S3 bucket!
			err := c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to save file",
				})
			}
		}

		// Return success response
		return c.JSON(fiber.Map{
			"message": "Order received successfully",
			"order":   order,
		})
	})
}
