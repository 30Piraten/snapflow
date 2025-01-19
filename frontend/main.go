package main

import (
	"log"
	"os"

	"github.com/30Piraten/snapflow/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Create a new instance of the template engine
	engine := html.New("./views", ".html")

	// Enable template engine reloading in development
	engine.Reload(true) // Enable this during development

	// Create a new Fiber app with the template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Setup static file serving if needed
	app.Static("/public", "./public")

	// Setup routes
	handlers.SetupRoutes(app)

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Listen(":9090"))
}
