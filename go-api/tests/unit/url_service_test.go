package unit

import (
	"errors"
	"testing"

	"github.com/nabilfikrisp/url-shortener/internal/common/helpers"
	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestURLService(t *testing.T) {
	t.Run("CreateShortToken", func(t *testing.T) {
		t.Run("Returns existing URL if token already exists", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := helpers.GenerateShortToken("https://exists.com")
			existing := &url.URLModel{Original: "https://exists.com", ShortToken: token}

			mockRepo.On("FindByShortToken", token).Return(existing, nil)

			result, err := service.CreateShortToken("https://exists.com")

			assert.NoError(t, err)
			assert.Equal(t, existing, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Success if token does not exist", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := helpers.GenerateShortToken("https://new.com")

			mockRepo.On("FindByShortToken", token).Return(nil, nil)
			mockRepo.On("Create", mock.AnythingOfType("*url.URLModel")).Return(nil)

			result, err := service.CreateShortToken("https://new.com")

			assert.NoError(t, err)
			assert.Equal(t, "https://new.com", result.Original)
			assert.Equal(t, token, result.ShortToken)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Returns error if repo.FindByShortToken fails", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := helpers.GenerateShortToken("https://error.com")

			mockRepo.On("FindByShortToken", token).Return(nil, errors.New("db error"))

			result, err := service.CreateShortToken("https://error.com")

			assert.Error(t, err)
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Returns error if repo.Create fails", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := helpers.GenerateShortToken("https://fail.com")

			mockRepo.On("FindByShortToken", token).Return(nil, nil)
			mockRepo.On("Create", mock.AnythingOfType("*url.URLModel")).Return(errors.New("insert failed"))

			result, err := service.CreateShortToken("https://fail.com")

			assert.Error(t, err)
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("FindByShortToken", func(t *testing.T) {
		t.Run("Success when URL exists", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := "abc123"
			existing := &url.URLModel{Original: "https://example.com", ShortToken: token}

			mockRepo.On("FindByShortToken", token).Return(existing, nil)

			result, err := service.FindByShortToken(token)

			assert.NoError(t, err)
			assert.Equal(t, existing, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Returns error when URL not found", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := "notfound"

			mockRepo.On("FindByShortToken", token).Return(nil, nil)

			result, err := service.FindByShortToken(token)

			assert.Error(t, err)
			assert.Equal(t, "short URL not found", err.Error())
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Returns error when repo fails", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := "error"

			mockRepo.On("FindByShortToken", token).Return(nil, errors.New("db error"))

			result, err := service.FindByShortToken(token)

			assert.Error(t, err)
			assert.Equal(t, "db error", err.Error())
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("RedirectService", func(t *testing.T) {
		t.Run("Success when URL exists and click count increments", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := "abc123"
			existing := &url.URLModel{Original: "https://example.com", ShortToken: token}

			mockRepo.On("FindByShortToken", token).Return(existing, nil)
			mockRepo.On("IncrementClickCount", token).Return(int64(1), nil)

			result, err := service.RedirectService(token)

			assert.NoError(t, err)
			assert.Equal(t, existing, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Returns error when URL not found", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := "notfound"

			mockRepo.On("FindByShortToken", token).Return(nil, nil)

			result, err := service.RedirectService(token)

			assert.Error(t, err)
			assert.Equal(t, "short URL not found", err.Error())
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Returns error when repo FindByShortToken fails", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := "error"

			mockRepo.On("FindByShortToken", token).Return(nil, errors.New("db error"))

			result, err := service.RedirectService(token)

			assert.Error(t, err)
			assert.Equal(t, "db error", err.Error())
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Returns error when increment click count fails", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := "abc123"
			existing := &url.URLModel{Original: "https://example.com", ShortToken: token}

			mockRepo.On("FindByShortToken", token).Return(existing, nil)
			mockRepo.On("IncrementClickCount", token).Return(int64(0), errors.New("update failed"))

			result, err := service.RedirectService(token)

			assert.Error(t, err)
			assert.Equal(t, "update failed", err.Error())
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Returns error when no rows affected by increment", func(t *testing.T) {
			mockRepo := new(MockURLRepo)
			service := url.NewURLService(mockRepo)

			token := "abc123"
			existing := &url.URLModel{Original: "https://example.com", ShortToken: token}

			mockRepo.On("FindByShortToken", token).Return(existing, nil)
			mockRepo.On("IncrementClickCount", token).Return(int64(0), nil)

			result, err := service.RedirectService(token)

			assert.Error(t, err)
			assert.Equal(t, "unable to update click statistics", err.Error())
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})
	})
}
