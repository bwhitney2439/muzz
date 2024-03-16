package main

import (
	"flag"
	"log"

	"github.com/bwhitney2439/muzz/database"
	"github.com/bwhitney2439/muzz/handlers"
	"github.com/bwhitney2439/muzz/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Connected with database
	database.Connect()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Create a /api/v1 endpoint
	v1 := app.Group("/api/v1")

	// Auth routes
	v1.Post("/user/create", handlers.CreateUser)
	v1.Post("/login", handlers.LoginUser)

	// Protected routes
	v1.Use(middleware.Protected())
	v1.Get("/discover", handlers.Discover)
	v1.Post("/swipe", handlers.Swipe)

	// Listen on port 3000
	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000
}
