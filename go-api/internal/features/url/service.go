package url

import (
	"crypto/sha1"
	"encoding/hex"
)

type URLService struct {
	repo URLRepo
}

func NewURLService(repo URLRepo) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) CreateShortToken(original string) (*URLModel, error) {
	ShortToken := generateShortToken(original)
	url := &URLModel{
		Original:   original,
		ShortToken: ShortToken,
	}
	if err := s.repo.Create(url); err != nil {
		return nil, err
	}
	return url, nil
}

func generateShortToken(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])[:8]
}
