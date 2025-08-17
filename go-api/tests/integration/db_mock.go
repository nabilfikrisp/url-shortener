package integration

import (
	"testing"

	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	// reset schema before each test
	db.Migrator().DropTable(&url.URLModel{})
	db.AutoMigrate(&url.URLModel{})

	// cleanup after test
	t.Cleanup(func() {
		db.Migrator().DropTable(&url.URLModel{})
	})

	return db
}
