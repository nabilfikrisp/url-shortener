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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// reset schema before each test
	err = db.Migrator().DropTable(&url.URLModel{})
	if err != nil {
		t.Fatalf("failed to drop table: %v", err)
	}
	err = db.AutoMigrate(&url.URLModel{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// cleanup after test
	t.Cleanup(func() {
		db.Migrator().DropTable(&url.URLModel{})
	})

	return db
}
