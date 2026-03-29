package handler

import (
	"net/http"

	"github.com/Dmitrii30002/url-shortener/internal/dto"
	"github.com/Dmitrii30002/url-shortener/internal/errors"
	"github.com/Dmitrii30002/url-shortener/internal/service"
	"github.com/Dmitrii30002/url-shortener/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service  service.Service
	baseURL  string
	validate *validator.Validate
	log      *logger.Logger
}

func NewURLHandler(service service.Service, log *logger.Logger, baseURL string) *Handler {
	return &Handler{
		service:  service,
		baseURL:  baseURL,
		validate: validator.New(),
		log:      log,
	}
}

func (h *Handler) CreateShortURL(c echo.Context) error {
	var req dto.CreateURLRequest

	if err := c.Bind(&req); err != nil {
		h.log.Warnf("Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warnf("Validation failed: %v", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Validation failed",
		})
	}

	shortURL, err := h.service.CreateShortURL(c.Request().Context(), req.URL)
	if err != nil {
		h.log.Errorf("Failed to create short URL: %v", err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	resp := dto.CreateURLResponse{
		ShortURL: h.baseURL + "/" + shortURL,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetOriginalURL(c echo.Context) error {
	shortURL := c.Param("short_url")
	if shortURL == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Short URL is required",
		})
	}

	url, err := h.service.GetOriginalURL(c.Request().Context(), shortURL)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			h.log.Warn("Short URL not found")
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: "Short URL not found",
			})
		default:
			h.log.Errorf("Failed to get original URL: %v", err)
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, dto.GetOriginalURLResponse{
		OriginalURL: url.OriginalURL,
	})
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "url-shortener",
	})
}
