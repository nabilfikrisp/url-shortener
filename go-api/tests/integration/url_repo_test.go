package integration

import (
	"testing"

	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

func TestURLRepo(t *testing.T) {
	db := setupTestDB(t)
	repo := url.NewURLRepo(db)

	t.Run("Create URL", func(t *testing.T) {
		newURL := &url.URLModel{
			Original:   "https://example.com",
			ShortToken: "abc123",
		}

		err := repo.Create(newURL)
		assert.NoError(t, err)

		var found url.URLModel
		err = db.First(&found, "short_token = ?", "abc123").Error
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", found.Original)
	})
}
