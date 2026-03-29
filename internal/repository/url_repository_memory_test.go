package repository

import (
	"context"
	"testing"

	"github.com/Dmitrii30002/url-shortener/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave_ShouldStoreURL_WhenOriginalAndShortURLAreNew(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123"

	err := repo.Save(ctx, originalURL, shortURL)

	assert.NoError(t, err)

	savedShortURL, _ := repo.urlToShort.Get(originalURL)
	assert.Equal(t, shortURL, savedShortURL)

	savedOriginalURL, _ := repo.shortToURL.Get(shortURL)
	assert.Equal(t, originalURL, savedOriginalURL)
}

func TestSave_ShouldReturnError_WhenOriginalURLAlreadyExists(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	originalURL := "https://example.com"

	err := repo.Save(ctx, originalURL, "abc123")
	require.NoError(t, err)

	err = repo.Save(ctx, originalURL, "xyz789")

	assert.Equal(t, errors.ErrDuplicateURL, err)

	savedShortURL, _ := repo.urlToShort.Get(originalURL)
	assert.Equal(t, "abc123", savedShortURL)
}

func TestSave_ShouldReturnError_WhenShortURLAlreadyExists(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	shortURL := "abc123"

	err := repo.Save(ctx, "https://example1.com", shortURL)
	require.NoError(t, err)

	err = repo.Save(ctx, "https://example2.com", shortURL)

	assert.Equal(t, errors.ErrDuplicateShortURL, err)

	ok := repo.urlToShort.Exist("https://example2.com")
	assert.False(t, ok)
}

func TestGetByShortURL_ShouldReturnURL_WhenShortURLExists(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123"

	err := repo.Save(ctx, originalURL, shortURL)
	require.NoError(t, err)

	url, err := repo.GetByShortURL(ctx, shortURL)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, url.OriginalURL)
	assert.Equal(t, shortURL, url.ShortURL)
}

func TestGetByShortURL_ShouldReturnNotFound_WhenShortURLMissing(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	url, err := repo.GetByShortURL(ctx, "notexist")

	assert.Equal(t, errors.ErrNotFound, err)
	assert.Nil(t, url)
}

func TestGetByOriginalURL_ShouldReturnURL_WhenOriginalURLExists(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123"

	err := repo.Save(ctx, originalURL, shortURL)
	require.NoError(t, err)

	url, err := repo.GetByOriginalURL(ctx, originalURL)

	assert.Equal(t, originalURL, url.OriginalURL)
	assert.Equal(t, shortURL, url.ShortURL)
}

func TestGetByOriginalURL_ShouldReturnNotFound_WhenOriginalURLMissing(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	url, err := repo.GetByOriginalURL(ctx, "https://notexist.com")

	assert.Equal(t, errors.ErrNotFound, err)
	assert.Nil(t, url)
}
