package integration

import (
	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"github.com/stretchr/testify/mock"
)

// MockURLService is a testify mock for URLService
type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) CreateShortToken(original string) (*url.URLModel, error) {
	args := m.Called(original)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*url.URLModel), args.Error(1)
}

func (m *MockURLService) FindByShortToken(shortToken string) (*url.URLModel, error) {
	args := m.Called(shortToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*url.URLModel), args.Error(1)
}

func (m *MockURLService) RedirectService(shortToken string) (*url.URLModel, error) {
	args := m.Called(shortToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*url.URLModel), args.Error(1)
}
