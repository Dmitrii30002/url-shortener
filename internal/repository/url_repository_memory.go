package repository

import (
	"context"

	"github.com/Dmitrii30002/url-shortener/internal/domain"
	"github.com/Dmitrii30002/url-shortener/internal/errors"
	"github.com/Dmitrii30002/url-shortener/pkg/storage/memory"
)

type MemoryRepository struct {
	urlToShort *memory.Store
	shortToURL *memory.Store
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		urlToShort: memory.New(),
		shortToURL: memory.New(),
	}
}

func (r *MemoryRepository) Save(ctx context.Context, originalURL string, shortURL string) error {
	if r.urlToShort.Exist(originalURL) {
		return errors.ErrDuplicateURL
	}

	if r.shortToURL.Exist(shortURL) {
		return errors.ErrDuplicateShortURL
	}

	r.urlToShort.Set(originalURL, shortURL)
	r.shortToURL.Set(shortURL, originalURL)

	return nil
}

func (r *MemoryRepository) GetByShortURL(ctx context.Context, shortURL string) (*domain.URL, error) {
	originalURL, ok := r.shortToURL.Get(shortURL)
	if !ok {
		return nil, errors.ErrNotFound
	}

	return &domain.URL{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}, nil
}

func (r *MemoryRepository) GetByOriginalURL(ctx context.Context, originalURL string) (*domain.URL, error) {
	shortURL, ok := r.urlToShort.Get(originalURL)
	if !ok {
		return nil, errors.ErrNotFound
	}

	return &domain.URL{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}, nil
}
