package url

import (
	"gorm.io/gorm"
)

type URLRepo interface {
	Create(url *URLModel) error
}

type urlRepo struct {
	db *gorm.DB
}

func NewURLRepo(db *gorm.DB) URLRepo {
	return &urlRepo{
		db: db,
	}
}

func (r *urlRepo) Create(url *URLModel) error {
	return r.db.Create(url).Error
}
