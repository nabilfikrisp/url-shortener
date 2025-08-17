package unit

import (
	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"github.com/stretchr/testify/mock"
)

type MockURLRepo struct {
	mock.Mock
}

func (m *MockURLRepo) Create(u *url.URLModel) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockURLRepo) FindByShortToken(token string) (*url.URLModel, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*url.URLModel), args.Error(1)
}

func (m *MockURLRepo) IncrementClickCount(token string) (int64, error) {
	args := m.Called(token)
	return args.Get(0).(int64), args.Error(1)
}
