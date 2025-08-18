package integration

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nabilfikrisp/url-shortener/internal/features/url"
)

func setupTestApp(t *testing.T) *fiber.App {
	db := SetupTestDB(t)

	// init handler + register routes
	handler := url.InitURLHandler(db)
	app := fiber.New()
	url.RegisterRoutes(app, handler)

	return app
}
