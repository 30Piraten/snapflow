package routes

import (
	h "github.com/30Piraten/snapflow/handlers"
	"github.com/gofiber/fiber/v2"
)

// Handler configured the routes for our photo upload service
func Handler(app *fiber.App) {
	// Save the upload form
	app.Get("/", ServeUploadForm)

	// Handle form submissions
	app.Post("/submit-order", HandleOrderSubmission)

	// Register the presigned URL route
	h.Upload(app)
}
