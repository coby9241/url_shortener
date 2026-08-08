package repositories

import (
	"url_shortener/models"

	"gorm.io/gorm"
)

// URLRepository defines the interface for URL database operations.
type URLRepository interface {
	Create(url *models.URL) error
	FindByShortURL(shortURL string) (*models.URL, error)
	FindByLongURL(longURL string) (*models.URL, error)
}

type sqlURLRepository struct {
	db *gorm.DB
}

// NewSQLURLRepository returns a new instance of a GORM-based URLRepository.
func NewSQLURLRepository(db *gorm.DB) URLRepository {
	return &sqlURLRepository{
		db: db,
	}
}

func (r *sqlURLRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r *sqlURLRepository) FindByShortURL(shortURL string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("short_url = ?", shortURL).First(&url).Error
	return &url, err
}

func (r *sqlURLRepository) FindByLongURL(longURL string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("long_url = ?", longURL).First(&url).Error
	return &url, err
}
