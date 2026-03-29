package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dmitrii30002/url-shortener/internal/domain"
	"github.com/Dmitrii30002/url-shortener/internal/dto"
	"github.com/Dmitrii30002/url-shortener/internal/errors"
	"github.com/Dmitrii30002/url-shortener/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockService struct {
	createShortURLResult string
	createShortURLErr    error
	getOriginalURLResult *domain.URL
	getOriginalURLErr    error
}

func (m *MockService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	return m.createShortURLResult, m.createShortURLErr
}

func (m *MockService) GetOriginalURL(ctx context.Context, shortURL string) (*domain.URL, error) {
	return m.getOriginalURLResult, m.getOriginalURLErr
}

func setupTestHandler(mockService *MockService) (*Handler, *echo.Echo) {
	log := logger.NewTestLogger()
	handler := NewURLHandler(mockService, log, "http://localhost:8080")
	e := echo.New()
	return handler, e
}

func TestCreateShortURL_ShouldReturnShortURL_WhenValidRequest(t *testing.T) {
	mockService := &MockService{
		createShortURLResult: "abc123",
		createShortURLErr:    nil,
	}
	handler, e := setupTestHandler(mockService)

	reqBody := dto.CreateURLRequest{URL: "https://example.com"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateShortURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.CreateURLResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/abc123", resp.ShortURL)
}

func TestCreateShortURL_ShouldReturnBadRequest_WhenRequestBodyIsInvalid(t *testing.T) {
	mockService := &MockService{}
	handler, e := setupTestHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateShortURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Invalid request body", resp.Error)
}

func TestCreateShortURL_ShouldReturnBadRequest_WhenURLIsMissing(t *testing.T) {
	mockService := &MockService{}
	handler, e := setupTestHandler(mockService)

	reqBody := dto.CreateURLRequest{URL: ""}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateShortURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Validation failed", resp.Error)
}

func TestCreateShortURL_ShouldReturnInternalError_WhenServiceFails(t *testing.T) {
	mockService := &MockService{
		createShortURLResult: "",
		createShortURLErr:    assert.AnError,
	}
	handler, e := setupTestHandler(mockService)

	reqBody := dto.CreateURLRequest{URL: "https://example.com"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateShortURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetOriginalURL_ShouldReturnOriginalURL_WhenShortURLExists(t *testing.T) {
	mockService := &MockService{
		getOriginalURLResult: &domain.URL{
			OriginalURL: "https://example.com",
			ShortURL:    "abc123",
		},
		getOriginalURLErr: nil,
	}
	handler, e := setupTestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("short_url")
	c.SetParamValues("abc123")

	err := handler.GetOriginalURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.GetOriginalURLResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", resp.OriginalURL)
}

func TestGetOriginalURL_ShouldReturnNotFound_WhenShortURLMissing(t *testing.T) {
	mockService := &MockService{
		getOriginalURLResult: nil,
		getOriginalURLErr:    errors.ErrNotFound,
	}
	handler, e := setupTestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/notexist", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("short_url")
	c.SetParamValues("notexist")

	err := handler.GetOriginalURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Short URL not found", resp.Error)
}

func TestGetOriginalURL_ShouldReturnInternalError_WhenServiceFails(t *testing.T) {
	mockService := &MockService{
		getOriginalURLResult: nil,
		getOriginalURLErr:    assert.AnError,
	}
	handler, e := setupTestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("short_url")
	c.SetParamValues("abc123")

	err := handler.GetOriginalURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetOriginalURL_ShouldReturnBadRequest_WhenShortURLIsEmpty(t *testing.T) {
	mockService := &MockService{}
	handler, e := setupTestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("short_url")
	c.SetParamValues("")

	err := handler.GetOriginalURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Short URL is required", resp.Error)
}

func TestHealthCheck_ShouldReturnOk(t *testing.T) {
	mockService := &MockService{}
	handler, e := setupTestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HealthCheck(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
	assert.Equal(t, "url-shortener", resp["service"])
}
