package url

import (
	"errors"

	"gorm.io/gorm"
)

type URLRepo interface {
	Create(url *URLModel) error
	FindByShortToken(shortToken string) (*URLModel, error)
	IncrementClickCount(shortToken string) (int64, error)
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

func (r *urlRepo) FindByShortToken(shortToken string) (*URLModel, error) {
	if shortToken == "" {
		return nil, errors.New("short token is required")
	}

	var url URLModel
	if err := r.db.Where("short_token = ?", shortToken).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &url, nil
}

func (r *urlRepo) IncrementClickCount(shortToken string) (int64, error) {
	if shortToken == "" {
		return 0, errors.New("short token is required")
	}

	result := r.db.Model(&URLModel{}).Where("short_token = ?", shortToken).UpdateColumn("click_count", gorm.Expr("click_count + 1"))
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return result.RowsAffected, result.Error
}
