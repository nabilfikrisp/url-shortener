package url

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitURLHandler(db *gorm.DB) URLHandler {
	repo := NewURLRepo(db)
	service := NewURLService(repo)
	handler := NewURLHandler(service)
	return handler
}

func RegisterRoutes(app *fiber.App, handler URLHandler) {
	app.Post("/shorten", handler.Create)
}
