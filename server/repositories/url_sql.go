package repositories

import (
	"url_shortener/models"

	"gorm.io/gorm"
)

type SQLURLRepository struct {
	db *gorm.DB
}

// NewSQLURLRepository returns a new instance of a GORM-based URLRepository.
func NewSQLURLRepository(db *gorm.DB) URLRepository {
	return &SQLURLRepository{
		db: db,
	}
}

func (r *SQLURLRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r *SQLURLRepository) FindByShortURL(shortURL string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("short_url = ?", shortURL).First(&url).Error
	return &url, err
}

func (r *SQLURLRepository) FindByLongURL(longURL string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("long_url = ?", longURL).First(&url).Error
	return &url, err
}
