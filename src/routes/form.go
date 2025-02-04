package routes

import "github.com/gofiber/fiber/v2"

// ServeUploadForm handles rendering the upload form
func ServeUploadForm(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Photo Upload Service",
	})
}
