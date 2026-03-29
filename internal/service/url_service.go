package service

import (
	"context"
	"fmt"

	"github.com/Dmitrii30002/url-shortener/internal/domain"
	"github.com/Dmitrii30002/url-shortener/internal/errors"
	"github.com/Dmitrii30002/url-shortener/internal/repository"
	"github.com/Dmitrii30002/url-shortener/pkg/generator"
)

type Service interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortURL string) (*domain.URL, error)
}

type urlService struct {
	repo      repository.UrlRepository
	generator generator.Generator
}

func NewService(repo repository.UrlRepository, generator generator.Generator) Service {
	return &urlService{
		repo:      repo,
		generator: generator,
	}
}

func (s *urlService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	var shortURL string
	var lastErr error
	for attempts := 0; attempts < 3; attempts++ {
		shortURL = s.generator.Generate()

		err := s.repo.Save(ctx, originalURL, shortURL)
		switch err {
		case nil:
			return shortURL, nil
		case errors.ErrDuplicateShortURL:
			lastErr = err
		case errors.ErrDuplicateURL:
			url, err := s.repo.GetByOriginalURL(ctx, originalURL)
			if err != nil {
				return "", err
			}
			return url.ShortURL, nil
		default:
			return "", fmt.Errorf("failed to save url: %w", err)
		}
	}

	return "", fmt.Errorf("failed to generate short url after 3 attempts: %w", lastErr)
}

func (s *urlService) GetOriginalURL(ctx context.Context, shortURL string) (*domain.URL, error) {
	url, err := s.repo.GetByShortURL(ctx, shortURL)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find url: %w", err)
	}

	return url, nil
}
