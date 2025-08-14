package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nabilfikrisp/url-shortener/internal/config"
	"github.com/nabilfikrisp/url-shortener/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", db.Name())
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	app.Listen(":" + cfg.Port)
}
