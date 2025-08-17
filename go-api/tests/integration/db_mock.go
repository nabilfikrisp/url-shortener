package integration

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	_ = godotenv.Load("../../.env.test")

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL is not set in .env.test")
	}

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
