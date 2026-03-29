package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Dmitrii30002/url-shortener/internal/domain"
	domainErrors "github.com/Dmitrii30002/url-shortener/internal/errors"
)

type urlRepositoryPostgres struct {
	db *sql.DB
}

func NewURLRepositoryPostgres(db *sql.DB) UrlRepository {
	return &urlRepositoryPostgres{db: db}
}

func (r *urlRepositoryPostgres) Save(ctx context.Context, originalURL, shortURL string) error {
	query := `
		INSERT INTO urls (original_url, short_url) 
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(ctx, query, originalURL, shortURL)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			if strings.Contains(err.Error(), "original_url") {
				return domainErrors.ErrDuplicateURL
			}
			if strings.Contains(err.Error(), "short_url") {
				return domainErrors.ErrDuplicateShortURL
			}
		}

		return err
	}

	return nil
}

func (r *urlRepositoryPostgres) GetByShortURL(ctx context.Context, shortURL string) (*domain.URL, error) {
	query := `
		SELECT original_url, short_url 
		FROM urls 
		WHERE short_url = $1
	`
	row := r.db.QueryRowContext(ctx,
		query,
		shortURL,
	)

	var url domain.URL
	if err := row.Scan(&url.OriginalURL, &url.ShortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrNotFound
		}

		return nil, err
	}

	return &url, nil
}

func (r *urlRepositoryPostgres) GetByOriginalURL(ctx context.Context, originalURL string) (*domain.URL, error) {
	query := `
		SELECT original_url, short_url 
		FROM urls 
		WHERE original_url = $1
	`
	row := r.db.QueryRowContext(ctx,
		query,
		originalURL,
	)

	var url domain.URL
	if err := row.Scan(&url.OriginalURL, &url.ShortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrNotFound
		}

		return nil, err
	}

	return &url, nil
}
