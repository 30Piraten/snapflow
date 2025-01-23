package routes

import "github.com/gofiber/fiber/v2"

// ConfirmationHandler renders a confirmation page for an order submission.
// It retrieves form data, counts uploaded photos and renders the confirmation
// template with the retrieved data.
func ConfirmationHandler(c *fiber.Ctx) error {
	// Retrieve form data (this is a simplified example)
	fullName := c.FormValue("fullName")
	email := c.FormValue("email")
	location := c.FormValue("location")
	size := c.FormValue("size")
	paperType := c.FormValue("paperType")

	// Count uploaded photos
	form, err := c.MultipartForm()
	if err != nil {
		return c.Render("error", fiber.Map{"message": "Error processing order"})
	}

	photoCount := len(form.File["photos"])

	// Render confirmation template
	return c.Render("confirmation", fiber.Map{
		"FullName":   fullName,
		"Email":      email,
		"Location":   location,
		"Size":       size,
		"PaperType":  paperType,
		"PhotoCount": photoCount,
	})
}
