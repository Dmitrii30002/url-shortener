package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Dmitrii30002/url-shortener/internal/config"
	"github.com/Dmitrii30002/url-shortener/internal/errors"
	"github.com/Dmitrii30002/url-shortener/pkg/migrator"
	"github.com/Dmitrii30002/url-shortener/pkg/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *urlRepositoryPostgres {
	cfg, _ := config.GetTestConfig()
	db, err := postgres.New(&cfg.PostgresCfg)
	require.NoError(t, err)
	defer cleanDB(t, db)

	err = migrator.Up(db, "../../migrations")
	require.NoError(t, err)

	repo := NewURLRepositoryPostgres(db)

	return repo.(*urlRepositoryPostgres)
}

func cleanDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE urls")
	require.NoError(t, err)
}

func TestPostgresSave_ShouldStoreURL_WhenOriginalAndShortURLAreNew(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanDB(t, repo.db)

	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123"

	err := repo.Save(ctx, originalURL, shortURL)

	assert.NoError(t, err)
}

func TestPostgresSave_ShouldReturnError_WhenOriginalURLAlreadyExists(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanDB(t, repo.db)

	ctx := context.Background()

	originalURL := "https://example.com"

	err := repo.Save(ctx, originalURL, "abc123")
	require.NoError(t, err)

	err = repo.Save(ctx, originalURL, "xyz789")

	assert.Equal(t, errors.ErrDuplicateURL, err)
}

func TestPostgresSave_ShouldReturnError_WhenShortURLAlreadyExists(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanDB(t, repo.db)

	ctx := context.Background()

	shortURL := "abc123"

	err := repo.Save(ctx, "https://example1.com", shortURL)
	require.NoError(t, err)

	err = repo.Save(ctx, "https://example2.com", shortURL)

	assert.Equal(t, errors.ErrDuplicateShortURL, err)
}

func TestPostgresGetByShortURL_ShouldReturnURL_WhenShortURLExists(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanDB(t, repo.db)

	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123"

	err := repo.Save(ctx, originalURL, shortURL)
	require.NoError(t, err)

	url, err := repo.GetByShortURL(ctx, shortURL)

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, originalURL, url.OriginalURL)
	assert.Equal(t, shortURL, url.ShortURL)
}

func TestPostgresGetByShortURL_ShouldReturnNotFound_WhenShortURLMissing(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanDB(t, repo.db)

	ctx := context.Background()

	url, err := repo.GetByShortURL(ctx, "notexist")

	assert.Equal(t, errors.ErrNotFound, err)
	assert.Nil(t, url)
}

func TestPostgresGetByOriginalURL_ShouldReturnURL_WhenOriginalURLExists(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanDB(t, repo.db)

	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123"

	err := repo.Save(ctx, originalURL, shortURL)
	require.NoError(t, err)

	url, err := repo.GetByOriginalURL(ctx, originalURL)

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, originalURL, url.OriginalURL)
	assert.Equal(t, shortURL, url.ShortURL)
}

func TestPostgresGetByOriginalURL_ShouldReturnNotFound_WhenOriginalURLMissing(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanDB(t, repo.db)

	ctx := context.Background()

	url, err := repo.GetByOriginalURL(ctx, "https://notexist.com")

	assert.Equal(t, errors.ErrNotFound, err)
	assert.Nil(t, url)
}
