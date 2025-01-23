package main

import (
	"log"
	"os"
	"time"

	"github.com/30Piraten/snapflow/routes"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

// Main configures the route for the Photo Upload service
func main() {
	// Create a new instance of the template engine
	engine := html.New("./views", ".html")

	// Enable template engine reloading in development
	engine.Reload(true) // Enable this during development
	// Create a new Fiber app with the template engine
	app := fiber.New(fiber.Config{
		Views: engine,

		// Fiber uses bodyLimit to enforce request size limit
		// which for some reason might be too low. Thus if the
		// uploaded file exceeds this limit, the request is rejected
		// before the application logic runs. Hence the direct use here
		BodyLimit: 50 * 1024 * 1024, // 50MB
	})

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load .env file")
	}

	// Get PORT from .env
	PORT := os.Getenv("PORT")

	// Initialize logger
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("Failed to initalize logger: %v", err)
	}
	defer utils.Logger.Sync()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:1234",
		AllowMethods: "POST,OPTIONS",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
	}))

	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Handling request for %s", c.Path())
		return c.Next()
	})

	// Register route
	routes.Handler(app)

	// Setup static file serving if needed
	app.Static("/public", "./public")

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Listen(":" + PORT))
}
