package repository

import (
	"context"

	"github.com/Dmitrii30002/url-shortener/internal/domain"
)

type UrlRepository interface {
	Save(ctx context.Context, originalURL string, shortURL string) error
	GetByShortURL(ctx context.Context, shortURL string) (*domain.URL, error)
	GetByOriginalURL(ctx context.Context, originalURL string) (*domain.URL, error)
}
