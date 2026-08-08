package models

type URL struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ShortURL string `gorm:"uniqueIndex;not null" json:"short_url"`
	LongURL  string `gorm:"not null" json:"long_url"`
}
