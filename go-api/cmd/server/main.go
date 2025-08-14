package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nabilfikrisp/url-shortener/internal/config"
	"github.com/nabilfikrisp/url-shortener/internal/database"
	"github.com/nabilfikrisp/url-shortener/internal/features/url"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", db.Name())
	}

	// Auto-migrate model
	if err := db.AutoMigrate(&url.URLModel{}); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	urlHandler := url.InitURLHandler(db)
	url.RegisterRoutes(app, urlHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	log.Fatal(app.Listen(":" + cfg.Port))
}
