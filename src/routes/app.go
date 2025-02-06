package routes

import (
	h "github.com/30Piraten/snapflow/handlers"
	"github.com/gofiber/fiber/v2"
)

// Handler configures the routes for the application.
// It sets up the endpoint for rendering the upload form,
// handling order submissions, and registering the route
// for generating presigned URLs.
func Handler(app *fiber.App) {
	// Save the upload form
	app.Get("/", ServeUploadForm)

	// Handle form submissions
	app.Post("/submit-order", HandleOrderSubmission)

	// Register the presigned URL route
	h.Upload(app)
}
