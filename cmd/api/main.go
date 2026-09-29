package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/adapters/events"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/adapters/httpapi"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/adapters/postgres"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/application"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("service stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	databaseURL := os.Getenv("DATABASE_URL")
	apiKey := os.Getenv("API_KEY")
	webhookURL := os.Getenv("EVENT_WEBHOOK_URL")
	webhookToken := os.Getenv("EVENT_WEBHOOK_TOKEN")
	if databaseURL == "" || len(apiKey) < 20 {
		return errors.New("DATABASE_URL and API_KEY of at least 20 characters are required")
	}
	if os.Getenv("APP_ENV") == "production" && (webhookURL == "" || webhookToken == "" || len(webhookURL) < 8 || webhookURL[:8] != "https://") {
		return errors.New("production requires HTTPS EVENT_WEBHOOK_URL and EVENT_WEBHOOK_TOKEN")
	}
	lagSeconds := 300
	if value := os.Getenv("OUTBOX_MAX_LAG_SECONDS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return errors.New("OUTBOX_MAX_LAG_SECONDS must be positive")
		}
		lagSeconds = parsed
	}
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return err
	}
	service := application.NewJobService(postgres.NewJobStore(pool))
	handler := httpapi.NewHandler(service, apiKey, func(ctx context.Context) error {
		return postgres.Ready(ctx, pool, time.Duration(lagSeconds)*time.Second)
	})
	server := &http.Server{Addr: address, Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 8192}
	if webhookURL != "" {
		publisher, err := events.NewHTTPPublisher(webhookURL, webhookToken)
		if err != nil {
			return err
		}
		go events.NewWorker(pool, publisher, logger).Run(ctx)
	}
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- server.ListenAndServe() }()
	logger.Info("service started", "address", address)
	select {
	case <-ctx.Done():
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
