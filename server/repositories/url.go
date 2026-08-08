package repositories

import (
	"url_shortener/models"
)

// URLRepository defines the interface for URL database operations.
type URLRepository interface {
	Create(url *models.URL) error
	FindByShortURL(shortURL string) (*models.URL, error)
	FindByLongURL(longURL string) (*models.URL, error)
}
