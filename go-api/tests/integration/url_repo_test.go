package integration

import (
	"testing"

	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
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

func TestURLRepo(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

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

		t.Run("Duplicate Short Token", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

			firstURL := &url.URLModel{
				Original:   "https://first.com",
				ShortToken: "duplicate123",
			}
			err := repo.Create(firstURL)
			assert.NoError(t, err)

			secondURL := &url.URLModel{
				Original:   "https://second.com",
				ShortToken: "duplicate123",
			}
			err = repo.Create(secondURL)
			assert.Error(t, err)
			// Should be a database constraint error due to unique short token index
		})
	})

	t.Run("FindByShortToken", func(t *testing.T) {
		t.Run("Found", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

			// Create a URL first
			testURL := &url.URLModel{
				Original:   "https://google.com",
				ShortToken: "goog123",
				ClickCount: 5,
			}
			err := repo.Create(testURL)
			assert.NoError(t, err)

			// Find it
			found, err := repo.FindByShortToken("goog123")
			assert.NoError(t, err)
			assert.NotNil(t, found)
			assert.Equal(t, "https://google.com", found.Original)
			assert.Equal(t, "goog123", found.ShortToken)
			assert.Equal(t, 5, found.ClickCount)
		})

		t.Run("Not Found", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

			found, err := repo.FindByShortToken("nonexistent")
			assert.NoError(t, err)
			assert.Nil(t, found)
		})

		t.Run("Empty Token", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

			found, err := repo.FindByShortToken("")
			assert.Error(t, err)
			assert.Nil(t, found)
			assert.Equal(t, "short token is required", err.Error())
		})
	})

	t.Run("IncrementClickCount", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

			// Create a URL first
			testURL := &url.URLModel{
				Original:   "https://github.com",
				ShortToken: "gh123",
				ClickCount: 10,
			}
			err := repo.Create(testURL)
			assert.NoError(t, err)

			// Increment click count
			rowsAffected, err := repo.IncrementClickCount("gh123")
			assert.NoError(t, err)
			assert.Equal(t, int64(1), rowsAffected)

			// Verify the count was incremented
			found, err := repo.FindByShortToken("gh123")
			assert.NoError(t, err)
			assert.Equal(t, 11, found.ClickCount)
		})

		t.Run("Record Not Found", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

			rowsAffected, err := repo.IncrementClickCount("nonexistent")
			assert.Error(t, err)
			assert.Equal(t, int64(0), rowsAffected)
			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		})

		t.Run("Empty Token", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

			rowsAffected, err := repo.IncrementClickCount("")
			assert.Error(t, err)
			assert.Equal(t, int64(0), rowsAffected)
			assert.Equal(t, "short token is required", err.Error())
		})

		t.Run("Multiple Increments", func(t *testing.T) {
			db := setupTestDB(t)
			repo := url.NewURLRepo(db)

			// Create a URL first
			testURL := &url.URLModel{
				Original:   "https://stackoverflow.com",
				ShortToken: "so123",
				ClickCount: 0,
			}
			err := repo.Create(testURL)
			assert.NoError(t, err)

			// Increment multiple times
			for i := 0; i < 3; i++ {
				rowsAffected, err := repo.IncrementClickCount("so123")
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			}

			// Verify final count
			found, err := repo.FindByShortToken("so123")
			assert.NoError(t, err)
			assert.Equal(t, 3, found.ClickCount)
		})
	})

}
