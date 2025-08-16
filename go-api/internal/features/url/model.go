package url

import (
	"time"

	"gorm.io/gorm"
)

// URL represents the mapping between the original long URL and its short token.
type URLModel struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ShortToken string         `gorm:"uniqueIndex;size:20;not null" json:"short_token"`
	Original   string         `gorm:"not null" json:"original"`
	ClickCount int            `gorm:"default:0" json:"click_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (URLModel) TableName() string {
	return "urls"
}
