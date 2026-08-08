package repositories

import (
	"errors"
	"strings"
	"testing"
	"url_shortener/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSQLURLRepository_Create(t *testing.T) {
	// Initialize in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&models.URL{}); err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	repo := NewSQLURLRepository(db)

	t.Run("successful insertion", func(t *testing.T) {
		url := &models.URL{
			ShortURL: "abc1234",
			LongURL:  "https://example.com",
		}

		err := repo.Create(url)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify record is created
		var found models.URL
		if err := db.First(&found, "short_url = ?", "abc1234").Error; err != nil {
			t.Errorf("failed to find created URL in database: %v", err)
		}
		if found.LongURL != url.LongURL {
			t.Errorf("expected long URL %s, got %s", url.LongURL, found.LongURL)
		}
	})

	t.Run("duplicate key violation error", func(t *testing.T) {
		url1 := &models.URL{
			ShortURL: "duplicate",
			LongURL:  "https://first.com",
		}
		url2 := &models.URL{
			ShortURL: "duplicate",
			LongURL:  "https://second.com",
		}

		// First insert should succeed
		if err := repo.Create(url1); err != nil {
			t.Fatalf("failed to insert first URL: %v", err)
		}

		// Second insert with duplicate short_url should fail
		err := repo.Create(url2)
		if err == nil {
			t.Error("expected duplicate key violation error, got nil")
		}

		// Ensure it is a unique key violation (checking via gorm.ErrDuplicatedKey or string containing UNIQUE)
		if !errors.Is(err, gorm.ErrDuplicatedKey) && !strings.Contains(err.Error(), "UNIQUE") {
			t.Errorf("expected duplicate key error or UNIQUE constraint error, got %v", err)
		}
	})
}
