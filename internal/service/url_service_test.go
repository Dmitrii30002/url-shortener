package service

import (
	"context"
	"testing"

	"github.com/Dmitrii30002/url-shortener/internal/domain"
	"github.com/Dmitrii30002/url-shortener/internal/errors"
	"github.com/stretchr/testify/assert"
)

type MockGenerator struct {
	nextCode string
}

func (m *MockGenerator) Generate() string {
	return m.nextCode
}

type MockRepository struct {
	saveErr                error
	getByShortURLResult    *domain.URL
	getByShortURLErr       error
	getByOriginalURLResult *domain.URL
	getByOriginalURLErr    error
}

func (m *MockRepository) Save(ctx context.Context, originalURL, shortURL string) error {
	return m.saveErr
}

func (m *MockRepository) GetByShortURL(ctx context.Context, shortURL string) (*domain.URL, error) {
	return m.getByShortURLResult, m.getByShortURLErr
}

func (m *MockRepository) GetByOriginalURL(ctx context.Context, originalURL string) (*domain.URL, error) {
	return m.getByOriginalURLResult, m.getByOriginalURLErr
}

func TestCreateShortURL_ShouldReturnShortURL_WhenSaveSucceeds(t *testing.T) {
	mockRepo := &MockRepository{saveErr: nil}
	mockGen := &MockGenerator{nextCode: "abc123"}

	service := NewService(mockRepo, mockGen)
	ctx := context.Background()

	shortURL, err := service.CreateShortURL(ctx, "https://example.com")

	assert.NoError(t, err)
	assert.Equal(t, "abc123", shortURL)
}

func TestCreateShortURL_ShouldRetry_WhenShortURLAlreadyExists(t *testing.T) {
	mockRepo := &MockRepository{saveErr: errors.ErrDuplicateShortURL}
	mockGen := &MockGenerator{nextCode: "abc123"}

	service := NewService(mockRepo, mockGen)
	ctx := context.Background()

	shortURL, err := service.CreateShortURL(ctx, "https://example.com")

	assert.Error(t, err)
	assert.Empty(t, shortURL)
	assert.Contains(t, err.Error(), "failed to generate short url after 5 attempts")
}

func TestCreateShortURL_ShouldReturnExisting_WhenOriginalURLExists(t *testing.T) {
	mockRepo := &MockRepository{
		saveErr: errors.ErrDuplicateURL,
		getByOriginalURLResult: &domain.URL{
			OriginalURL: "https://example.com",
			ShortURL:    "existing123",
		},
		getByOriginalURLErr: nil,
	}
	mockGen := &MockGenerator{nextCode: "abc123"}

	service := NewService(mockRepo, mockGen)
	ctx := context.Background()

	shortURL, err := service.CreateShortURL(ctx, "https://example.com")

	assert.NoError(t, err)
	assert.Equal(t, "existing123", shortURL)
}

func TestCreateShortURL_ShouldFail_WhenFetchingExistingURLFails(t *testing.T) {
	mockRepo := &MockRepository{
		saveErr:             errors.ErrDuplicateURL,
		getByOriginalURLErr: errors.ErrNotFound,
	}
	mockGen := &MockGenerator{nextCode: "abc123"}

	service := NewService(mockRepo, mockGen)
	ctx := context.Background()

	shortURL, err := service.CreateShortURL(ctx, "https://example.com")

	assert.Error(t, err)
	assert.Empty(t, shortURL)
}

func TestCreateShortURL_ShouldFail_WhenUnexpectedErrorOccurs(t *testing.T) {
	mockRepo := &MockRepository{saveErr: assert.AnError}
	mockGen := &MockGenerator{nextCode: "abc123"}

	service := NewService(mockRepo, mockGen)
	ctx := context.Background()

	shortURL, err := service.CreateShortURL(ctx, "https://example.com")

	assert.Error(t, err)
	assert.Empty(t, shortURL)
	assert.Contains(t, err.Error(), "failed to save url")
}

func TestGetOriginalURL_ShouldReturnURL_WhenShortURLExists(t *testing.T) {
	mockRepo := &MockRepository{
		getByShortURLResult: &domain.URL{
			OriginalURL: "https://example.com",
			ShortURL:    "abc123",
		},
		getByShortURLErr: nil,
	}
	mockGen := &MockGenerator{nextCode: ""}

	service := NewService(mockRepo, mockGen)
	ctx := context.Background()

	url, err := service.GetOriginalURL(ctx, "abc123")

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, "https://example.com", url.OriginalURL)
	assert.Equal(t, "abc123", url.ShortURL)
}

func TestGetOriginalURL_ShouldReturnNotFound_WhenShortURLMissing(t *testing.T) {
	mockRepo := &MockRepository{
		getByShortURLResult: nil,
		getByShortURLErr:    errors.ErrNotFound,
	}
	mockGen := &MockGenerator{nextCode: ""}

	service := NewService(mockRepo, mockGen)
	ctx := context.Background()

	url, err := service.GetOriginalURL(ctx, "notexist")

	assert.Equal(t, errors.ErrNotFound, err)
	assert.Nil(t, url)
}

func TestGetOriginalURL_ShouldFail_WhenDatabaseErrorOccurs(t *testing.T) {
	mockRepo := &MockRepository{
		getByShortURLResult: nil,
		getByShortURLErr:    assert.AnError,
	}
	mockGen := &MockGenerator{nextCode: ""}

	service := NewService(mockRepo, mockGen)
	ctx := context.Background()

	url, err := service.GetOriginalURL(ctx, "abc123")

	assert.Error(t, err)
	assert.Nil(t, url)
	assert.Contains(t, err.Error(), "failed to find url")
}
