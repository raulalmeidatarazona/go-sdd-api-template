package events

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Publisher interface {
	Publish(context.Context, string, string, []byte) error
}

type HTTPPublisher struct {
	endpoint string
	token    string
	client   *http.Client
}

func NewHTTPPublisher(endpoint, token string) (*HTTPPublisher, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, errors.New("invalid event webhook URL")
	}
	return &HTTPPublisher{endpoint: endpoint, token: token, client: &http.Client{Timeout: 5 * time.Second}}, nil
}

func (publisher *HTTPPublisher) Publish(ctx context.Context, id, eventType string, payload []byte) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, publisher.endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Event-Id", id)
	request.Header.Set("X-Event-Type", eventType)
	if publisher.token != "" {
		request.Header.Set("Authorization", "Bearer "+publisher.token)
	}
	response, err := publisher.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("receiver returned %d", response.StatusCode)
	}
	return nil
}

type Worker struct {
	pool      *pgxpool.Pool
	publisher Publisher
	logger    *slog.Logger
}

func NewWorker(pool *pgxpool.Pool, publisher Publisher, logger *slog.Logger) *Worker {
	return &Worker{pool: pool, publisher: publisher, logger: logger}
}

func (worker *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		if err := worker.deliverOne(ctx); err != nil && ctx.Err() == nil {
			worker.logger.Error("outbox delivery failed", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (worker *Worker) deliverOne(ctx context.Context) error {
	tx, err := worker.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var id, eventType string
	var payload []byte
	var attempts int
	err = tx.QueryRow(ctx, `SELECT id,event_type,payload,attempts FROM outbox
        WHERE delivered_at IS NULL AND available_at <= now()
        ORDER BY available_at,id LIMIT 1 FOR UPDATE SKIP LOCKED`).Scan(&id, &eventType, &payload, &attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return tx.Commit(ctx)
	}
	if err != nil {
		return err
	}
	publishContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	publishErr := worker.publisher.Publish(publishContext, id, eventType, payload)
	cancel()
	if publishErr == nil {
		_, err = tx.Exec(ctx, `UPDATE outbox SET delivered_at=now(),last_error=NULL WHERE id=$1`, id)
	} else {
		backoff := time.Duration(1<<min(attempts, 8)) * time.Second
		_, err = tx.Exec(ctx, `UPDATE outbox SET attempts=attempts+1,available_at=now()+$2::interval,last_error=$3 WHERE id=$1`, id, backoff.String(), "receiver unavailable")
	}
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	if publishErr != nil {
		worker.logger.Warn("event queued for retry", "event_id", id, "attempt", attempts+1, "error", publishErr)
	}
	return nil
}
