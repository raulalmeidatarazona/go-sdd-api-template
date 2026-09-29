//go:build integration

package events

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/adapters/postgres"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/application"
)

type stubPublisher struct {
	fail  bool
	calls int
}

func (publisher *stubPublisher) Publish(_ context.Context, _, _ string, _ []byte) error {
	publisher.calls++
	if publisher.fail {
		return errors.New("receiver unavailable")
	}
	return nil
}

func TestOutboxRetriesAndReadinessRecovers(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if _, err := pool.Exec(ctx, `TRUNCATE outbox,jobs`); err != nil {
		t.Fatal(err)
	}
	service := application.NewJobService(postgres.NewJobStore(pool))
	job, _, err := service.Create(ctx, "retry example", "retry-test-"+time.Now().Format(time.RFC3339Nano))
	if err != nil {
		t.Fatal(err)
	}
	var eventID string
	if err := pool.QueryRow(ctx, `SELECT id FROM outbox WHERE payload->>'jobId'=$1`, job.ID).Scan(&eventID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `UPDATE outbox SET created_at=now()-interval '1 hour' WHERE id=$1`, eventID); err != nil {
		t.Fatal(err)
	}
	if err := postgres.Ready(ctx, pool, 5*time.Minute); err == nil {
		t.Fatal("stale outbox must fail readiness")
	}
	publisher := &stubPublisher{fail: true}
	worker := NewWorker(pool, publisher, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err := worker.deliverOne(ctx); err != nil {
		t.Fatal(err)
	}
	var attempts int
	if err := pool.QueryRow(ctx, `SELECT attempts FROM outbox WHERE id=$1`, eventID).Scan(&attempts); err != nil || attempts != 1 {
		t.Fatalf("attempts=%d error=%v", attempts, err)
	}
	if _, err := pool.Exec(ctx, `UPDATE outbox SET available_at=now() WHERE id=$1`, eventID); err != nil {
		t.Fatal(err)
	}
	publisher.fail = false
	if err := worker.deliverOne(ctx); err != nil {
		t.Fatal(err)
	}
	var delivered bool
	if err := pool.QueryRow(ctx, `SELECT delivered_at IS NOT NULL FROM outbox WHERE id=$1`, eventID).Scan(&delivered); err != nil || !delivered {
		t.Fatalf("delivered=%v error=%v", delivered, err)
	}
	if publisher.calls != 2 {
		t.Fatalf("publish attempts=%d", publisher.calls)
	}
	if err := postgres.Ready(ctx, pool, 5*time.Minute); err != nil {
		t.Fatalf("readiness should recover: %v", err)
	}
}
