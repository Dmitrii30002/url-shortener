package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dmitrii30002/url-shortener/internal/config"
	"github.com/Dmitrii30002/url-shortener/internal/handler"
	"github.com/Dmitrii30002/url-shortener/internal/repository"
	"github.com/Dmitrii30002/url-shortener/internal/service"
	"github.com/Dmitrii30002/url-shortener/pkg/generator"
	"github.com/Dmitrii30002/url-shortener/pkg/logger"
	"github.com/Dmitrii30002/url-shortener/pkg/migrator"
	"github.com/Dmitrii30002/url-shortener/pkg/storage/postgres"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.GetConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
		return
	}

	log, err := logger.New(&cfg.LoggerCfg)
	if err != nil {
		fmt.Printf("Failed to setup logger: %v", err)
		return
	}

	var repo repository.UrlRepository
	switch cfg.StorageType {
	case "memory":
		repo = repository.NewMemoryRepository()
	default:
		db, err := postgres.New(&cfg.PostgresCfg)
		if err != nil {
			log.Fatalf("Failed to setup postgres: %v", err)
			return
		}

		err = migrator.Up(db, "migrations")
		if err != nil {
			log.Fatalf("Failed to migrate: %v", err)
			return
		}

		repo = repository.NewURLRepositoryPostgres(db)
	}

	service := service.NewService(repo, generator.New())

	e := echo.New()
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 5 * time.Second

	urlHandler := handler.NewURLHandler(service, log, cfg.BaseURL)
	e.POST("/short", urlHandler.CreateShortURL)
	e.GET("/:short_url", urlHandler.GetOriginalURL)
	e.GET("/health", urlHandler.HealthCheck)

	go func() {
		log.Info("Starting server")
		if err := e.Start(cfg.Server.Host + ":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed")
	}

	log.Info("Server stopped")
}
