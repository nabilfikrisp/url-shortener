package url

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

type URLService interface {
	CreateShortToken(original string) (*URLModel, error)
	FindByShortToken(shortToken string) (*URLModel, error)
	RedirectService(shortToken string) (*URLModel, error)
}
type urlService struct {
	repo URLRepo
}

func NewURLService(repo URLRepo) URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) CreateShortToken(original string) (*URLModel, error) {
	shortToken := generateShortToken(original)

	existingURL, err := s.repo.FindByShortToken(shortToken)
	if err != nil {
		return nil, err
	}
	if existingURL != nil {
		return existingURL, nil
	}

	url := &URLModel{
		Original:   original,
		ShortToken: shortToken,
	}
	if err := s.repo.Create(url); err != nil {
		return nil, err
	}
	return url, nil
}

func (s *urlService) FindByShortToken(shortToken string) (*URLModel, error) {
	url, err := s.repo.FindByShortToken(shortToken)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, errors.New("short URL not found")
	}
	return url, nil
}

func (s *urlService) RedirectService(shortToken string) (*URLModel, error) {
	url, err := s.repo.FindByShortToken(shortToken)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, errors.New("short URL not found")
	}

	affectedRows, err := s.repo.IncrementClickCount(shortToken)
	if err != nil {
		return nil, err
	}
	if affectedRows == 0 {
		return nil, errors.New("unable to update click statistics")
	}

	return url, nil
}

func generateShortToken(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])[:16]
}
