package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"url_shortener/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

// InitDB initializes GORM database using PostgreSQL and performs auto-migration.
func InitDB() (*gorm.DB, error) {
	// Load .env file (for local development)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Check for DATABASE_URL first (used by Render, Heroku, etc.)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		// Render's PostgreSQL may provide "postgres://" but gorm expects "postgresql://"
		if strings.HasPrefix(databaseURL, "postgres://") {
			databaseURL = strings.Replace(databaseURL, "postgres://", "postgresql://", 1)
		}
		database, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		// Auto-migrate URL model
		if err := database.AutoMigrate(&models.URL{}); err != nil {
			return nil, fmt.Errorf("failed to auto-migrate URL schema: %w", err)
		}
		return database, nil
	}

	// Fallback to individual parameters (existing logic)
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "url_shortener"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, password, dbname, port)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate URL model
	if err := database.AutoMigrate(&models.URL{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate URL schema: %w", err)
	}

	return database, nil
}