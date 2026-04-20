package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BellOriba/go-search-service/internal/api"
	"github.com/BellOriba/go-search-service/internal/database"
	"github.com/BellOriba/go-search-service/internal/products"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting go-search-service", "version", "v0.1.0", "env", "development")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbPool, err := database.NewPostgresPool(ctx)
	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	productRepo := products.NewPostgresRepository(dbPool)

	meiliClient := database.NewMeilisearchClient()

	if err := database.SetupMeilisearchIndex(meiliClient); err != nil {
		slog.Error("failed to setup meilisearch indexes", "error", err)
	}

	searchRepo := products.NewMeilisearchRepository(meiliClient)

	productService := products.NewProductService(productRepo, searchRepo)

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      api.Handler(productService, productRepo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("server started", "port", 8000)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received, draining requests...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}

	slog.Info("server stopped gracefully")
}
